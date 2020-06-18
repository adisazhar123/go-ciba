package grant

import "github.com/adisazhar123/ciba-server/domain"

type CibaGrantTypeInterface interface {
	InitRepositories(repo1, repo2 string)
	GrantTypeInterface
}

type CibaGrant struct {
	authenticationRequestId string
	idToken string
	accessToken domain.AccessToken
	clientApplication domain.ClientApplication
}

func NewCibaGrant() *CibaGrant {
	return &CibaGrant{
		authenticationRequestId: "",
		idToken:                 "",
		accessToken:             domain.AccessToken{},
		clientApplication:       domain.ClientApplication{},
	}
}

func (cg *CibaGrant) GetIdentifier() string {
	return "urn:openid:params:grant-type:ciba"
}

func (cg *CibaGrant) ValidateAuthenticationRequest() {

}

func (cg *CibaGrant) HandleAuthenticationRequest() {

}

func (cg *CibaGrant) InitRepositories(repo1, repo2 string) {

}