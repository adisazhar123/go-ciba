package domain

import (
	"encoding/json"
)

type Key struct {
	Id       string `db:"id" json:"id"`
	ClientId string `db:"client_id" json:"client_id"`
	Alg      string `db:"alg" json:"alg"`
	Public   string `db:"public" json:"public"`
	Private  string `db:"private" json:"private"`
}

func (k *Key) MarshalBinary() ([]byte, error) {
	return json.Marshal(k)
}

func (k *Key) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}

	return nil
}
