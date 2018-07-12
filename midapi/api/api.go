package api

import (
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	OnClose()
	HandlerPost(c *gin.Context)
	HandlerGet(c *gin.Context)
	HandlerGetDataByHash(c *gin.Context)
	HandlerTransactions(c *gin.Context)
	HandlerThorizeData(c *gin.Context)
	HandlerGetThorizeData(c *gin.Context)
	HandlerGetPubByPriv(c *gin.Context)
	HandlerCheckPub(c *gin.Context)
}

type DataBase interface {
	OnClose()
	Insert(...interface{}) error
	Get(...interface{}) ([]byte, error)
	IsExist(...interface{}) (bool, error)
	Del(...interface{}) error

	HInsert(...interface{}) error
	HGet(...interface{}) ([]byte, error)
	HIsExist(...interface{}) (bool, error)
	HDel(...interface{}) error
}

type DataBig interface {
	Upload(...interface{}) error
	Download(...interface{}) error
}
