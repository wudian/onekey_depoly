package handlers

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto/sha3"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/types"
)

func rlpHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func TestHash(t *testing.T) {

	tx := types.NewTransaction(1, big.NewInt(100), ethcmn.HexToAddress("033bd2084f02c7dd26615b38cf8ed1833f556010967acc723972f1e159b2eeb805"),
		types.TimeBounds{Min: 0, Max: 10000},
		nil, []byte("123"), 9999999999)

	var privkeys []*ecdsa.PrivateKey

	pri := "bb8a20086b3956c0c369da93d33f4af9a513967655b40d347633838a601897cf"

	privkey := crypto.ToECDSA(ethcmn.Hex2Bytes(pri))

	privkeys = append(privkeys, privkey)

	fmt.Println(tx.Hash().Hex(), rlpHash(tx).Hex(), tx.Data.Sigs)

	sigtx, err := tx.Sign(privkeys)
	if err != nil {
		t.Log(err)
	}

	fmt.Println(sigtx.Hash().Hex(), rlpHash(sigtx).Hex(), sigtx.Data.Sigs)
	fmt.Println(tx.Hash().Hex(), rlpHash(tx).Hex(), tx.Data.Sigs)
}

// 结论：签名前和签名后算出的交易哈希是一样的 --- 吴典
