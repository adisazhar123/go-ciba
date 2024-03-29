package domain

import (
	"encoding/json"
	"strings"

	"github.com/adisazhar123/go-ciba/util"
)

const (
	ModePing = "ping"
	ModePoll = "poll"
	ModePush = "push"
)

type ClientApplication struct {
	Id                              string `db:"id" json:"id"`
	Secret                          string `db:"secret" json:"secret"`
	Name                            string `db:"name" json:"name"`
	Scope                           string `db:"scope" json:"scope"`
	TokenMode                       string `db:"token_mode" json:"token_mode"`
	ClientNotificationEndpoint      string `db:"client_notification_endpoint" json:"client_notification_endpoint"`
	AuthenticationRequestSigningAlg string `db:"authentication_request_signing_alg" json:"authentication_request_signing_alg"`
	UserCodeParameterSupported      bool   `db:"user_code_parameter_supported" json:"user_code_parameter_supported"`

	RedirectUri                 string `db:"redirect_uri" json:"redirect_uri"`
	TokenEndpointAuthMethod     string `db:"token_endpoint_auth_method" json:"token_endpoint_auth_method"`
	TokenEndpointAuthSigningAlg string `db:"token_endpoint_auth_signing_alg" json:"token_endpoint_auth_signing_alg"`
	GrantTypes                  string `db:"grant_types" json:"grant_types"`
	PublicKeyUri                string `db:"public_key_uri" json:"public_key_uri"`
}

func NewClientApplication(name, scope, tokenMode, clientNotificationEndpoint, authenticationRequestSigningAlg string, userCode bool) *ClientApplication {
	return &ClientApplication{
		Id:                              util.GenerateUuid(),
		Secret:                          util.GenerateRandomString(),
		Name:                            name,
		Scope:                           scope,
		TokenMode:                       tokenMode,
		ClientNotificationEndpoint:      clientNotificationEndpoint,
		AuthenticationRequestSigningAlg: authenticationRequestSigningAlg,
		UserCodeParameterSupported:      userCode,
	}
}

func (ca *ClientApplication) MarshalBinary() ([]byte, error) {
	return json.Marshal(ca)
}

func (ca *ClientApplication) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, ca); err != nil {
		return err
	}

	return nil
}

func (ca *ClientApplication) GetGrantTypes() string {
	return ca.GrantTypes
}

func (ca *ClientApplication) GetClientNotificationEndpoint() string {
	return ca.ClientNotificationEndpoint
}

func (ca *ClientApplication) GetTokenEndpointAuthMethod() string {
	return ca.TokenEndpointAuthMethod
}

func (ca *ClientApplication) GetId() string {
	return ca.Id
}

func (ca *ClientApplication) GetSecret() string {
	return ca.Secret
}

func (ca *ClientApplication) GetName() string {
	return ca.Name
}

func (ca *ClientApplication) GetScope() string {
	return ca.Scope
}

func (ca *ClientApplication) GetTokenMode() string {
	return ca.TokenMode
}

func (ca *ClientApplication) SetId(id string) {
	ca.Id = id
}

func (ca *ClientApplication) SetSecret(secret string) {
	ca.Secret = secret
}

func (ca *ClientApplication) SetName(name string) {
	ca.Name = name
}

func (ca *ClientApplication) SetScope(scope string) {
	ca.Scope = scope
}

func (ca *ClientApplication) SetTokenMode(mode string) {
	ca.TokenMode = mode
}

func (ca *ClientApplication) SetUserCodeSupported(supported bool) {
	ca.UserCodeParameterSupported = supported
}

func (ca *ClientApplication) IsUserCodeSupported() bool {
	return ca.UserCodeParameterSupported
}

func (ca *ClientApplication) GetAuthenticationRequestSigningAlg() string {
	return ca.AuthenticationRequestSigningAlg
}

func (ca *ClientApplication) GetUserCodeParameterSupported() bool {
	return ca.UserCodeParameterSupported
}

func (ca *ClientApplication) IsRegisteredToUseGrantType(grantType string) bool {
	return util.SliceStringContains(strings.Split(ca.GetGrantTypes(), " "), grantType)
}
