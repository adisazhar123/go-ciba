package grant

type GrantTypeInterface interface {
	GetIdentifier() string
}

type GrantConfig struct {
	Issuer                       string
	IdTokenLifetimeInSeconds     int64
	AccessTokenLifetimeInSeconds int64
	PollingIntervalInSeconds     *int64
	AuthReqIdLifetimeInSeconds   int64
	TokenEndpointUrl             string
}
