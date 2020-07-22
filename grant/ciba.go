package grant

const (
	IdentifierCiba = "urn:openid:params:grant-type:ciba"
)

type CibaGrantTypeInterface interface {
	InitRepositories(repo1, repo2 string)
	GrantTypeInterface
	SetInterval(val int)
}

type CibaGrant struct{
	PollInterval *int
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
