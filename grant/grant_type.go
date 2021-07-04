package grant

type GrantTypeInterface interface {
	GetIdentifier() string
}

type GrantConfig struct {
	Issuer                            string
	IdTokenLifetimeInSeconds          int64
	AccessTokenLifetimeInSeconds      int64
	DefaultAuthReqIdLifetimeInSeconds int64
	TokenEndpointUrl                  string
}
