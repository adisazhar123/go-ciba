package domain

import "github.com/adisazhar123/go-ciba/util"

type AccessToken struct {
	value     string
	tokenType string
	expiresIn string
}

type DecodedIdToken struct {
}

type EncodedIdToken struct {
	value string
}

type IdTokenManager struct {
	e util.EncryptionInterface
}

type Tokens struct {
	IdToken     EncodedIdToken
	AccessToken AccessToken
}

type IdTokenInterface interface {
	CreateIdToken(claims interface{}, key, alg, keyId string) *EncodedIdToken
}

func (tkn *IdTokenManager) CreateIdToken(claims interface{}, key, alg, keyId string) *EncodedIdToken {
	return &EncodedIdToken{value: tkn.e.Encode(claims, key, alg, keyId)}
}

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
	// String value used to associate a Client session with an ID Token, and to mitigate replay attacks.
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
	// Authentication request id value
	AuthReqId string `json:"urn:openid:params:jwt:claim:auth_req_id,omitempty"`
}
