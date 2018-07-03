package handlers

import (
	"fmt"
	"strings"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	"gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) GetNonceData(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if !ethcmn.IsHexAddress(addressHex) {
		responseWrite(ctx, false, fmt.Sprintf("Invalid address %s", addressHex))
		return
	}

	if strings.Index(addressHex, "0x") == 0 {
		addressHex = addressHex[2:]
	}

	nonce, err := hd.getNonce(ethcmn.Hex2Bytes(addressHex))
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	responseWrite(ctx, true, nonce)
}

func (hd *Handler) getNonce(address []byte) (uint64, error) {
	query := types.API_QUERY_NONCE.AppendBytes(address)
	res, err := hd.jsonRPC("query", query)
	if err != nil {
		return 0, err
	} else {
		nonce := new(uint64)
		rlp.DecodeBytes(res.(*at.ResultQuery).Result.Data, nonce)
		return *nonce, nil
	}
}
