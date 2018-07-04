package encrypt

import (
	"encoding/hex"
	"strings"
	"github.com/btcsuite/btcd/btcec"
	ethcmn "gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/common"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/crypto/sha3"
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/eth/rlp"
)

func GenKey() (privKey, publicKey string, err error) {

	privkey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return
	}

	privKey = hex.EncodeToString(privkey.ToECDSA().D.Bytes())

	pubkey := privkey.PubKey()

	publicKey = hex.EncodeToString(pubkey.SerializeCompressed())

	return
}

func PrivToPub(privKey string) (publicKey string, err error) {
	if strings.Index(privKey, "0x") == 0 {
		privKey = privKey[2:]
	}
	pkBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return
	}
	_, publicK := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	publicKey = hex.EncodeToString(publicK.SerializeCompressed())

	return
}

func ValidPublicKey(publicKey string) (err error) {
	var (
		pubKeyBytes []byte
	)
	if pubKeyBytes, err = hex.DecodeString(publicKey); err != nil {
		return
	}
	_, err = btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	return
}

func EncodeData(publicKey string, in []byte) ([]byte, error) {
	if strings.Index(publicKey, "0x") == 0 {
		publicKey = publicKey[2:]
	}
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcec.Encrypt(pubKey, in)

}

func DecodeData(privateKey string, in []byte) ([]byte, error) {

	pkBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	return btcec.Decrypt(privKey, in)

}

func RlpHash(x interface{}) string {
	var h ethcmn.Hash
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])

	s := ethcmn.Bytes2Hex(h.Bytes())
	return s
}
