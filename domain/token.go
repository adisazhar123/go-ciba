package domain

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"time"

	"github.com/adisazhar123/go-ciba/util"
)

type DefaultIdTokenClaims struct {
	// Required
	// --------
	// Issuer Identifier for the Issuer of the response.
	Iss string `json:"iss"`
	// Subject Identifier.
	Sub string `json:"sub"`
	//  Audience(s) that this ID Token is intended for.
	Aud string `json:"aud"`
	// Expiration time on or after which the ID Token MUST NOT be accepted for processing.
	Exp int64 `json:"exp"`
	// 	Time at which the JWT was issued.
	Iat int64 `json:"iat"`
	// Time when the End-User authentication occurred.
	AuthTime int64 `json:"auth_time"`
	// String Value used to associate a Client session with an ID Token, and to mitigate replay attacks.
	Nonce string `json:"nonce"`

	// Optional
	// --------
	// Authentication Context Class Reference.
	Acr string `json:"act,omitempty"`
	// Authentication Methods References.
	Amr string `json:"amr,omitempty"`
	// Authorized party - the party to which the ID Token was issued.
	Azp string `json:"azp,omitempty"`
}

type DefaultCibaIdTokenClaims struct {
	DefaultIdTokenClaims

	// Access token hash
	AtHash string `json:"at_hash,omitempty"`
	// Refresh token hash
	RtHash string `json:"urn:openid:params:jwt:claim:rt_hash,omitempty"`
	// Authentication request id Value
	AuthReqId string `json:"urn:openid:params:jwt:claim:auth_req_id,omitempty"`
}

type AccessToken struct {
	Value    string    `db:"access_token"`
	ClientId string    `db:"client_id"`
	Expires  time.Time `db:"expires"`
	UserId   string    `db:"user_id"`
	Scope    string    `db:"scope"`
}

func (at *AccessToken) MarshalBinary() ([]byte, error) {
	return json.Marshal(at)
}

func (at *AccessToken) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	return nil
}

func (at *AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	return now.After(at.Expires)
}

func NewAccessToken(value, clientId, userId, scope string, expires time.Time) *AccessToken {
	return &AccessToken{
		Value:    value,
		ClientId: clientId,
		Expires:  expires,
		UserId:   userId,
		Scope:    scope,
	}
}

type AccessTokenInternal struct {
	Value     string
	TokenType string
	ExpiresIn int64
}

type DecodedIdToken struct {
}

type EncodedIdToken struct {
	Value string
}

type TokenManager struct {
	e util.EncryptionInterface
}

type Tokens struct {
	IdToken     EncodedIdToken
	AccessToken AccessTokenInternal
}

type TokenInterface interface {
	CreateIdToken(claims map[string]interface{}, key, alg, keyId, accessToken string) EncodedIdToken
	CreateAccessToken() string
}

func NewTokenManager() *TokenManager {
	return &TokenManager{e: util.NewGoJoseEncryption()}
}

func (tkn *TokenManager) CreateIdToken(claims map[string]interface{}, key, alg, keyId, accessToken string) EncodedIdToken {
	addTokenHashClaim(claims, accessToken, alg)
	token, _ := tkn.e.Encode(claims, key, alg, keyId)
	return EncodedIdToken{Value: token}
}

func (tkn *TokenManager) CreateAccessToken() string {
	return util.GenerateUuid()
}

func addTokenHashClaim(claims map[string]interface{}, token, alg string) {
	claims["at_hash"] = createTokenHash(token, alg)
}

func createTokenHash(token, alg string) string {
	alg = alg[2:]
	hashAlg := fmt.Sprintf("SHA%s", alg)

	var h hash.Hash

	if hashAlg == "SHA256" {
		h = sha256.New()
	} else if hashAlg == "SHA512" {
		h = sha512.New()
	} else {
		panic("hash algorithm not supported")
	}

	h.Write([]byte(token))
	hashed := h.Sum(nil)
	hashStr := fmt.Sprintf("%x", hashed)
	tokenHash := hashStr[:(len(hashStr)/2)-1]

	return base64.StdEncoding.EncodeToString([]byte(tokenHash))
}
