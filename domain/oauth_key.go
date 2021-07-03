package domain

import (
	"encoding/json"
)

type Key struct {
	Id       string `db:"id"`
	ClientId string `db:"client_id"`
	Alg      string `db:"alg"`
	Public   string `db:"public"`
	Private  string `db:"private"`
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
