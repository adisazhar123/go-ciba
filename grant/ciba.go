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
}

func NewCibaGrant() *CibaGrant {
	return &CibaGrant{}
}

func (cg *CibaGrant) GetIdentifier() string {
	return IdentifierCiba
}

func (cg *CibaGrant) SetInterval(val *int) {
	cg.PollInterval = val
}

func (cg *CibaGrant) CreateAccessTokenAndIdToken() *domain.Tokens {
	return nil
}

func (cg *CibaGrant) CreateIdToken(userId, clientId, accessToken string) {
	//claims := domain.DefaultCibaIdTokenClaims{
	//	DefaultIdTokenClaims: domain.DefaultIdTokenClaims{
	//		Iss:      cg.Config.Issuer,
	//		Sub:      userId,
	//		Aud:      clientId,
	//		Exp:      int(time.Now().Unix()) + cg.Config.IdTokenLifetime,
	//		Iat:      0,
	//		AuthTime: int(time.Now().Unix()),
	//	},
	//	AtHash:               "",
	//	RtHash:               "",
	//	AuthReqId:            "",
	//}
}

func accessTokenHash(accessToken, clientId string) string {
	return ""
}
