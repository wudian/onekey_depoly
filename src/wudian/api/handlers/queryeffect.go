package handlers

import (
	"strconv"

	"encoding/json"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryAllEffects(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	begin := ctx.Query("time_begin")
	end := ctx.Query("time_end")
	typei := ctx.Query("type_i")
	hd.queryEffects(ctx, typei, "", "", begin, end, cursor, limit, order)
}

func (hd *Handler) QueryAccountEffects(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	begin := ctx.Query("time_begin")
	end := ctx.Query("time_end")
	typei := ctx.Query("type_i")
	account := ctx.Param("address")

	if account == "" {
		responseWrite(ctx, false, "param account is required")
		return
	}

	hd.queryEffects(ctx, typei, account, "", begin, end, cursor, limit, order)
}

func (hd *Handler) QueryTxEffects(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	begin := ctx.Query("time_begin")
	end := ctx.Query("time_end")
	typei := ctx.Query("type_i")
	txhash := ctx.Param("txhash")

	if txhash == "" {
		responseWrite(ctx, false, "param txhash is required")
		return
	}

	hd.queryEffects(ctx, typei, "", txhash, begin, end, cursor, limit, order)
}

func (hd *Handler) queryEffects(ctx *gin.Context, typei, account, txhash, begin, end, cursor, limit, order string) {
	var err error
	var query types.EffectsQuery
	if len(cursor) != 0 {
		query.Cursor, err = strconv.ParseUint(cursor, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	if len(limit) != 0 {
		templimit, err := strconv.ParseUint(limit, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
		query.Limit = templimit
	}
	query.Order = order

	if account != "" {
		query.Account = ethcmn.HexToAddress(account)
	}
	if txhash != "" {
		query.TxHash = ethcmn.HexToHash(txhash)
	}
	if typei == "" {
		query.Typei = types.TypeiUndefined
	} else {
		query.Typei, err = strconv.ParseUint(typei, 10, 64)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}

	if begin == "" {
		query.Begin = 0
	} else {
		query.Begin, err = strconv.ParseUint(begin, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	if end == "" {
		query.End = 0
	} else {
		query.End, err = strconv.ParseUint(end, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_EFFECTS.AppendBytes(bys)
	tmResult := new(at.RPCResult)
	err = hd.sendTxCall("query", []interface{}{queryData}, tmResult)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (*tmResult).(*at.ResultQuery)
	if 0 != res.Result.Code {
		responseWrite(ctx, false, string(res.Result.Log))
		return
	}
	resData := make([]map[string]interface{}, 0)
	json.Unmarshal(res.Result.Data, &resData)
	responseWrite(ctx, true, &resData)
}
