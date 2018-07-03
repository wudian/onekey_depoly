package handlers

import (
	"strconv"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"

	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryActions(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	begin := ctx.Query("time_begin")
	end := ctx.Query("time_end")
	typei := ctx.Query("type_i")
	hd.queryActions(ctx, typei, "", "", begin, end, cursor, limit, order)
}

func (hd *Handler) QueryAccActions(ctx *gin.Context) {
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

//	if err := ValidPublicKey(account); err != nil {
//		responseWrite(ctx, false, "account is not right")
//		return
//	}

	hd.queryActions(ctx, typei, account, "", begin, end, cursor, limit, order)
}

func (hd *Handler) QueryTxActions(ctx *gin.Context) {
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

	hd.queryActions(ctx, typei, "", txhash, begin, end, cursor, limit, order)
}

func (hd *Handler) queryActions(ctx *gin.Context, typei, account, txhash, begin, end, cursor, limit, order string) {
	var err error
	var query types.ActionsQuery
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
	queryData := types.API_QUERY_ACTION.AppendBytes(bys)
	hd.queryAndResponse(ctx, queryData)
}
