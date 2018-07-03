package handlers

import (
	"strconv"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"

	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryPayments(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	typei := ctx.Query("type_i")
	hd.queryPayments(ctx, typei, "", "", cursor, limit, order)
}

func (hd *Handler) QueryAccPayments(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	typei := ctx.Query("type_i")
	account := ctx.Param("address")

	if account == "" {
		responseWrite(ctx, false, "param account is required")
		return
	}

	hd.queryPayments(ctx, typei, account, "", cursor, limit, order)
}

func (hd *Handler) QueryTxPayments(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	typei := ctx.Query("type_i")
	txhash := ctx.Param("txhash")

	if txhash == "" {
		responseWrite(ctx, false, "param txhash is required")
		return
	}

	hd.queryPayments(ctx, typei, "", txhash, cursor, limit, order)
}

func (hd *Handler) queryPayments(ctx *gin.Context, typei, account, txhash, cursor, limit, order string) {
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

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_PAYMENT.AppendBytes(bys)
	hd.queryAndResponse(ctx, queryData)
}
