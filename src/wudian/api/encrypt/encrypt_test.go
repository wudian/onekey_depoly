package encrypt

import (
	"encoding/hex"
	"fmt"
	"testing"

	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
)


func TestRlpHash(t *testing.T) {
	pri, pub, _ := GenKey()
	pub2, _ := PrivToPub(pri)
	fmt.Println(pub, pub2)

	a := "123456789012345678901234567890123456789012345678901234567890"
	fmt.Println(RlpHash(a), RlpHash(123))
}

func TestGenKey(t *testing.T) {
	a := "123456789012345678901234567890123456789012345678901234567890"
	data := []byte(a)
	data2 := ethcmn.Hex2Bytes(a)
	fmt.Println(data, data2)
	fmt.Println(string(data), ethcmn.Bytes2Hex(data))

	pri, pub, _ := GenKey()
	sk := crypto.ToECDSA(ethcmn.Hex2Bytes(pri))
	pk := &sk.PublicKey
	pub2 := hex.EncodeToString(crypto.FromECDSAPubCompressed(pk))
	pub3 := crypto.PubkeyToAddress(*pk).Hex()
	pri2 := hex.EncodeToString(crypto.FromECDSA(sk))
	pub4, _ := PrivToPub(pri)

	pub5 := ethcmn.Bytes2Hex(crypto.FromECDSAPubCompressed(pk))
	pri5 := ethcmn.Bytes2Hex(crypto.FromECDSA(sk))

	fmt.Println(pub, pub2, pub3, pub4, pub5)
	fmt.Println(pri, pri2, pri5)

	b := ethcmn.IsHexAddress(pub5)
	fmt.Println(b)
}

// 测试加密、解密
func TestEncode(t *testing.T) {
	data := []byte("123456789012345678901234567890123456789012345678901234567890")
	t.Log("data", string(data))

	pri, pub, err := GenKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(pub), pri)

	enData, err := EncodeData(pub, data)
	if err != nil {
		t.Fatal(err)
	}
	hex1 := ethcmn.Bytes2Hex(enData)
	hex2 := hex.EncodeToString(enData)
	t.Log("hex1", hex1, len(hex1), len(enData)) 
	t.Log("hex2", hex2, len(hex2), len(enData)) 
	// t.Log("enData", string(enData))

	bytes := ethcmn.Hex2Bytes(hex1)
	deData, err := DecodeData(pri, bytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("deData", string(deData))

	// 用另一个私钥来解密
	srcpriv := "bbcc74ebcc7d731e1b7724c9279e2caed1d8e324cbd1eb3586126848366eb822"
	deData2, err := DecodeData(srcpriv, enData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(deData2), len(deData2))
}
