package grant

type GrantTypeInterface interface {
	GetIdentifier() string
}

type GrantConfig struct {
	Issuer              string
	IdTokenLifetime     int64
	AccessTokenLifetime int64
}
