package handlers

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	// "bytes"
	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-wire"
	apiconf "gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/config"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func responseWrite(ctx *gin.Context, isSuccess bool, result interface{}) {
	ret := gin.H{
		"isSuccess": isSuccess,
	}

	// 用在api->midware的情形
	// response, _ := result.(map[string]interface{})
	// if response["codetype"] == at.CodeType_Timeout || response["codetype"] == at.CodeType_NonceTooLow {
	// 	ret["isSuccess"] = false
	// }

	if isSuccess {
		ret["result"] = result
	} else {
		ret["message"] = result
	}

	ctx.JSON(200, ret)

	fmt.Printf("===========raw request url: %s\n", ctx.Request.URL.String())
	fmt.Printf("===========raw response result: %v\n", result)
}

func (hd *Handler) SignAndSendTx(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey) {
	sigTx, err := tx.Sign(privkey)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	ret, err := hd.jsonRPC("broadcast_tx_sync", txBytes)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (ret).(*at.ResultBroadcastTx)

	//fmt.Printf("resultbroadcastTx : %v", string(res.Data[:]))

	if res.Code != at.CodeType_OK {
		responseWrite(ctx, false, string(res.Log))
		return
	}

	response := map[string]interface{}{
		"res": string(res.Data),
		"tx":  sigTx.Hash().Hex(),
	}
	bigDataHashes, exist := ctx.Get("bigdata_hashs")
	if exist {
		if data := bigDataHashes.(*[]map[string]interface{}); len(*data) > 0 {
			response["bigDatas"] = *data
		}
	}

	responseWrite(ctx, true, response)
}

func (hd *Handler) SignSendToCommitTx(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey) {
	var (
		code    int
		message string
	)
	sigTx, err := tx.Sign(privkey)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	rptRes := new(at.RPCResult)
	err = hd.sendTxCall("broadcast_tx_sync", []interface{}{txBytes}, rptRes)

	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (*rptRes).(*at.ResultBroadcastTx)

	if res.Code != at.CodeType_OK {
		responseWrite(ctx, false, string(res.Log))
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	chExit := time.After(time.Second * time.Duration(apiconf.TimeOut()))

	for {
		select {
		case <-ticker.C:
			code, message = hd.queryTxForHash(sigTx.Hash())
			switch code {
			case 0:
				response := map[string]interface{}{
					"res": string(res.Data),
					"tx":  sigTx.Hash().Hex(),
				}
				bigDataHashes, exist := ctx.Get("bigdata_hashs")
				if exist {
					if data := bigDataHashes.(*[]map[string]interface{}); len(*data) > 0 {
						response["bigDatas"] = *data
					}
				}
				responseWrite(ctx, true, response)
				return
			case 1:
				continue
			default:
				responseWrite(ctx, false, message)
				return
			}
		case <-chExit:
			responseWrite(ctx, false, fmt.Sprintf("hash:%v,%v", sigTx.Hash().Hex(), message))
			return
		}
	}
	return
}

func (hd *Handler) SignAndCommitTx(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey) {
	// fmt.Println("tx.Hash().Hex()", tx.Hash().Hex())
	// 加签
	sigTx, err := tx.Sign(privkey)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	// fmt.Println("sigTx.Hash().Hex()", sigTx.Hash().Hex())
	// 如果从redis缓存里取到结果，则直接返回给前端
	lhash := sigTx.Hash().Hex()
	r_key := fmt.Sprintf("%s%s", Redis_Key, lhash)
	result, _ := hd.redisApi.HGet(r_key, Redis_Result)
	if len(result) > 0 {
		fmt.Println("get from redis: txHash(%s)", r_key)
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.Write(result)
		return
	}

	hd.recursiveTransactions(ctx, sigTx, 3)
}

// 为了解决交易超时的bug，而让交易递归地进行。 递归的次数为n
func (hd *Handler) recursiveTransactions(ctx *gin.Context, sigTx *types.Transaction, n int) {
	fmt.Println("recursiveTransactions %s ", n)
	if n <= 0 {
		responseWrite(ctx, false, fmt.Sprintf("recursiveTransactions more than %s times", n))
		return
	}

	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	ret, err := hd.jsonRPC("broadcast_tx_commit", txBytes)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (ret).(*at.ResultBroadcastTxCommit)
	if res.Code != at.CodeType_OK && res.Code != at.CodeType_Timeout && res.Code != at.CodeType_NonceTooLow {
		responseWrite(ctx, false, string(res.Log))
		return
	}

	txHash := sigTx.Hash().Hex()
	response, ctype, err := hd.handleTimeout(ctx, res, txHash)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	if ctype == at.CodeType_Timeout { // 说明交易超时的情况下，还没查到交易结果
		hd.recursiveTransactions(ctx, sigTx, n-1)
	} else if ctype == at.CodeType_NonceTooLow { // 在nonce too low的情况下
		responseWrite(ctx, false, "nonce too low")
	} else { // == at.CodeType_Ok
		ret := gin.H{
			"isSuccess": true,
		}
		ret["result"] = response
		result, _ := json.Marshal(ret)
		r_key := fmt.Sprintf("%s%s", Redis_Key, txHash)
		hd.redisApi.HInsert(r_key, Redis_Result, result, 600)
		responseWrite(ctx, true, response)
	}
}

func (hd *Handler) handleTimeout(ctx *gin.Context, res *at.ResultBroadcastTxCommit, txHash string) (interface{}, at.CodeType, error) {
	if res.Code != at.CodeType_Timeout && res.Code != at.CodeType_NonceTooLow {
		response := map[string]interface{}{
			"res":      string(res.Data),
			"tx":       txHash,
			"codetype": res.Code,
		}
		bigDataHashes, exist := ctx.Get("bigdata_hashs")
		if exist {
			if data := bigDataHashes.(*[]map[string]interface{}); len(*data) > 0 {
				response["bigDatas"] = *data
			}
		}
		return response, at.CodeType_OK, nil
	}

	var n int
	if res.Code == at.CodeType_Timeout {
		n = 5
	} else {
		n = 1
	}
	cCodetype := make(chan at.CodeType)
	go hd.QueryRes(ctx, cCodetype, txHash, n)
	codetype := <-cCodetype

	if codetype == at.CodeType_Timeout {
		// 查了几次，最终结果还是超时
		fmt.Println("handleTimeout CodeType_Timeout")
		return nil, res.Code, nil

	} else {
		// 有交易结果了,在redis缓存结果后，返回给前端
		fmt.Println("handleTimeout CodeType_Ok")
		response := map[string]interface{}{
			"res":      "",
			"tx":       txHash,
			"codetype": at.CodeType_OK,
		}
		bigDataHashes, exist := ctx.Get("bigdata_hashs")
		if exist {
			if data := bigDataHashes.(*[]map[string]interface{}); len(*data) > 0 {
				response["bigDatas"] = *data
			}
		}
		return response, at.CodeType_OK, nil
	}
}

func (hd *Handler) QueryRes(ctx *gin.Context, cCodetype chan at.CodeType, txhash string, n int) {
	for i := 0; i < n; i++ {
		fmt.Println("QueryRes %s ", i)

		b := hd.isTxExist(ctx, txhash)
		// 查到结果，说明交易已经执行
		if b {
			cCodetype <- at.CodeType_OK
			return
		}
		time.Sleep(time.Second)
	}

	cCodetype <- at.CodeType_Timeout
}

func (hd *Handler) SignAndSendTxForBigData(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey, hashs string) {
	sigTx, err := tx.Sign(privkey)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	rptRes := new(at.RPCResult)
	err = hd.sendTxCall("broadcast_tx_sync", []interface{}{txBytes}, rptRes)

	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (*rptRes).(*at.ResultBroadcastTx)

	//fmt.Printf("resultbroadcastTx : %v", string(res.Data[:]))

	if res.Code != at.CodeType_OK {
		responseWrite(ctx, false, string(res.Log))
		return
	}

	response := map[string]interface{}{
		"res":   string(res.Data),
		"tx":    sigTx.Hash().Hex(),
		"hashs": hashs,
	}
	responseWrite(ctx, true, response)
}

func (hd *Handler) SignAndCommitTxForBigData(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey, hashs string) {
	sigTx, err := tx.Sign(privkey)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	rptRes := new(at.RPCResult)
	err = hd.sendTxCall("broadcast_tx_commit", []interface{}{txBytes}, rptRes)

	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (*rptRes).(*at.ResultBroadcastTxCommit)

	//fmt.Printf("resultbroadcastTx : %v", string(res.Data[:]))

	if res.Code != at.CodeType_OK {
		responseWrite(ctx, false, string(res.Log))
		return
	}

	response := map[string]interface{}{
		"res":   string(res.Data),
		"tx":    sigTx.Hash().Hex(),
		"hashs": hashs,
	}
	responseWrite(ctx, true, response)
}

func parseAsset(ctx *gin.Context) (selling, buying types.Asset, err error) {
	var temp uint64

	sellingAssetType := ctx.Query("selling_asset_type")
	if sellingAssetType == "" {
		selling.SetNull()
	} else {
		temp, err = strconv.ParseUint(sellingAssetType, 10, 64)
		if err != nil {
			return
		}
		selling.Type = uint8(temp)
		if selling.Type == types.ASSET_TYPE_CREDIT_ALPHANUM4 || selling.Type == types.ASSET_TYPE_CREDIT_ALPHANUM12 {
			selling.Code = ctx.Query("selling_asset_code")
			selling.Issuer = ethcmn.HexToAddress(ctx.Query("selling_asset_issuer"))
		}
	}

	buyingAssetType := ctx.Query("buying_asset_type")
	if buyingAssetType == "" {
		buying.SetNull()
	} else {
		temp, err = strconv.ParseUint(buyingAssetType, 10, 64)
		if err != nil {
			return
		}
		buying.Type = uint8(temp)
		if buying.Type == types.ASSET_TYPE_CREDIT_ALPHANUM4 || buying.Type == types.ASSET_TYPE_CREDIT_ALPHANUM12 {
			buying.Code = ctx.Query("buying_asset_code")
			buying.Issuer = ethcmn.HexToAddress(ctx.Query("buying_asset_issuer"))
		}
	}

	if err = buying.Valid(); err != nil {
		return
	}
	if err = selling.Valid(); err != nil {
		return
	}

	return
}

func parseQueryBase(ctx *gin.Context) (qb types.QueryBase, err error) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	begin := ctx.Query("time_begin")
	end := ctx.Query("time_end")

	qb.Order = order
	if cursor != "" {
		qb.Cursor, err = strconv.ParseUint(cursor, 10, 64)
		if err != nil {
			return
		}
	}
	if limit != "" {
		qb.Limit, err = strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return
		}
	}

	if begin == "" {
		qb.Begin = 0
	} else {
		if qb.Begin, err = strconv.ParseUint(begin, 10, 0); err != nil {
			return
		}
	}
	if end == "" {
		qb.End = 0
	} else {
		if qb.End, err = strconv.ParseUint(end, 10, 0); err != nil {
			return
		}
	}

	return
}

func unmarshalResponseBytes(responseBytes []byte, result interface{}) (interface{}, error) {
	var err error
	response := &RPCResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling rpc response: %v", err)
	}
	errorStr := response.Error
	if errorStr != "" {
		return nil, fmt.Errorf("Response error: %v", errorStr)
	}
	// unmarshal the RawMessage into the result
	result = wire.ReadJSONPtr(result, *response.Result, &err)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling rpc response result: %v", err)
	}
	return result, nil
}
