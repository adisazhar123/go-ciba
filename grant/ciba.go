package grant

const (
	IdentifierCiba = "urn:openid:params:grant-type:ciba"
)

type CibaGrantTypeInterface interface {
	InitRepositories(repo1, repo2 string)
	GrantTypeInterface
}

type CibaGrant struct{}

func NewCibaGrant() *CibaGrant {
	return &CibaGrant{}
}

func (cg *CibaGrant) GetIdentifier() string {
	return IdentifierCiba
}
