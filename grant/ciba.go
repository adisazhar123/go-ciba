package grant

import (
	"github.com/adisazhar123/go-ciba/domain"
)

const (
	IdentifierCiba = "urn:openid:params:grant-type:ciba"
)

var (
	DefaultPollIntervalInSeconds        int64 = 5
	DefaultIdTokenLifeTimeInSeconds     int64 = 3600
	DefaultAccessTokenLifeTimeInSeconds int64 = 3600
)

type CibaGrantTypeInterface interface {
	GrantTypeInterface

	InitRepositories(repo1, repo2 string)
	SetInterval(val int64)
}

type CibaGrant struct {
	PollInterval *int64
	Config       *GrantConfig
	TokenManager domain.TokenInterface
}

func NewCibaGrant() *CibaGrant {
	return &CibaGrant{
		PollInterval: &DefaultPollIntervalInSeconds,
		Config: &GrantConfig{
			Issuer:              "issuer-ciba.example.com",
			IdTokenLifetime:     DefaultIdTokenLifeTimeInSeconds,
			AccessTokenLifetime: DefaultAccessTokenLifeTimeInSeconds,
		},
		TokenManager: domain.NewTokenManager(),
	}
}

func NewCustomCibaGrant(pollInterval *int64, grantConfig *GrantConfig) *CibaGrant {
	return &CibaGrant{
		PollInterval: pollInterval,
		Config:       grantConfig,
		TokenManager: domain.NewTokenManager(),
	}
}

func (cg *CibaGrant) GetIdentifier() string {
	return IdentifierCiba
}

func (cg *CibaGrant) SetInterval(val *int64) {
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
