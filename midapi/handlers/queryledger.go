package handlers

import (
	"strconv"

	"encoding/json"

	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryAllLedgers(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")

	var err error
	var query types.LedgerQuery
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
	if order != "" {
		query.Order = order
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_ALL_LEDGER.AppendBytes(bys)
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
	for idx := range resData {
		resData[idx]["hash"] = "0x" + string(types.TrimZero([]byte(resData[idx]["hash"].(string))))
		resData[idx]["prev_hash"] = "0x" + string(types.TrimZero([]byte(resData[idx]["prev_hash"].(string))))
	}
	responseWrite(ctx, true, &resData)
}

func (hd *Handler) QuerySeqLedger(ctx *gin.Context) {
	sequence := ctx.Param("sequence")
	var err error
	var query types.LedgerQuery
	if len(sequence) != 0 {
		seq, err := strconv.ParseUint(sequence, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
		query.Sequence = seq
	}

	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_SEQ_LEDGER.AppendBytes(bys)
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
	resData := make(map[string]interface{}, 0)
	json.Unmarshal(res.Result.Data, &resData)
	resData["hash"] = "0x" + string(types.TrimZero([]byte(resData["hash"].(string))))
	resData["prev_hash"] = "0x" + string(types.TrimZero([]byte(resData["prev_hash"].(string))))

	responseWrite(ctx, true, &resData)

}
