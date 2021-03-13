package domain

import (
	"encoding/json"
)

type Key struct {
	ID       string
	ClientId string
	Alg      string
	Public   string
	Private  string
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
