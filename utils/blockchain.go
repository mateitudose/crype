package utils

import (
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

func GeneratePaymentAddress(currency string) (string, error) {
	switch currency {
	case "USDC_BASE":
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return "", errors.New("could not generate private key")
		}
		publicKey := privateKey.PublicKey
		address := crypto.PubkeyToAddress(publicKey)
		return address.Hex(), nil
	default:
		return "", errors.New("currency is not supported")
	}
}
