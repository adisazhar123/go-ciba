package domain

import (
	"encoding/json"
	"time"
)

type UserAccount struct {
	Id        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	UserCode  string    `db:"user_code" json:"user_code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
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
