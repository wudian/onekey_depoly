package handlers

import (
	"fmt"
	// "bytes"
	// at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	// "gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-wire"

	// ethcmn "wudian_go/eth/common"
	// "wudian_go/eth/rlp"
	// "wudian_go/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func responseWrite(ctx *gin.Context, isSuccess bool, result interface{}) {
	ret := gin.H{
		"isSuccess": isSuccess,
	}

	// 用在api->midware的情形
	// response, _ := result.(map[string]interface{})
	// if response["codetype"] == at.CodeType_Timeout || response["codetype"] == at.CodeType_NonceTooLow {
	// 	ret["isSuccess"] = false
	// }

	if isSuccess {
		ret["result"] = result
	} else {
		ret["message"] = result
	}

	ctx.JSON(200, ret)

	fmt.Printf("===========raw request url: %s\n", ctx.Request.URL.String())
	fmt.Printf("===========raw response result: %v\n", result)
}
