package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	delostypes "gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
)

type StTransactions struct {
	TimeBounds
	PrivKey    []string          `json:"privkey"`
	BaseFee    string            `json:"basefee"`
	Memo       string            `json:"memo"`
	Operations []json.RawMessage `json:"operations"`
	Nonce      string            `json:"nonce"`
}

func (tdata *StTransactions) check() error {
	if len(tdata.PrivKey) == 0 {
		// responseWrite(ctx, false, "please input privkey")
		return errors.New("please input privkey")
	}
	if len(tdata.Memo) > 1024 {
		// responseWrite(ctx, false, "memo lenth error")
		return errors.New("memo lenth error")
	}

	for i := range tdata.PrivKey {
		pri := tdata.PrivKey[i]
		if strings.Index(pri, "0x") == 0 {
			pri = pri[2:]
		}
		if len(pri) != 64 {
			// responseWrite(ctx, false, fmt.Sprintf("Invalid sk, length %d. %s", len(pri), pri))
			return errors.New(fmt.Sprintf("Invalid sk, length %d. %s", len(pri), pri))
		}
		tdata.PrivKey[i] = pri
	}

	return nil
}

type Asset struct {
	Issuer string `json:"issuer"`
	Type   uint8  `json:"type"`
	Code   string `json:"code"`
}

type TimeBounds struct {
	Lower uint64 `json:"lowertime"` // rlp does not support int64, so use uint64 here
	Upper uint64 `json:"uppertime"` // rlp does not support int64, so use uint64 here
}

type OperationBase struct {
	OpType string `json:"optype"`
	Source string `json:"source"`
}

type StChangeTrust struct {
	OperationBase
	Line  Asset  `json:"asset"`
	Limit string `json:"limit"`
}

type StAllowTrust struct {
	OperationBase
	Trustor   string `json:"trustor"`
	AssetType Asset  `json:"asset"`
	Authorize bool   `json:"authorize"`
}

type StCreateAccount struct {
	OperationBase
	Destination     string `json:"destination"`
	StartingBalance string `json:"startingBalance"`
}

type StPayment struct {
	OperationBase
	Destination string `json:"destination"`
	AssetInfo   Asset  `json:"asset"`
	Amount      string `json:"amount"`
}

type StPathPayment struct {
	OperationBase
	Destination string  `json:"destination"`
	SendAsset   Asset   `json:"send_asset"`
	SendMax     string  `json:"send_max"`
	DestAsset   Asset   `json:"dest_asset"`
	DestAmount  string  `json:"dest_amount"`
	Path        []Asset `json:"path"`
}

type StInflation struct {
	OperationBase
	// Time TimeBounds `json:"time"`
}

type StBigData struct {
	OperationBase
	IsPub bool `json:"is_pub"`
	// Value    []byte `json:"value"`
	Value string `json:"value"`
	// DataType string `json:"type"`
	// Ext      string `json:"ext"`
	Memo string `json:"memo"`
}

type StOffer struct {
	OperationBase
	Selling Asset  `json:"selling"`
	Buying  Asset  `json:"buying"`
	Amount  string `json:"amount"`
	Price   string `json:"price"`
	OfferID uint64 `json:"offerid"`
}

type StGetBigData struct {
	Source  string   `json:"source"`
	Privkey string   `json:"privkey"`
	Hashs   []string `json:"hashs"`
	// Keys    []struct {
	// 	Hash  string `json:"hash"`
	// 	IsPub bool   `json:"is_pub"`
	// } `json:"keys"`
}

type StQueryManageData struct {
	Source  string   `json:"source"`
	PrivKey string   `json:"privkey"`
	Keys    []string `json:"keys"`
}

type StManageData struct {
	OperationBase
	Source  string `json:"source"`
	Privkey string `json:"privkey"`
	Keypair []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		IsPub bool   `json:"is_pub"`
	} `json:"keyPair"`
}

type SpencialOp struct {
	IsCA         bool     `json:"isCA"`
	ValidatorPub string   `json:"validator_pubkey"`
	Sigs         []string `json:"sigs"`
	// AccountPub   string `json:"account_pubkey"`
	Power uint64 `json:"power"`
}

type Result_BigData struct {
	IsSuccess bool           `json:"isSuccess"`
	Result    *RspGetBigData `json:"result"`
}
type RspGetBigData struct {
	Account    string `json:"account"`
	Created_at int64  `json:"created_at"`
	Ext        string `json:"ext"`
	Hash       string `json:"hash"`
	Id         int64  `json:"id"`
	Memo       string `json:"memo"`
	Type       string `json:"type"`
	Value      string `json:"value"`
}

type RPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
}

type ReqThorizeData struct {
	// Id          string `json:"id"`
	Hashs   []string `json:"hashs"`
	Privkey string   `json:"privkey"`
	// Source      string `json:"source"`
	Destination string `json:"destination"`
	// Encrypt     bool   `json:"encrypt"`
}
type ReqGetThorizeData struct {
	// Id      string `json:"id"`
	Privkey string `json:"privkey"`
	// Destination string `json:"destination"`
	Hashs []string `json:"hashs"`
}

type StCreateContract struct {
	OperationBase
	ByteCode string        `json:"bytecode"`
	Abi      string        `json:"abi"`
	Params   []interface{} `json:"params"`
	Price    string        `json:"gas_price"`
	GasLimit string        `json:"gas_limit"`
	Amount   string        `json:"amount"`
}

type StExcuteContract struct {
	OperationBase
	ContractAddr string        `json:"contract_address"`
	Abi          string        `json:"abi"`
	Func         string        `json:"function"`
	Params       []interface{} `json:"params"`
	Amount       string        `json:"amount"`
	Price        string        `json:"gas_price"`
	GasLimit     string        `json:"gas_limit"`
}

type StQueryContract struct {
	PrivKey      string        `json:"privkey"`
	Source       string        `json:"source"`
	ContractAddr string        `json:"contract_address"`
	Abi          string        `json:"abi"`
	Func         string        `json:"function"`
	Params       []interface{} `json:"params"`
}

type StReceipt struct {
	OpType          string            `json:"type"`
	Source          string            `json:"source"`
	TxHash          string            `json:"hash"`
	TxReceiptStatus bool              `json:"txReceiptStatus"`
	Message         string            `json:"msg"`
	Res             interface{}       `json:"res"`
	Height          uint64            `json:"height"`
	ContractAddress string            `json:"contract_address"`
	Function        string            `json:"function"`
	Params          []interface{}     `json:"params"`
	GasPrice        string            `json:"gas_price"`
	GasLimit        string            `json:"gas_limit"`
	GasUsed         *big.Int          `json:"gasUsed"`
	Logs            []*delostypes.Log `json:"logs"`
}

type stTx struct {
	Res      string      `json:"res"`
	TxHash   string      `json:"txhash"`
	CodeType at.CodeType `json:"codetype"`
}
type stRes struct {
	IsSuccess bool `json:"isSuccess, bool"`
	Result    stTx `json:"result"`
}
