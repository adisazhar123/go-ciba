package grant

type GrantTypeInterface interface {
	GetIdentifier() string
}

type GrantConfig struct {
	Issuer              string
	IdTokenLifetime     int
	AccessTokenLifetime int
}
