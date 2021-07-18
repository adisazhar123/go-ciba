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
	DefaultAuthReqIdLifetimeInSeconds   int64 = 120
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
		Config: &GrantConfig{
			IdTokenLifetimeInSeconds:     DefaultIdTokenLifeTimeInSeconds,
			AccessTokenLifetimeInSeconds: DefaultAccessTokenLifeTimeInSeconds,
			AuthReqIdLifetimeInSeconds:   DefaultAuthReqIdLifetimeInSeconds,
			PollingIntervalInSeconds:     &DefaultPollIntervalInSeconds,
			Issuer:                       "issuer-ciba.example.com",
			TokenEndpointUrl:             "issuer-ciba.example.com/token",
		},
		TokenManager: domain.NewTokenManager(),
	}
}

func NewCustomCibaGrant(grantConfig *GrantConfig) *CibaGrant {
	return &CibaGrant{
		Config:       grantConfig,
		TokenManager: domain.NewTokenManager(),
	}
}

func (cg *CibaGrant) GetIdentifier() string {
	return IdentifierCiba
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
			ExpiresIn: cg.Config.AccessTokenLifetimeInSeconds,
		},
	}
}
