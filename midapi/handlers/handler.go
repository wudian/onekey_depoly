package handlers

import (
	"encoding/json"
	"errors"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	rpc "gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-rpc/client"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/api"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/config"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/utils/redis"
	"go.uber.org/zap"
	gin "gopkg.in/gin-gonic/gin.v1"
)

var logger *zap.Logger

func InitLogger(log *zap.Logger) {
	logger = log
}

type Handler struct {
	client   *rpc.ClientJSONRPC
	redisApi api.DataBase
	expire   int
}

func NewHandler() *Handler {
	var h Handler
	h.client = rpc.NewClientJSONRPC(logger, config.BackendCallAddress())
	h.redisApi = redis.NewClient(config.Redis(), int(config.Redisidls()), int(config.Redistimeout()))
	h.expire = int(config.Expire())
	return &h
}

const (
	Redis_Key string = "midapi:"
	// Redis_Txdata string = "txdata"
	// Redis_Nonce  string = "nonce"
	Redis_Result string = "result"
)

func (hd *Handler) jsonRPC(action string, p []byte) (interface{}, error) {
	tmResult := new(at.RPCResult)
	var err error
	if p != nil {
		_, err = hd.client.Call(action, []interface{}{p}, tmResult)
	} else {
		_, err = hd.client.Call(action, []interface{}{}, tmResult)
	}

	if err != nil {
		return nil, err
	}

	if action == "query" {
		res := (*tmResult).(*at.ResultQuery)
		if 0 != res.Result.Code {
			return nil, errors.New(res.Result.Log)
		}
		return res, nil
	} else {
		return *tmResult, nil
	}
}

func (hd *Handler) queryAndResponse(ctx *gin.Context, queryData []byte) {
	res, err := hd.jsonRPC("query", queryData)
	if err != nil {
		responseWrite(ctx, false, err.Error())
	} else {
		var ret []map[string]interface{}
		err = json.Unmarshal(res.(*at.ResultQuery).Result.Data, &ret)
		if err != nil {
			responseWrite(ctx, false, err.Error())
		} else {
			responseWrite(ctx, true, ret)
		}
	}
}

func (hd *Handler) sendTxCall(method string, params []interface{}, result interface{}) error {
	_, err := hd.client.Call(method, params, result)
	if err != nil {
		return err
	}

	return nil
}
