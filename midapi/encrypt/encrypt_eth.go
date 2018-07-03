package encrypt

import (
	"encoding/hex"

	// ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto"
	// "crypto/ecdsa"
)

func GenKey() (pri, pub string, err error) {
	sk, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	pk := &sk.PublicKey
	// pub = crypto.PubkeyToAddress(*pk).Hex()
	pub = hex.EncodeToString(crypto.FromECDSAPubCompressed(pk))
	pri = hex.EncodeToString(crypto.FromECDSA(sk))
	return
}

func PrivToPub(pri string) (pub string, err error) {
	// sk := crypto.ToECDSA(ethcmn.Hex2Bytes(pri))
	sk, err := crypto.HexToECDSA(pri)
	if err != nil {
		return "", err
	}
	pk := &sk.PublicKey
	pub = hex.EncodeToString(crypto.FromECDSAPubCompressed(pk))
	return
}

// func ValidPublicKey(pub string) (err error) {
// 	_, err = crypto.ToECDSAPubCompressed(ethcmn.Hex2Bytes(pub))
// 	return
// }

// func EncodeData(pub string, in []byte) ([]byte, error) {
// 	pk, err := crypto.ToECDSAPubCompressed(ethcmn.Hex2Bytes(pub))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return crypto.Encrypt(pk, in)
// }

// func DecodeData(pri string, in []byte) ([]byte, error) {
// 	sk := crypto.ToECDSA(ethcmn.Hex2Bytes(pri))
// 	return crypto.Decrypt(sk, in)
// }
