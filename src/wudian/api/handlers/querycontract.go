package handlers

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/abi"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/utils"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
	"gopkg.in/gin-gonic/gin.v1"
)

func (hd *Handler) QueryContractExist(ctx *gin.Context) {

	addressHex := ctx.Param("address")
	if !ethcmn.IsHexAddress(addressHex) {
		responseWrite(ctx, false, fmt.Sprintf("Invalid address %s", addressHex))
		return
	}
	if strings.Index(addressHex, "0x") == 0 {
		addressHex = addressHex[2:]
	}
	query := types.API_QUERY_CONTRACT_EXIST.AppendBytes(ethcmn.Hex2Bytes(addressHex))
	tmResult := new(at.RPCResult)
	err := hd.sendTxCall("query", []interface{}{query}, tmResult)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}

	res := (*tmResult).(*at.ResultQuery)
	if 0 != res.Result.Code {
		responseWrite(ctx, false, string(res.Result.Log))
		return
	}
	resData := make(map[string]interface{}, 0)
	json.Unmarshal(res.Result.Data, &resData)
	responseWrite(ctx, true, &resData)
}

func (hd *Handler) QueryReceipt(ctx *gin.Context) {

	txHash := ctx.Param("txhash")

	btyHash := ethcmn.HexToHash(txHash)

	query := types.API_QUERY_RECEIPT.AppendBytes(btyHash.Bytes())

	tmResult := new(at.RPCResult)

	err := hd.sendTxCall("query", []interface{}{query}, tmResult)
	if err != nil {
		responseWrite(ctx, false, err.Error())
		return
	}
	res := (*tmResult).(*at.ResultQuery)

	if 0 != res.Result.Code {
		responseWrite(ctx, false, string(res.Result.Log))
		return
	}

	var receipt types.Receipt

	if err := rlp.DecodeBytes(res.Result.Data, &receipt); err != nil {
		responseWrite(ctx, false, err)
		return
	}

	stReceipt := &StReceipt{
		OpType:          receipt.OpType,
		Source:          receipt.Source.Hex(),
		TxHash:          receipt.TxHash.Hex(),
		Height:          receipt.Height,
		ContractAddress: receipt.ContractAddress,
		Function:        receipt.Function,
		GasPrice:        receipt.GasPrice,
		GasLimit:        receipt.GasLimit,
		GasUsed:         receipt.GasUsed,
		TxReceiptStatus: receipt.TxReceiptStatus,
		Message:         receipt.Message,
	}

	json.Unmarshal(receipt.Params, &stReceipt.Params)
	json.Unmarshal(receipt.Logs, &stReceipt.Logs)
	json.Unmarshal(receipt.Res, &stReceipt.Res)

	responseWrite(ctx, true, stReceipt)

	return
}

func (hd *Handler) QueryContract(ctx *gin.Context) {
	var (
		stQContract    StQueryContract
		jAbi           abi.ABI
		args           []interface{}
		parseResult    interface{}
		packData       []byte
		query, queryTx []byte
		qContract      types.QueryContract
		tmResult       *at.RPCResult
		res            *at.ResultQuery
		tx             *types.Transaction
		err            error
		op             types.Operation
		privkey        *ecdsa.PrivateKey
	)
	if err = ctx.BindJSON(&stQContract); err != nil {
		goto errDeal
	}

	if jAbi, err = abi.JSON(strings.NewReader(stQContract.Abi)); err != nil {
		err = errors.New("format abi error:" + err.Error())
		goto errDeal
	}
	if !jAbi.Methods[stQContract.Func].Const {
		err = fmt.Errorf("Contract Func %v Is Not Query Method")
		goto errDeal
	}
	if len(stQContract.Params) > 0 {
		if args, err = utils.ParseArgs(stQContract.Func, jAbi, stQContract.Params); err != nil {
			goto errDeal
		}
	}
	if packData, err = jAbi.Pack(stQContract.Func, args...); err != nil {
		goto errDeal
	}
	qContract.Amount = "0"
	qContract.Price = "0"
	qContract.GasLimit = "0"
	qContract.FuncData = packData
	qContract.ContractAddr = ethcmn.HexToAddressPointer(stQContract.ContractAddr)
	qContract.Source = ethcmn.HexToAddressPointer("0x00000000000000000000000000000000")

	op.Type = types.OP_QUERY_CONTRACT.ToUint()
	op.BodySer = qContract.Bytes()

	tx = types.NewTransaction(0, big.NewInt(0), *qContract.Source,
		types.TimeBounds{Min: 0, Max: 0},
		[]types.Operation{op}, []byte(""), time.Now().UnixNano())

	if stQContract.PrivKey != "" {
		stQContract.PrivKey = formAddress(stQContract.PrivKey)

		if len(stQContract.PrivKey) != 64 {
			err = fmt.Errorf("Invalid privkey, length %d", len(stQContract.PrivKey))
			goto errDeal
		}
		privkey = crypto.ToECDSA(ethcmn.Hex2Bytes(stQContract.PrivKey))
		if tx, err = tx.Sign([]*ecdsa.PrivateKey{privkey}); err != nil {
			goto errDeal
		}
	}

	if queryTx, err = rlp.EncodeToBytes(tx); err != nil {
		goto errDeal
	}

	query = types.API_QUERY_CONTRACT.AppendBytes(queryTx)

	tmResult = new(at.RPCResult)

	if err = hd.sendTxCall("query", []interface{}{query}, tmResult); err != nil {
		goto errDeal
	}

	res = (*tmResult).(*at.ResultQuery)

	if res.Result.Code != at.CodeType_OK {
		err = errors.New(res.Result.Log)
		goto errDeal
	}

	if parseResult, err = unpackResultToArray(stQContract.Func, jAbi, res.Result.Data); err != nil {
		goto errDeal
	}
	responseWrite(ctx, true, parseResult)
	return
errDeal:
	responseWrite(ctx, false, err.Error())
	return
}

func formAddress(addr string) string {
	addr = strings.ToUpper(addr)
	return strings.TrimPrefix(addr, "0X")
}

func unpackResultToArray(method string, abiDef abi.ABI, output []byte) (interface{}, error) {
	if len(output) == 0 {
		return nil, nil
	}
	m, ok := abiDef.Methods[method]
	if !ok {
		return nil, errors.New("No such method")
	}
	if len(m.Outputs) == 0 {
		return nil, errors.New("method " + m.Name + " doesn't have any returns")
	}
	if len(m.Outputs) == 1 {
		var result interface{}
		d := ethcmn.ParseData(output)
		if err := abiDef.Unpack(&result, method, d); err != nil {
			return nil, err
		}
		return result, nil
	}
	var result []interface{}
	d := ethcmn.ParseData(output)
	if err := abiDef.Unpack(&result, method, d); err != nil {
		return nil, err
	}
	return result, nil
}

func unpackResult(method string, abiDef abi.ABI, output string) (interface{}, error) {
	m, ok := abiDef.Methods[method]
	if !ok {
		return nil, errors.New("No such method")
	}

	if len(m.Outputs) == 0 {
		return nil, errors.New("method " + m.Name + " doesn't have any returns")
	}
	if len(m.Outputs) == 1 {
		var result interface{}
		parsedData := ethcmn.ParseData(output)
		if err := abiDef.Unpack(&result, method, parsedData); err != nil {
			return nil, err
		}
		if strings.Index(m.Outputs[0].Type.String(), "bytes") == 0 {
			b := result.([]byte)
			idx := 0
			for i := 0; i < len(b); i++ {
				if b[i] == 0 {
					idx = i
				} else {
					break
				}
			}
			b = b[idx+1:]
			return fmt.Sprintf("%s", b), nil
		}
		return result, nil
	}
	d := ethcmn.ParseData(output)
	var result []interface{}
	if err := abiDef.Unpack(&result, method, d); err != nil {
		return nil, err
	}

	retVal := map[string]interface{}{}
	for i, output := range m.Outputs {
		if strings.Index(output.Type.String(), "bytes") == 0 {
			b := result[i].([]byte)
			idx := 0
			for i := 0; i < len(b); i++ {
				if b[i] == 0 {
					idx = i
				} else {
					break
				}
			}
			b = b[idx+1:]
			retVal[output.Name] = fmt.Sprintf("%s", b)
		} else {
			retVal[output.Name] = result[i]
		}
	}
	return retVal, nil
}
