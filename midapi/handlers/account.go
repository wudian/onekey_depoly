package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	"net/http"

	"strings"

	"github.com/btcsuite/btcd/btcec"

	// ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	// "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/encrypt"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) GenKey(ctx *gin.Context) {
	var err error
	var result struct {
		Privkey string `json:"privkey"`
		Address string `json:"address"`
	}

	// privkey, err := crypto.GenerateKey()
	result.Privkey, result.Address , err = encrypt.GenKey()
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	// var result struct {
	// 	Privkey string `json:"privkey"`
	// 	Address string `json:"address"`
	// }

	// result.Privkey = ethcmn.Bytes2Hex(crypto.FromECDSA(privkey))

	// result.Address = "0x" + ethcmn.Bytes2Hex(crypto.FromECDSAPubCompressed(&privkey.PublicKey))

	responseWrite(ctx, true, result)
}

func (hd *Handler) HandlerCheckPub(c *gin.Context) {
	var (
		err error
	)

	pub := c.Param("publickey")

	if err = encrypt.ValidPublicKey(pub); err != nil {
		goto errDeal
	}
	responseWrite(c, true, "有效")
	return
errDeal:
	HandleErrorMsg(c, "HandlerCheckPub", err.Error())
	return
}

func (hd *Handler) HandlerGetPubByPriv(c *gin.Context) {
	var (
		err    error
		pub    string
		priv   string
		result *types.RspConvert
	)

	result = new(types.RspConvert)
	result.ResultInfo = new(types.Result)

	priv = c.Param("priv")

	if pub, err = encrypt.PrivToPub(priv); err != nil {
		goto errDeal
	}

	c.JSON(http.StatusOK, &types.RspConvert{
		IsSuccess: true,
		ResultInfo: &types.Result{
			Address: pub,
			Privkey: priv,
		},
	})
	return
errDeal:
	HandleErrorMsg(c, "HandlerGetPubByPriv", err.Error())
	return
}

func HandleErrorMsg(c *gin.Context, requestType string, result string) {
	msg := getReturnMessage(requestType, c.Request.RemoteAddr, result)
	responseWrite(c, false, msg)
}

func getReturnMessage(requestType, remoteAddr, message string) (logMsg string) {
	logMsg = fmt.Sprintf("type[%s] From [%s] Error [%s] ", requestType, remoteAddr, message)
	return
}

func Encrypt(origin, key []byte) ([]byte, error) {
	if len(origin)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("origin text not full blocks")
	}

	key = ZeroPadding(key, aes.BlockSize)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	crypted := make([]byte, aes.BlockSize+len(origin))
	iv := crypted[:aes.BlockSize]

	blockMode := cipher.NewCBCEncrypter(cipherBlock, iv)
	blockMode.CryptBlocks(crypted[aes.BlockSize:], origin)

	return crypted, nil
}

func Decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < aes.BlockSize {
		return nil, fmt.Errorf("crypted text too short")
	}

	key = ZeroPadding(key, aes.BlockSize)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := crypted[:aes.BlockSize]
	origin := make([]byte, len(crypted[aes.BlockSize:]))

	blockMode := cipher.NewCBCDecrypter(cipherBlock, iv)
	blockMode.CryptBlocks(origin, crypted[aes.BlockSize:])

	return origin, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext) < blockSize {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{0}, padding)
		ciphertext = append(ciphertext, padtext...)
	} else {
		ciphertext = ciphertext[:blockSize]
	}
	return ciphertext
}

func ValidPublicKey(publicKey string) (err error) {
	var (
		pubKeyBytes []byte
	)

	if pubKeyBytes, err = hex.DecodeString(strings.TrimPrefix(publicKey, "0x")); err != nil {
		return
	}
	_, err = btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	return
}
