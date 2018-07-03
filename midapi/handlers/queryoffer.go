package handlers

import (
	"encoding/json"
	"strconv"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

//QueryAccountOffers query offers belonging to a specific account
func (hd *Handler) QueryAccountOffers(ctx *gin.Context) {
	var err error
	var query types.AccountOfferQuery

	query.Account = ethcmn.HexToAddress(ctx.Param("address"))

	if query.Selling, query.Buying, err = parseAsset(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	if query.QueryBase, err = parseQueryBase(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_OFFER_AC.AppendBytes(bys)

	hd.queryAndResponse(ctx, queryData)
}

//QueryOrderbook query orderbook of specific assert
func (hd *Handler) QueryOrderbook(ctx *gin.Context) {
	var err error
	var query types.OrderbookQuery

	if query.Selling, query.Buying, err = parseAsset(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	if query.Selling.IsNull() || query.Buying.IsNull() {
		responseWrite(ctx, false, "Invalid asset parameter")
		return
	}

	if amt := ctx.Query("amount"); amt == "" {
		query.Amount = 20
	} else {
		query.Amount, err = strconv.ParseUint(amt, 10, 64)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}

	if query.Amount > 100 {
		query.Amount = 100
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_ORDER_BOOK.AppendBytes(bys)

	// hd.queryAndResponse(ctx, queryData)
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

	var ret map[string]interface{}
	err = json.Unmarshal(res.Result.Data, &ret)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	responseWrite(ctx, true, ret)
}

//QueryTrades query trades of specific assert
func (hd *Handler) QueryTrades(ctx *gin.Context) {
	var err error
	var query types.TradesQuery

	if query.Selling, query.Buying, err = parseAsset(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	if query.QueryBase, err = parseQueryBase(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_TRADES.AppendBytes(bys)

	hd.queryAndResponse(ctx, queryData)
}

//QueryAccountTrades query trades for given account
func (hd *Handler) QueryAccountTrades(ctx *gin.Context) {
	var err error
	offerid := ctx.Query("offerid")
	var query types.AccountTradesQuery

	query.Account = ethcmn.HexToAddress(ctx.Param("address"))
	if len(offerid) != 0 {
		query.OfferID, err = strconv.ParseUint(offerid, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}

	if query.Selling, query.Buying, err = parseAsset(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	if query.QueryBase, err = parseQueryBase(ctx); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_TRADES_AC.AppendBytes(bys)

	hd.queryAndResponse(ctx, queryData)
}
