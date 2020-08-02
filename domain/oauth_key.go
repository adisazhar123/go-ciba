package domain

import (
	"crypto"
	"encoding/json"
)

type Key struct {
	ClientId string
	Alg      string
	Public   crypto.PublicKey
	Private  crypto.PrivateKey
}

func (k *Key) MarshalBinary() ([]byte, error) {
	return json.Marshal(k)
}

func (k *Key) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	return nil
}
