package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	"gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryAccount(ctx *gin.Context) {
	addressHex := ctx.Param("address")
	if !ethcmn.IsHexAddress(addressHex) {
		responseWrite(ctx, false, fmt.Sprintf("Invalid address %s", addressHex))
		return
	}
	if strings.Index(addressHex, "0x") == 0 {
		addressHex = addressHex[2:]
	}

	query := types.API_QUERY_ACCOUNT.AppendBytes(ethcmn.Hex2Bytes(addressHex))
	res, err := hd.jsonRPC("query", query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
	} else {
		resData := make(map[string]interface{}, 0)
		json.Unmarshal(res.(*at.ResultQuery).Result.Data, &resData)
		responseWrite(ctx, true, &resData)
	}
}
