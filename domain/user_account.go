package domain

import (
	"encoding/json"
	"time"
)

type UserAccount struct {
	id        string
	name      string
	email     string
	password  string
	userCode  string
	createdAt time.Time
	updatedAt time.Time
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
	return ua.userCode
}

func (ua *UserAccount) SetUserCode(code string) {
	ua.userCode = code
}

func (ua *UserAccount) GetId() string {
	return ua.id
}

func (ua *UserAccount) SetId(id string) {
	ua.id = id
}

func (ua *UserAccount) SetName(name string) {
	ua.name = name
}

func (ua *UserAccount) GetName() string {
	return ua.name
}

func (ua *UserAccount) SetEmail(email string) {
	ua.email = email
}

func (ua *UserAccount) GetEmail() string {
	return ua.email
}

func (ua *UserAccount) SetPassword(password string) {
	ua.password = password
}

func (ua *UserAccount) GetPassword() string {
	return ua.password
}
