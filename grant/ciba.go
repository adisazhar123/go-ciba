package grant

import (
	"github.com/adisazhar123/go-ciba/domain"
)

const (
	IdentifierCiba = "urn:openid:params:grant-type:ciba"
)

type CibaGrantTypeInterface interface {
	GrantTypeInterface

	InitRepositories(repo1, repo2 string)
	SetInterval(val int)
}

type CibaGrant struct {
	PollInterval *int
	Config       GrantConfig
	TokenManager domain.TokenInterface
}

// TODO: add passable config
func NewCibaGrant() *CibaGrant {
	return &CibaGrant{
		PollInterval: nil,
		Config: GrantConfig{
			Issuer:              "issuer-ciba.example.com",
			IdTokenLifetime:     3600,
			AccessTokenLifetime: 3600,
		},
		TokenManager: domain.NewTokenManager(),
	}
}

func (cg *CibaGrant) GetIdentifier() string {
	return IdentifierCiba
}

func (cg *CibaGrant) SetInterval(val *int) {
	cg.PollInterval = val
}

func formatCibaClaims(defaultClaims domain.DefaultCibaIdTokenClaims, extraClaims map[string]interface{}) map[string]interface{} {
	combinedClaims := make(map[string]interface{})

	combinedClaims["auth_req_id"] = defaultClaims.AuthReqId
	combinedClaims["aud"] = defaultClaims.Aud
	combinedClaims["auth_time"] = defaultClaims.AuthTime
	combinedClaims["iat"] = defaultClaims.Iat
	combinedClaims["exp"] = defaultClaims.Exp
	combinedClaims["iss"] = defaultClaims.Iss
	combinedClaims["sub"] = defaultClaims.Sub
	combinedClaims["nonce"] = defaultClaims.Nonce

	for k, v := range extraClaims {
		combinedClaims[k] = v
	}

	return combinedClaims
}

func (cg *CibaGrant) CreateAccessTokenAndIdToken(defaultClaims domain.DefaultCibaIdTokenClaims, extraClaims map[string]interface{}, key, alg, keyId string) *domain.Tokens {
	claims := formatCibaClaims(defaultClaims, extraClaims)
	accessToken := cg.TokenManager.CreateAccessToken()
	idToken := cg.TokenManager.CreateIdToken(claims, key, alg, keyId, accessToken)

	return &domain.Tokens{
		IdToken: idToken,
		AccessToken: domain.AccessTokenInternal{
			Value:     accessToken,
			TokenType: "bearer",
			ExpiresIn: cg.Config.AccessTokenLifetime,
		},
	}
}
