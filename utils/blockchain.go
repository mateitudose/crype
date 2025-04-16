package utils

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Address      string
	PrivateKey   string
}

func GeneratePaymentAddress(currency string) (Wallet, error) {
	switch currency {
	case "USDC_BASE":
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return Wallet{}, errors.New("could not generate private key")
		}
		publicKey := privateKey.PublicKey
		address := crypto.PubkeyToAddress(publicKey)
		return Wallet{
			Address:      address.Hex(),
			PrivateKey:   fmt.Sprintf("%x", crypto.FromECDSA(privateKey)),
		}, nil
	default:
		return Wallet{}, errors.New("currency is not supported")
	}
}
