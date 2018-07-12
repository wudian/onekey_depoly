package server

import (
	"os"

	apiconf "wudian_go/midapi/config"
	"wudian_go/midapi/handlers"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type Server struct {
	handler *handlers.Handler
}

func NewServer() *Server {
	handler := handlers.NewHandler()

	return &Server{
		handler: handler,
	}
}

func (s *Server) Start() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/TestGet", s.handler.TestGet)
		// v1.POST("/transactions", s.handler.SendTransactions)
		// v1.POST("/specialop/validator/change", s.handler.ChangeValidator)

		// v1.GET("/genkey", s.handler.GenKey)
		// v1.GET("/nonce/:address", s.handler.GetNonceData)
		// v1.GET("/accounts/:address", s.handler.QueryAccount)

		// v1.GET("/check/:publickey", s.handler.HandlerCheckPub)
		// v1.GET("/convert/:priv", s.handler.HandlerGetPubByPriv)
		// v1.GET("/ledger/:height/transactions", s.handler.QueryHeightTx)

		// v1.GET("/transactions/:txhash", s.handler.QuerySingleTx)
		// v1.GET("/transactions", s.handler.QueryTxs)

		// v1.GET("/contract/:address", s.handler.QueryContractExist)
		// v1.POST("/contract/query", s.handler.QueryContract)
		// v1.GET("/receipt/:txhash", s.handler.QueryReceipt)

		// v1.GET("/accounts/:address/transactions", s.handler.QueryAccTxs)

		// v1.GET("/accounts/:address/offers", s.handler.QueryAccountOffers)
		// v1.GET("/accounts/:address/trades", s.handler.QueryAccountTrades)
		// v1.GET("/accounts/:address/managedata", s.handler.QueryAccountManagedata)
		// v1.POST("/managedata", s.handler.GetManageDataListByPrik)

		// v1.GET("/order_book", s.handler.QueryOrderbook)
		// v1.GET("/order_book/trades", s.handler.QueryTrades)

		// v1.GET("/accounts/:address/bigdata", s.handler.QueryAccountsBigdata)
		// v1.POST("/bigdata", s.handler.GetDataByHash)
		// v1.POST("/bigdata/thorizedata", s.handler.ThorizeData)
		// v1.POST("/bigdata/gethorizedata", s.handler.GetThorizeData)

		// v1.GET("/payments", s.handler.QueryPayments)
		// v1.GET("/accounts/:address/payments", s.handler.QueryAccPayments)
		// v1.GET("/transactions/:txhash/payments", s.handler.QueryTxPayments)

		// v1.GET("/actions", s.handler.QueryActions)
		// v1.GET("/accounts/:address/actions", s.handler.QueryAccActions)
		// v1.GET("/transactions/:txhash/actions", s.handler.QueryTxActions)

		// v1.GET("/effects", s.handler.QueryAllEffects)
		// v1.GET("/accounts/:address/effects", s.handler.QueryAccountEffects)
		// v1.GET("/transactions/:txhash/effects", s.handler.QueryTxEffects)

		// v1.GET("/ledgers", s.handler.QueryAllLedgers)
		// v1.GET("/ledgers/:sequence", s.handler.QuerySeqLedger)
	}

	if len(os.Args) > 1 && os.Args[1] == "version" {
		return
	}
	// s.handler.ReqServerInfo()
	router.Run(apiconf.ListenAddress())
}
