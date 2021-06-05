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
	Exp int `json:"exp"`
	// 	Time at which the JWT was issued.
	Iat int `json:"iat"`
	// Time when the End-User authentication occurred.
	AuthTime int `json:"auth_time"`
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
	Value    string
	ClientId string
	Expires  int
	UserId   string
	Scope    string
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
	now := int(time.Now().Unix())
	return at.Expires < now
}

func NewAccessToken(value, clientId, userId, scope string, expires int) *AccessToken {
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
	ExpiresIn int
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

	return base64.URLEncoding.EncodeToString([]byte(tokenHash))
}
