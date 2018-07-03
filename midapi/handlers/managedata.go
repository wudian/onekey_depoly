package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/encrypt"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) GetManageDataListByPrik(ctx *gin.Context) {
	var (
		tdata StQueryManageData
		ret   []map[string]types.AccDataIsPub
		res   map[string]types.AccDataIsPub
		err   error
	)

	if err = ctx.BindJSON(&tdata); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	for _, key := range tdata.Keys {
		if res, err = hd.querySingleManageData(tdata.Source, key); err != nil {
			fmt.Println("querySingleManageData key:", key, "error", err.Error())
			responseWrite(ctx, false, err.Error())
			return
		}

		if len(res) == 0 {
			continue
		}

		if len(tdata.PrivKey) == 0 && !res[key].IsPub {
			continue
		}

		if !res[key].IsPub {

			if res[key].Value != "" {

				byValue := ethcmn.Hex2Bytes(res[key].Value)
				deValue, err := encrypt.DecodeData(tdata.PrivKey, byValue)
				if err != nil {
					fmt.Println("DecodeData error:", err.Error())
				}
				dData := types.AccDataIsPub{
					Value: string(deValue),
					IsPub: res[key].IsPub,
				}
				res[key] = dData
			}
		}
		ret = append(ret, res)
	}
	responseWrite(ctx, true, ret)
}

func (hd *Handler) querySingleManageData(source string, keys string) (ret map[string]types.AccDataIsPub, err error) {
	var query types.SingleManageData
	query.Account = ethcmn.HexToAddress(source)
	query.Keys = keys
	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		return
	}
	queryData := types.API_QUERY_SINGLE_MANAGEDATA.AppendBytes(bys)
	res, err := hd.jsonRPC("query", queryData)
	if err != nil {
		return
	}

	if err = json.Unmarshal(res.(*at.ResultQuery).Result.Data, &ret); err != nil {
		return
	}

	if len(ret) == 0 {
		err = errors.New("Data not found")
		return
	}
	return
}
