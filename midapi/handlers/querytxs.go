package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"

	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type TxDetail struct {
	ResultCode  uint   `json:"result_code"`
	ResultCodes string `json:"result_code_s"`
}

func (hd *Handler) QueryTxs(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	hd.queryTxs(ctx, "", "", "", cursor, limit, order)
}

func (hd *Handler) QueryAccTxs(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	account := ctx.Param("address")

	if account == "" {
		responseWrite(ctx, false, "param account is required")
		return
	}

	hd.queryTxs(ctx, account, "", "", cursor, limit, order)
}

func (hd *Handler) QuerySingleTx(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	txhash := ctx.Param("txhash")

	if txhash == "" {
		responseWrite(ctx, false, "param txhash is required")
		return
	}

	hd.queryTxs(ctx, "", txhash, "", cursor, limit, order)
}

//code = -1 查询出错，直接返回错误; 1 hash不存在，tx 还没执行需要重试查询; 其他值表示查到结果返回到上层
func (hd *Handler) queryTxForHash(txhash ethcmn.Hash) (code int, message string) {
	var (
		query          types.TxQuery
		bys, queryData []byte
		txDetails      []TxDetail
		err            error
		tmResult       *at.RPCResult
		res            *at.ResultQuery
	)

	query.TxHash = txhash
	if bys, err = rlp.EncodeToBytes(&query); err != nil {
		goto errDeal
	}
	queryData = types.API_QUERY_TX.AppendBytes(bys)

	tmResult = new(at.RPCResult)

	if err = hd.sendTxCall("query", []interface{}{queryData}, tmResult); err != nil {
		goto errDeal
	}

	fmt.Println(txhash.Hex())

	res = (*tmResult).(*at.ResultQuery)

	if res.Result.Code == 0 {
		if err = json.Unmarshal(res.Result.Data, &txDetails); err != nil {
			goto errDeal
		}
		fmt.Println(string(res.Result.Data), txDetails)
		return int(txDetails[0].ResultCode), txDetails[0].ResultCodes
	} else {
		return 1, res.Result.Log
	}
	return
errDeal:
	return -1, err.Error()

}
func (hd *Handler) QueryHeightTx(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	height := ctx.Param("height")

	if height == "" {
		responseWrite(ctx, false, "param height is required")
		return
	}
	hd.queryTxs(ctx, "", "", height, cursor, limit, order)
}

func (hd *Handler) queryTxs(ctx *gin.Context, account, txhash, height, cursor, limit, order string) {
	var err error
	var query types.TxQuery
	if len(cursor) != 0 {
		query.Cursor, err = strconv.ParseUint(cursor, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	if len(limit) != 0 {
		var tmplmt uint64
		tmplmt, err = strconv.ParseUint(limit, 10, 0)
		query.Limit = tmplmt
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	query.Order = order

	if account != "" {
		query.Account = ethcmn.HexToAddress(account)
	}
	if txhash != "" {
		query.TxHash = ethcmn.HexToHash(txhash)
	}
	query.Height = height

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_TX.AppendBytes(bys)
	hd.queryAndResponse(ctx, queryData)
	//return
}

func (hd *Handler) isTxExist(ctx *gin.Context, txhash string) bool {
	var err error
	var query types.TxQuery

	if txhash == "" {
		return false
	}
	query.Order = ""
	query.TxHash = ethcmn.HexToHash(txhash)

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		return false
	}
	queryData := types.API_QUERY_TX.AppendBytes(bys)

	if res, err := hd.jsonRPC("query", queryData); err != nil {
		return false
	} else {
		var ret []map[string]interface{}
		err = json.Unmarshal(res.(*at.ResultQuery).Result.Data, &ret)
		if err != nil {
			return false
		} else {
			return true
		}
	}
}
