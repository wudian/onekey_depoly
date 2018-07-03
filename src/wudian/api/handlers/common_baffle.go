package handlers

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	// "bytes"
	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) SignAndSendTxBaffle(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey) {
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

	//	ret, err := hd.jsonRPC("broadcast_tx_sync", txBytes)
	//	if err != nil {
	//		responseWrite(ctx, false, err.Error())
	//		return
	//	}
	//	res := (ret).(*at.ResultBroadcastTx)

	//fmt.Printf("resultbroadcastTx : %v", string(res.Data[:]))
	fmt.Println(txBytes)
	res := &at.ResultBroadcastTxCommit{0, []byte("test"), "test"}
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

func (hd *Handler) SignAndCommitTxBaffle(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey) {
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

	hd.recursiveTransactionsBaffle(ctx, sigTx, 3)
}

// 为了解决交易超时的bug，而让交易递归地进行。 递归的次数为n
func (hd *Handler) recursiveTransactionsBaffle(ctx *gin.Context, sigTx *types.Transaction, n int) {
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
	fmt.Println(txBytes)
	//	ret, err := hd.jsonRPC("broadcast_tx_commit", txBytes)
	//	if err != nil {
	//		responseWrite(ctx, false, err.Error())
	//		return
	//	}
	//	res := (ret).(*at.ResultBroadcastTxCommit)
	res := &at.ResultBroadcastTxCommit{0, []byte("test"), "test"}
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
		hd.recursiveTransactionsBaffle(ctx, sigTx, n-1)
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

func (hd *Handler) SignAndSendTxForBigDataBaffle(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey, hashs string) {
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

	//	rptRes := new(at.RPCResult)
	//	err = hd.sendTxCall("broadcast_tx_sync", []interface{}{txBytes}, rptRes)

	//	if err != nil {
	//		responseWrite(ctx, false, err.Error())
	//		return
	//	}
	//	res := (*rptRes).(*at.ResultBroadcastTx)
	fmt.Println(txBytes)
	res := &at.ResultBroadcastTxCommit{0, []byte("test"), "test"}
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

func (hd *Handler) SignAndCommitTxForBigDataBaffle(ctx *gin.Context, tx *types.Transaction, privkey []*ecdsa.PrivateKey, hashs string) {
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

	//	rptRes := new(at.RPCResult)
	//	err = hd.sendTxCall("broadcast_tx_commit", []interface{}{txBytes}, rptRes)

	//	if err != nil {
	//		responseWrite(ctx, false, err.Error())
	//		return
	//	}
	//	res := (*rptRes).(*at.ResultBroadcastTxCommit)

	//fmt.Printf("resultbroadcastTx : %v", string(res.Data[:]))
	fmt.Println(txBytes)
	res := &at.ResultBroadcastTxCommit{0, []byte("test"), "test"}
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
