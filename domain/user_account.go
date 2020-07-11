package domain

import "encoding/json"

type UserAccount struct {
	Name string
	Email string
	Address string
	Password string
	UserCode string
}

func (ua *UserAccount) MarshalBinary() ([]byte, error) {
	return json.Marshal(ua)
}

func (ua *UserAccount) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	return nil
}

func (ua *UserAccount) GetUseCode() string {
	return ua.UserCode
}

func (ua *UserAccount) SetUserCode(code string) string {
	ua.UserCode = code
}

