package handlers

import (
	"encoding/hex"
	"encoding/json"
	"time"

	atypes "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	gcommon "gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-common"
	"gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-crypto"
	"gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-wire"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/config"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) ChangeValidator(ctx *gin.Context) {
	var tdata SpencialOp
	if err := ctx.BindJSON(&tdata); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	if tdata.ValidatorPub == "" {
		responseWrite(ctx, false, "validator_pubkey cannot be empty")
		return
	}
	validatorPub := gcommon.SanitizeHex(tdata.ValidatorPub)

	pubkeybyte32, err := atypes.StringTo32byte(validatorPub)
	if err != nil {
		responseWrite(ctx, false, err)
		return
	}
	vb := wire.JSONBytes(atypes.ValidatorAttr{
		PubKey:     crypto.PubKeyEd25519(pubkeybyte32).Bytes(),
		Power:      tdata.Power,
		IsCA:       tdata.IsCA,
		RPCAddress: config.BackendCallAddress(), /*"tcp://0.0.0.0:46657"*/
	})

	cmd := &atypes.SpecialOPCmd{
		CmdCode: atypes.SpecialOP,
		CmdType: atypes.SpecialOP_ChangeValidator,
		Msg:     vb,
		Time:    time.Now(),
		Nonce:   0, // TODO
	}
	for _, sig := range tdata.Sigs {
		sigb, err := hex.DecodeString(sig)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
		cmd.Sigs = append(cmd.Sigs, sigb)
	}

	rawbytes, _ := json.Marshal(cmd)
	cmdbytes := atypes.TagSpecialOPTx(rawbytes)

	tmResult := new(atypes.RPCResult)
	err = hd.sendTxCall("request_special_op", []interface{}{cmdbytes}, tmResult)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	//res := (*tmResult).(*atypes.ResultRequestSpecialOP)
	//response := map[string]interface{}{
	//"data": string(res.Data),
	//"log":  res.Log,
	//}
	responseWrite(ctx, true, tmResult)
}
