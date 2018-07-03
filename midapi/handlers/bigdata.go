package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	// "gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/config"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	// "gitlab.zhonganinfo.com/tech_bighealth/go-sdk/ti"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/encrypt"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryAccountsBigdata(ctx *gin.Context) {
	addr := ctx.Param("address")
	cursor := ctx.Query("cursor")
	limit := ctx.Query("limit")
	order := ctx.Query("order")
	memo := ctx.Query("memo")
	hd.queryAccountsBigdata(ctx, addr, memo, cursor, limit, order)
}

func (hd *Handler) queryAccountsBigdata(ctx *gin.Context, addr, memo, cursor, limit, order string) {
	var err error
	var query types.BigdataQuery
	if len(cursor) != 0 {
		query.Cursor, err = strconv.ParseUint(cursor, 10, 0)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	if len(limit) != 0 {
		var tmplmt uint64
		tmplmt, err = strconv.ParseUint(limit, 10, 0)
		query.Limit = tmplmt
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}
	query.Order = order
	query.Account = ethcmn.HexToAddress(addr)
	query.Memo = memo
	var bys []byte
	bys, err = rlp.EncodeToBytes(&query)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	queryData := types.API_QUERY_ALL_BIGDATA.AppendBytes(bys)
	hd.queryAndResponse(ctx, queryData)
}

func (hd *Handler) getSingleData(Source, Hash string) (map[string]interface{}, error) {
	var query types.SingleData
	query.Account = ethcmn.HexToAddress(Source)
	query.Hash = Hash
	var bys []byte
	bys, err := rlp.EncodeToBytes(&query)
	if err != nil {
		return nil, err
	}
	queryData := types.API_QUERY_SINGLE_BIGDATA.AppendBytes(bys)
	res, err := hd.jsonRPC("query", queryData)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}
	err = json.Unmarshal(res.(*at.ResultQuery).Result.Data, &ret)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, errors.New("Data not found")
	}
	return ret, nil
}

func (hd *Handler) GetDataByHash(ctx *gin.Context) {
	var (
		tdata       StGetBigData
		err         error
		pri_is_null bool
	)
	if err := ctx.BindJSON(&tdata); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	pri_is_null = (tdata.Privkey == "")
	if !pri_is_null {
		tdata.Source, err = encrypt.PrivToPub(tdata.Privkey)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
	}

	var ret []map[string]interface{}

	for _, hash := range tdata.Hashs {
		// modify by zhaoyang.
		res, err := hd.getSingleData(tdata.Source, hash)
		if err != nil {
			//ret = append(ret, map[string]interface{}{"err": err.Error()})

			fmt.Println("GetDataByHash hash:", hash, "error", err.Error())
			responseWrite(ctx, false, err.Error())
			return
		}

		if len(res) == 0 {
			continue
		}

		is_pub := res["is_pub"].(bool)
		if !is_pub && pri_is_null {
			//	ret = append(ret, map[string]interface{}{"err": errors.New("privkey is null but is_pub is false")})
			continue
		}

		if !is_pub {
			byValue := ethcmn.Hex2Bytes(res["value"].(string))
			deValue, err := encrypt.DecodeData(tdata.Privkey, byValue)
			if err != nil {
				ret = append(ret, map[string]interface{}{"err": err.Error()})
				continue
			} else {
				res["value"] = string(deValue)
			}
		}

		ret = append(ret, res)

	}

	responseWrite(ctx, true, ret)
}

func (hd *Handler) ThorizeData(c *gin.Context) {
	var (
		tdata  ReqThorizeData
		err    error
		Source string
	)

	if err = c.BindJSON(&tdata); err != nil {
		responseWrite(c, false, err.Error())
		return
	}

	if Source, err = encrypt.PrivToPub(tdata.Privkey); err != nil {
		responseWrite(c, false, err.Error())
		return
	}

	for _, hash := range tdata.Hashs {
		res, err := hd.getSingleData(Source, hash)
		if err != nil {
			responseWrite(c, false, err.Error())
			return
		}

		value := res["value"].(string)
		is_pub := res["is_pub"].(bool)
		if !is_pub {
			byValue := ethcmn.Hex2Bytes(value)
			deValue, err := encrypt.DecodeData(tdata.Privkey, byValue)
			if err != nil {
				responseWrite(c, false, err.Error())
				return
			}
			value = string(deValue)
		}

		newValue, err := encrypt.EncodeData(tdata.Destination, []byte(value))
		if err != nil {
			responseWrite(c, false, err.Error())
			return
		}
		res["value"] = ethcmn.Bytes2Hex(newValue)

		js, _ := json.Marshal(res)
		err = hd.redisApi.Insert(hash, js, 600)
		if err != nil {
			responseWrite(c, false, err.Error())
			return
		}
	}

	//  h.putRedis(id, value)
	responseWrite(c, true, nil)
}

func (hd *Handler) GetThorizeData(c *gin.Context) {
	var (
		tdata ReqGetThorizeData
		err   error
		ret   []map[string]interface{}
	)

	if err = c.BindJSON(&tdata); err != nil {
		responseWrite(c, false, err.Error())
		return
	}

	for _, hash := range tdata.Hashs {
		res, err := hd.redisApi.Get(hash)
		if err != nil {
			ret = append(ret, map[string]interface{}{"err": err.Error()})
		} else {
			js := map[string]interface{}{}
			err = json.Unmarshal(res, &js)
			if err != nil {
				ret = append(ret, map[string]interface{}{"err": err.Error()})
			} else {
				value := js["value"].(string)
				byValue := ethcmn.Hex2Bytes(value)
				deValue, err := encrypt.DecodeData(tdata.Privkey, byValue)
				if err != nil {
					responseWrite(c, false, err.Error())
					return
				}
				js["value"] = string(deValue)
				delete(js, "created_at")
				delete(js, "js_pub")
				ret = append(ret, js)
			}
		}
	}

	responseWrite(c, true, ret)
}
