package handlers

import (
	"strconv"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryAccountManagedata(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	name := ctx.Query("name")
	isPub := ctx.Query("is_pub")
	account := ctx.Param("address")
	hd.queryManageData(ctx, account, isPub, name, cursor, limit, order)
}

func (hd *Handler) queryManageData(ctx *gin.Context, account string, isPub, name, cursor, limit, order string) {
	var err error
	var query types.ManageDataQuery
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
	query.Name = name
	query.IsPub = isPub

	if account != "" {
		query.Account = ethcmn.HexToAddress(account)
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_MANAGEDATA.AppendBytes(bys)

	hd.queryAndResponse(ctx, queryData)
}
