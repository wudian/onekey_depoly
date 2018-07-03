package handlers

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	// "gitlab.zhonganinfo.com/tech_bighealth/go-sdk/ti"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/abi"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/utils"
	// "gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/config"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/encrypt"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	gin "gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) SendTransactions(ctx *gin.Context) {
	var tdata StTransactions
	if err := ctx.BindJSON(&tdata); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	if err := tdata.check(); err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	var privkeys []*ecdsa.PrivateKey
	for i := range tdata.PrivKey {
		pri := tdata.PrivKey[i]
		privkey, err := crypto.HexToECDSA(pri)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
		privkeys = append(privkeys, privkey)
	}

	accountAddress := crypto.PubkeyToAddress(privkeys[0].PublicKey)

	bigDataHashes := make([]map[string]interface{}, 0)
	ctx.Set("bigdata_hashs", &bigDataHashes)

	var ops []types.Operation
	for _, action := range tdata.Operations {
		op, err := getOperation(ctx, action)
		if err != nil {
			responseWrite(ctx, false, err.Error())
			return
		}
		ops = append(ops, op)
	}

	baseFee, ok := new(big.Int).SetString(tdata.BaseFee, 10)
	if tdata.BaseFee != "" && !ok {
		responseWrite(ctx, false, "Cannot parse basefee")
		return
	}
	nonceNum, err := strconv.ParseUint(tdata.Nonce, 10, 64)
	if err != nil {
		responseWrite(ctx, false, "Cannot parse nonceNum")
		return
	}

	tx := types.NewTransaction(nonceNum, baseFee, accountAddress,
		types.TimeBounds{Min: tdata.Lower, Max: tdata.Upper},
		ops, []byte(tdata.Memo), time.Now().UnixNano())

	//hd.SignAndCommitTx(ctx, tx, privkeys)
	hd.SignSendToCommitTx(ctx, tx, privkeys)
}

func getOperation(ctx *gin.Context, action json.RawMessage) (op types.Operation, err error) {
	type stOpType struct {
		OpType string `json:"optype"`
	}
	opt := stOpType{}
	err = json.Unmarshal(action, &opt)
	if err != nil {
		return
	}

	switch opt.OpType {
	case types.OP_CREATE_ACCOUNT.String():
		op.Type = types.OP_CREATE_ACCOUNT.ToUint()
		cadata := StCreateAccount{}
		err = json.Unmarshal(action, &cadata)
		if err != nil {
			return
		}
		if !ethcmn.IsHexAddress(cadata.Destination) {
			return op, errors.New("destination not right")
		}
		caop := types.CreateAccountOp{}
		tempsour := ethcmn.HexToAddress(cadata.Source)
		caop.Source = &tempsour
		caop.TargetAddress = ethcmn.HexToAddress(cadata.Destination)

		if stb, succ := new(big.Int).SetString(cadata.StartingBalance, 10); succ {
			caop.StartBalance = stb
		} else {
			return op, errors.New("starting balance error")
		}
		op.BodySer = caop.Bytes()
	case types.OP_PAYMENT.String():
		op.Type = types.OP_PAYMENT.ToUint()
		paydata := StPayment{}
		err = json.Unmarshal(action, &paydata)
		if err != nil {
			return
		}
		payop := types.PaymentOp{}
		tempsour := ethcmn.HexToAddress(paydata.Source)
		payop.Source = &tempsour
		tempdest := ethcmn.HexToAddress(paydata.Destination)
		payop.Destination = &tempdest
		payop.Asset.Code = paydata.AssetInfo.Code
		payop.Asset.Type = paydata.AssetInfo.Type
		payop.Asset.Issuer = ethcmn.HexToAddress(paydata.AssetInfo.Issuer)
		if err := payop.Asset.Valid(); err != nil {
			return op, err
		}
		iss := ethcmn.HexToAddress(paydata.AssetInfo.Issuer)
		payop.Issuer = &iss

		if strings.Contains(paydata.Amount, "-") {
			return op, errors.New("Negative illegal")
		}

		var ok bool
		payop.Amount, ok = new(big.Int).SetString(paydata.Amount, 10)
		if !ok {
			return op, errors.New("pay amount invalid")
		}
		op.BodySer = payop.Bytes()
	case types.OP_ALLOW_TRUST.String():
		op.Type = types.OP_ALLOW_TRUST.ToUint()
		atdata := StAllowTrust{}
		err = json.Unmarshal(action, &atdata)
		if err != nil {
			return
		}
		atop := types.AllowTrustOp{}
		tempsour := ethcmn.HexToAddress(atdata.Source)
		atop.Source = &tempsour
		atop.Trustor = ethcmn.HexToAddress(atdata.Trustor)
		atop.Asset = types.Asset{
			Code:   atdata.AssetType.Code,
			Type:   atdata.AssetType.Type,
			Issuer: ethcmn.HexToAddress(atdata.AssetType.Issuer),
		}
		if err := atop.Asset.Valid(); err != nil {
			return op, err
		}
		atop.Authorize = atdata.Authorize
		op.BodySer = atop.Bytes()
	case types.OP_CHANGE_TRUST.String():
		op.Type = types.OP_CHANGE_TRUST.ToUint()
		ctdata := StChangeTrust{}
		err = json.Unmarshal(action, &ctdata)
		if err != nil {
			return
		}
		ctop := types.ChangeTrustOp{}
		tempsour := ethcmn.HexToAddress(ctdata.Source)
		ctop.Source = &tempsour
		ctop.Line.Code = ctdata.Line.Code
		ctop.Line.Type = ctdata.Line.Type
		ctop.Line.Issuer = ethcmn.HexToAddress(ctdata.Line.Issuer)
		if err := ctop.Line.Valid(); err != nil {
			return op, err
		}
		if strings.Contains(ctdata.Limit, "-") {
			return op, errors.New("Negative illegal")
		}

		if ctdata.Limit == "" {
			ctop.Limit = big.NewInt(math.MaxInt64)
		} else {
			ctop.Limit, _ = new(big.Int).SetString(ctdata.Limit, 10)
		}
		op.BodySer = ctop.Bytes()
	case types.OP_MANAGE_DATA.String():
		op.Type = types.OP_MANAGE_DATA.ToUint()
		mandata := StManageData{}
		err = json.Unmarshal(action, &mandata)
		if err != nil {
			return
		}
		manop := types.ManageDataOp{}
		tempsour := ethcmn.HexToAddress(mandata.Source)
		manop.Source = &tempsour

		names := make([]string, len(mandata.Keypair))
		values := make([]string, len(mandata.Keypair))
		isPubs := make([]bool, len(mandata.Keypair))

		var value string
		for idx, pair := range mandata.Keypair {
			//if len(pair.Name) > types.AccDataLength || len(pair.Name) == 0 || len(pair.Value) > types.AccDataLength || len(pair.Value) == 0 || len(pair.IsPub) == types.AccDataLength || len(pair.IsPub) == 0 {
			if len(pair.Name) > types.AccDataLength || len(pair.Name) == 0 {
				err = fmt.Errorf("name is not right")
				return
			}

			if len(pair.Value) > types.AccDataLength {

				err = fmt.Errorf("value is not right")

			} else if len(pair.Value) == 0 {

				pair.Value = ""
				pair.IsPub = true
			}

			if pair.IsPub {
				value = pair.Value
			} else {
				valuesEnc := []byte(pair.Value)
				valuesEncrypt, err := encrypt.EncodeData(mandata.Source, valuesEnc)
				if err != nil {
					fmt.Println("error:", err)
				}
				value = ethcmn.Bytes2Hex(valuesEncrypt)
				//fmt.Println(len(mandata.Source), mandata.Source, valuesEncrypt)
			}

			names[idx] = pair.Name
			values[idx] = value
			isPubs[idx] = pair.IsPub
		}

		manop.Data = values
		manop.DataName = names
		manop.IsPub = isPubs

		op.BodySer = manop.Bytes()
	case types.OP_SET_OPTIONS.String():
		op.Type = types.OP_SET_OPTIONS.ToUint()
		soop := types.SetOptionsOp{}

		var mData map[string]interface{}
		err = json.Unmarshal(action, &mData)
		if err != nil {
			return
		}

		source, _ := getStr(mData, "source")
		soop.Source = ethcmn.HexToAddressPointer(source)

		if inflationDest, exist := getStr(mData, "inflationDest"); exist {
			soop.InflationDest = ethcmn.HexToAddress(inflationDest)
		}
		soop.SetFlags, soop.SetFlagsValid = getUint8(mData, "setFlags")
		soop.MasterWeight, soop.MasterWeightValid = getUint8(mData, "masterWeight")
		soop.LowThreshold, soop.LowThresholdValid = getUint8(mData, "lowThreshold")
		soop.MedThreshold, soop.MedThresholdValid = getUint8(mData, "medThreshold")
		soop.HighThreshold, soop.HighThresholdValid = getUint8(mData, "highThreshold")

		if signerItr, exist := mData["signer"]; exist {
			if signer, ok := signerItr.(map[string]interface{}); ok {
				if signerAccount, exist := getStr(signer, "signerAccount"); exist {
					soop.Signer.AccountID = ethcmn.HexToAddress(signerAccount)
				}
				if typeV, exist := getStr(signer, "type"); exist {
					soop.Signer.Type = typeV
				}
				if weight, exist := getUint8(signer, "weight"); exist {
					soop.Signer.Weight = weight
				}
			}
		}

		op.BodySer = soop.Bytes()
	case types.OP_MANAGE_OFFER.String():
		op.Type = types.OP_MANAGE_OFFER.ToUint()
		modata := StOffer{}
		err = json.Unmarshal(action, &modata)
		if err != nil {
			return
		}
		moop := types.ManageOfferOp{}
		tempsour := ethcmn.HexToAddress(modata.Source)
		moop.Source = &tempsour
		moop.Selling.Type = modata.Selling.Type
		moop.Selling.Code = modata.Selling.Code
		moop.Selling.Issuer = ethcmn.HexToAddress(modata.Selling.Issuer)
		moop.Buying.Type = modata.Buying.Type
		moop.Buying.Code = modata.Buying.Code
		moop.Buying.Issuer = ethcmn.HexToAddress(modata.Buying.Issuer)

		if moop.Buying.Code == moop.Selling.Code {
			return op, errors.New("It is not allowed to hang the same asset.")
		}

		var ok bool
		moop.Amount, ok = new(big.Int).SetString(modata.Amount, 10)
		if modata.Amount != "" && !ok {
			return op, errors.New("can't transfer to big int")
		}
		moop.OfferID = modata.OfferID
		moop.Price.V = modata.Price
		moop.Flag = uint8(types.OfferTypeActive)
		op.BodySer = moop.Bytes()
	case types.OP_CREATE_PASSIVE_OFFER.String():
		op.Type = types.OP_CREATE_PASSIVE_OFFER.ToUint()
		modata := StOffer{}
		err = json.Unmarshal(action, &modata)
		if err != nil {
			return
		}
		moop := types.ManageOfferOp{}
		tempsour := ethcmn.HexToAddress(modata.Source)
		moop.Source = &tempsour
		moop.Selling.Type = modata.Selling.Type
		moop.Selling.Code = modata.Selling.Code
		moop.Selling.Issuer = ethcmn.HexToAddress(modata.Selling.Issuer)
		moop.Buying.Type = modata.Buying.Type
		moop.Buying.Code = modata.Buying.Code
		moop.Buying.Issuer = ethcmn.HexToAddress(modata.Buying.Issuer)
		var ok bool
		moop.Amount, ok = new(big.Int).SetString(modata.Amount, 10)
		if modata.Amount != "" && !ok {
			return op, errors.New("can't transfer to big int")
		}
		moop.OfferID = modata.OfferID
		moop.Price.V = modata.Price
		moop.Flag = uint8(types.OfferTypePassive)
		op.BodySer = moop.Bytes()
	case types.OP_PATH_PAYMENT.String():
		op.Type = types.OP_PATH_PAYMENT.ToUint()
		ppdata := StPathPayment{}
		err = json.Unmarshal(action, &ppdata)
		if err != nil {
			return
		}
		ppop := types.PathPaymentOp{}
		tempsour := ethcmn.HexToAddress(ppdata.Source)
		ppop.Source = &tempsour
		dest := ethcmn.HexToAddress(ppdata.Destination)
		ppop.Destination = &dest
		ppop.SendAsset.Code = ppdata.SendAsset.Code
		ppop.SendAsset.Type = ppdata.SendAsset.Type
		ppop.SendAsset.Issuer = ethcmn.HexToAddress(ppdata.SendAsset.Issuer)
		if err := ppop.SendAsset.Valid(); err != nil {
			return op, err
		}
		ppop.DestAsset.Code = ppdata.DestAsset.Code
		ppop.DestAsset.Type = ppdata.DestAsset.Type
		ppop.DestAsset.Issuer = ethcmn.HexToAddress(ppdata.DestAsset.Issuer)
		if err := ppop.DestAsset.Valid(); err != nil {
			return op, err
		}
		var ok bool
		ppop.SendMax, ok = new(big.Int).SetString(ppdata.SendMax, 10)
		if !ok {
			return op, errors.New("can't transfer send_max to big int")
		}
		ppop.DestAmount, ok = new(big.Int).SetString(ppdata.DestAmount, 10)
		if !ok {
			return op, errors.New("can't transfer dest_amount to big int")
		}
		for _, a := range ppdata.Path {
			newAsset := types.Asset{
				Type:   a.Type,
				Code:   a.Code,
				Issuer: ethcmn.HexToAddress(a.Issuer),
			}
			if err := newAsset.Valid(); err != nil {
				return op, err
			}
			ppop.Path = append(ppop.Path, newAsset)
		}
		op.BodySer = ppop.Bytes()
	case types.OP_MANAGE_BIG_DATA.String():
		op.Type = types.OP_MANAGE_BIG_DATA.ToUint()
		predata := StBigData{}
		err = json.Unmarshal(action, &predata)
		if err != nil {
			return
		}
		// ticlient := ti.NewTiCapsuleClient(config.TiConnEndpoint(), config.TiConnKey(), config.TiConnSecret())
		// var upres ti.UploadResult
		// upres, err = ticlient.SaveData("uploadfile", []byte(predata.Value))
		// if err != nil {
		// 	return
		// }
		// if !upres.IsSuccess {
		// 	err = errors.New(upres.Info)
		// 	return
		// }
		mbop := types.ManageBigDataOp{}
		tempsour := ethcmn.HexToAddress(predata.Source)
		mbop.Source = &tempsour
		// if predata.DataType != types.DataType_File &&
		// 	predata.DataType != types.DataType_JSON &&
		// 	predata.DataType != types.DataType_Text {
		// 	err = errors.New("unknow file type")
		// 	return
		// }
		// if predata.DataType == types.DataType_File && predata.Ext == "" {
		// 	err = errors.New("data is FILE, EXT cannot empty")
		// 	return
		// }
		// mbop.Hash = upres.Hash
		// mbop.DataType = ""
		// mbop.Ext = ""
		if predata.IsPub {
			mbop.Value = predata.Value
		} else {
			value, err2 := encrypt.EncodeData(predata.Source, []byte(predata.Value))
			if err2 != nil {
				err = err2
				return
			}
			mbop.Value = ethcmn.Bytes2Hex(value)
		}
		mbop.Hash = encrypt.RlpHash(mbop.Value)
		mbop.Memo = predata.Memo
		mbop.IsPub = predata.IsPub
		op.BodySer = mbop.Bytes()
		bigDataHashes := ctx.MustGet("bigdata_hashs").(*[]map[string]interface{})
		*bigDataHashes = append(*bigDataHashes, map[string]interface{}{
			"is_pub": mbop.IsPub,
			"hash":   mbop.Hash,
			"value":  mbop.Value,
		})
	case types.OP_INFLATION.String():
		op.Type = types.OP_INFLATION.ToUint()
		var infParams StInflation
		err = json.Unmarshal(action, &infParams)
		if err != nil {
			return
		}
		infop := &types.InflationOp{}
		infop.Source = ethcmn.HexToAddressPointer(infParams.Source)
		op.BodySer = infop.Bytes()
	case types.OP_EXCUTE_CONTRACT.String():
		var (
			mExcuteContract    StExcuteContract
			jAbi               abi.ABI
			args               []interface{}
			packData, jsParams []byte
		)
		op.Type = types.OP_EXCUTE_CONTRACT.ToUint()

		if err = json.Unmarshal(action, &mExcuteContract); err != nil {
			return
		}

		//mExcuteContract.Abi = Abi

		if jAbi, err = abi.JSON(strings.NewReader(mExcuteContract.Abi)); err != nil {
			return
		}
		if len(mExcuteContract.Params) > 0 {
			if args, err = utils.ParseArgs(mExcuteContract.Func, jAbi, mExcuteContract.Params); err != nil {
				return
			}
			jsParams, _ = json.Marshal(mExcuteContract.Params)
		}
		if packData, err = jAbi.Pack(mExcuteContract.Func, args...); err != nil {
			return
		}

		gasLimit, err := strconv.ParseUint(mExcuteContract.GasLimit, 10, 64)
		if gasLimit > 8000000 || err != nil {
			return op, errors.New("Exceeds block gas limit(8000000)")
		}

		tempsour := ethcmn.HexToAddress(mExcuteContract.Source)
		opcontract := &types.ExcuteContractOp{}
		opcontract.Source = &tempsour
		opcontract.ContractAddr = ethcmn.HexToAddressPointer(mExcuteContract.ContractAddr)
		opcontract.FuncData = packData
		opcontract.Amount = mExcuteContract.Amount
		opcontract.GasLimit = mExcuteContract.GasLimit
		opcontract.Price = mExcuteContract.Price
		opcontract.FuncName = mExcuteContract.Func
		opcontract.Params = jsParams
		opcontract.Abi = mExcuteContract.Abi

		op.BodySer = opcontract.Bytes()

	case types.OP_CREATE_CONTRACT.String():
		var (
			mCreateContract    StCreateContract
			jAbi               abi.ABI
			args               []interface{}
			packData, jsParams []byte
		)
		op.Type = types.OP_CREATE_CONTRACT.ToUint()

		if err = json.Unmarshal(action, &mCreateContract); err != nil {
			return
		}

		if jAbi, err = abi.JSON(strings.NewReader(mCreateContract.Abi)); err != nil {
			return
		}
		bytecode := ethcmn.Hex2Bytes(mCreateContract.ByteCode)
		if len(bytecode) == 0 {
			err = fmt.Errorf("bytecode is null")
			return
		}
		if len(mCreateContract.Params) > 0 {
			if args, err = utils.ParseArgs("", jAbi, mCreateContract.Params); err != nil {
				return
			}
			if packData, err = jAbi.Pack("", args...); err != nil {
				return
			}
			bytecode = append(bytecode, packData...)

			jsParams, _ = json.Marshal(mCreateContract.Params)
		}

		gasLimit, err := strconv.ParseUint(mCreateContract.GasLimit, 10, 64)
		if gasLimit > 8000000 || err != nil {
			return op, errors.New("Exceeds block gas limit(8000000)")
		}

		tempsour := ethcmn.HexToAddress(mCreateContract.Source)
		opcontract := &types.CreateContractOp{}
		opcontract.Source = &tempsour
		opcontract.ByteCode = bytecode
		opcontract.GasLimit = mCreateContract.GasLimit
		opcontract.Price = mCreateContract.Price
		opcontract.Params = jsParams
		opcontract.Amount = mCreateContract.Amount

		op.BodySer = opcontract.Bytes()
	}
	return
}

func getStr(m map[string]interface{}, key string) (string, bool) {
	if v, exist := m[key]; exist {
		sv, ok := v.(string)
		return sv, ok
	}
	return "", false
}

func getUint8(m map[string]interface{}, key string) (uint8, bool) {
	if v, exist := m[key]; exist {
		uv, ok := v.(float64)
		return uint8(uv), ok
	}
	return 0, false
}
