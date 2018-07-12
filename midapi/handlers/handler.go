package handlers

import (
	// "encoding/json"
	// "errors"

	"wudian_go/midapi/api"
	"wudian_go/midapi/config"
	"wudian_go/lib/redis"

	// at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	// rpc "gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-rpc/client"
	"go.uber.org/zap"
	gin "gopkg.in/gin-gonic/gin.v1"
)

var logger *zap.Logger

func InitLogger(log *zap.Logger) {
	logger = log
}

type Handler struct {
	// client   *rpc.ClientJSONRPC
	redisApi api.DataBase
	expire   int
}

func NewHandler() *Handler {
	var h Handler
	// h.client = rpc.NewClientJSONRPC(logger, config.BackendCallAddress())
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



func (hd *Handler) TestGet(ctx *gin.Context) {
	responseWrite(ctx, true, 123)
}
