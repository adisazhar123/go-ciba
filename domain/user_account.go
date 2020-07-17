package domain

import (
	"encoding/json"
	"time"
)

type UserAccount struct {
	Id        string    `json:"Id"`
	Name      string    `json:"Name"`
	Email     string    `json:"Email"`
	Password  string    `json:"Password"`
	UserCode  string    `json:"UserCode"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func (ua *UserAccount) MarshalBinary() ([]byte, error) {
	return json.Marshal(ua)
}

func (ua *UserAccount) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, ua); err != nil {
		return err
	}

	return nil
}

func (ua *UserAccount) GetUseCode() string {
	return ua.UserCode
}

func (ua *UserAccount) SetUserCode(code string) {
	ua.UserCode = code
}

func (ua *UserAccount) GetId() string {
	return ua.Id
}

func (ua *UserAccount) SetId(id string) {
	ua.Id = id
}

func (ua *UserAccount) SetName(name string) {
	ua.Name = name
}

func (ua *UserAccount) GetName() string {
	return ua.Name
}

func (ua *UserAccount) SetEmail(email string) {
	ua.Email = email
}

func (ua *UserAccount) GetEmail() string {
	return ua.Email
}

func (ua *UserAccount) SetPassword(password string) {
	ua.Password = password
}

func (ua *UserAccount) GetPassword() string {
	return ua.Password
}
