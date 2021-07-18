package repository

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type anyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestClientApplicationSQLRepository_Register(t *testing.T) {
	clientApp := test_data.ClientAppPing
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &clientApplicationSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "client_applications",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO client_applications (id, secret, name, scope, token_mode, client_notification_endpoint, authentication_request_signing_alg, user_code_parameter_supported, redirect_uri, token_endpoint_auth_method, token_endpoint_auth_signing_alg, grant_types, public_key_uri) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).WithArgs(clientApp.Id, clientApp.Secret, clientApp.Name, clientApp.Scope, clientApp.TokenMode, clientApp.ClientNotificationEndpoint, clientApp.AuthenticationRequestSigningAlg, clientApp.UserCodeParameterSupported, clientApp.RedirectUri, clientApp.TokenEndpointAuthMethod, clientApp.TokenEndpointAuthSigningAlg, clientApp.GrantTypes, clientApp.PublicKeyUri).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Register(&clientApp)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
}

func TestClientApplicationSQLRepository_FindById(t *testing.T) {
	clientApp := test_data.ClientAppPing
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &clientApplicationSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "client_applications",
	}
	rows := sqlmock.NewRows([]string{"id", "secret", "name", "scope", "token_mode", "client_notification_endpoint", "authentication_request_signing_alg", "user_code_parameter_supported", "redirect_uri", "token_endpoint_auth_method", "token_endpoint_auth_signing_alg", "grant_types", "public_key_uri"}).
		AddRow(clientApp.Id, clientApp.Secret, clientApp.Name, clientApp.Scope, clientApp.TokenMode, clientApp.ClientNotificationEndpoint, clientApp.AuthenticationRequestSigningAlg, clientApp.UserCodeParameterSupported, clientApp.RedirectUri, clientApp.TokenEndpointAuthMethod, clientApp.TokenEndpointAuthSigningAlg, clientApp.GrantTypes, clientApp.PublicKeyUri)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM client_applications WHERE id = ?")).
		WithArgs(clientApp.Id).
		WillReturnRows(rows)

	ca, err := repo.FindById(clientApp.Id)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
	assert.NotNil(t, ca)
}

func TestAccessTokenSQLRepository_Create(t *testing.T) {
	accesToken := test_data.AccessTokenValid
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &accessTokenSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "access_tokens",
	}
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO access_tokens (access_token, client_id, expires, user_id, scope) VALUES (?, ?, ?, ?, ?)")).
		WithArgs(accesToken.Value, accesToken.ClientId, accesToken.Expires, accesToken.UserId, accesToken.Scope).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(&accesToken)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
}

func TestAccessTokenSQLRepository_Find(t *testing.T) {
	accessToken := test_data.AccessTokenValid
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &accessTokenSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "access_tokens",
	}

	rows := sqlmock.NewRows([]string{"access_token", "client_id", "expires", "user_id", "scope"}).
		AddRow(accessToken.Value, accessToken.ClientId, accessToken.Expires, accessToken.UserId, accessToken.Scope)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM access_tokens WHERE access_token = ?")).
		WithArgs(accessToken.Value).WillReturnRows(rows)

	at, err := repo.Find(accessToken.Value)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
	assert.NotNil(t, at)
}

func TestCibaSessionSQLRepository_Create(t *testing.T) {
	cibaSession := test_data.CibaSession6
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &cibaSessionSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "ciba_sessions",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ciba_sessions (auth_req_id, client_id, user_id, hint, binding_message, client_notification_token, expires_in, interval, valid, id_token, consented, scope, latest_token_requested_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
		WithArgs(cibaSession.AuthReqId, cibaSession.ClientId, cibaSession.UserId, cibaSession.Hint, cibaSession.BindingMessage, cibaSession.ClientNotificationToken, cibaSession.ExpiresIn, cibaSession.Interval, cibaSession.Valid, cibaSession.IdToken, cibaSession.Consented, cibaSession.Scope, cibaSession.LatestTokenRequestedAt, cibaSession.CreatedAt.Format(time.RFC3339)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(&cibaSession)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
}

func TestCibaSessionSQLRepository_FindById(t *testing.T) {
	cibaSession := test_data.CibaSession6
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &cibaSessionSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "ciba_sessions",
	}
	rows := sqlmock.NewRows([]string{"auth_req_id", "client_id", "user_id", "hint", "binding_message", "client_notification_token", "expires_in", "interval", "valid", "id_token", "consented", "scope", "latest_token_requested_at", "created_at"}).
		AddRow(cibaSession.AuthReqId, cibaSession.ClientId, cibaSession.UserId, cibaSession.Hint, cibaSession.BindingMessage, cibaSession.ClientNotificationToken, cibaSession.ExpiresIn, cibaSession.Interval, cibaSession.Valid, cibaSession.IdToken, cibaSession.Consented, cibaSession.Scope, cibaSession.LatestTokenRequestedAt, cibaSession.CreatedAt)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM ciba_sessions WHERE auth_req_id = ?")).
		WillReturnRows(rows)

	cs, err := repo.FindById(cibaSession.AuthReqId)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
	assert.NotNil(t, cs)
}

func TestCibaSessionSQLRepository_Update(t *testing.T) {
	cibaSession := test_data.CibaSession6
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &cibaSessionSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "ciba_sessions",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE ciba_sessions SET client_id = ?, user_id = ?, hint = ?, binding_message = ?, client_notification_token = ?, expires_in = ?, interval = ?, valid = ?, id_token = ?, consented = ?, scope = ?, latest_token_requested_at = ? WHERE auth_req_id = ?")).
		WithArgs(cibaSession.ClientId, cibaSession.UserId, cibaSession.Hint, cibaSession.BindingMessage, cibaSession.ClientNotificationToken, cibaSession.ExpiresIn, cibaSession.Interval, cibaSession.Valid, cibaSession.IdToken, cibaSession.Consented, cibaSession.Scope, cibaSession.LatestTokenRequestedAt, cibaSession.AuthReqId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(&cibaSession)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
}

func TestKeySQLRepository_FindPrivateKeyByClientId(t *testing.T) {
	key := test_data.Key6
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := &keySQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "keys",
	}
	rows := sqlmock.NewRows([]string{"id", "client_id", "alg", "public", "private"}).
		AddRow(key.Id, key.ClientId, key.Alg, key.Public, key.Private)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM keys WHERE client_id = ?")).
		WithArgs(key.ClientId).
		WillReturnRows(rows)

	keyRes, err := repo.FindPrivateKeyByClientId(key.ClientId)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
	assert.NotNil(t, keyRes)
}

func TestUserAccountSQLRepository_FindById(t *testing.T) {
	userAccount := test_data.User3
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := userAccountSQLRepository{
		db:        sqlx.NewDb(mockDb, ""),
		tableName: "user_accounts",
	}
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_code", "created_at", "updated_at"}).
		AddRow(userAccount.Id, userAccount.Name, userAccount.Email, userAccount.Password, userAccount.UserCode, userAccount.CreatedAt, userAccount.UpdatedAt)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM user_accounts WHERE id = ?")).
		WithArgs(userAccount.Id).
		WillReturnRows(rows)

	user, err := repo.FindById(userAccount.Id)
	mockErr := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NoError(t, mockErr)
	assert.NotNil(t, user)
}

func TestUserClaimSQLRepository_GetUserClaims(t *testing.T) {
	userAccount := test_data.User3
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	repo := userClaimSQLRepository{
		db:                   sqlx.NewDb(mockDb, ""),
		tableNameScopes:      "scopes",
		tableNameClaims:      "claims",
		tableNameUsers:       "user_accounts",
		tableNameScopeClaims: "scope_claims",
	}

	userRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "user_code", "created_at", "updated_at"}).
		AddRow(userAccount.Id, userAccount.Name, userAccount.Email, userAccount.Password, userAccount.UserCode, userAccount.CreatedAt, userAccount.UpdatedAt)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM user_accounts WHERE id = ? LIMIT 1")).
		WithArgs(userAccount.Id).
		WillReturnRows(userRows)

	// scopeIdSql := u.db.Rebind(fmt.Sprintf("SELECT id FROM %s WHERE name = ? LIMIT 1", u.tableNameScopes))
	scopeRows := sqlmock.NewRows([]string{"id"}).
		AddRow("id.timestamp.read")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM scopes WHERE name = ? LIMIT 1")).
		WithArgs("timestamp.read").
		WillReturnRows(scopeRows)
	claimsRows := sqlmock.NewRows([]string{"name"}).
		AddRow("created_at").
		AddRow("updated_at")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM claims WHERE id IN (SELECT id in scope_claims WHERE scope_id = ?)")).
		WithArgs("id.timestamp.read").
		WillReturnRows(claimsRows)

	claims, err := repo.GetUserClaims(userAccount.Id, "timestamp.read")
	mockErr := mock.ExpectationsWereMet()

	for _, c := range []string{"created_at", "updated_at"} {
		v, ok := claims[c]
		assert.True(t, ok)
		assert.NotNil(t, v)
	}
	assert.Nil(t, err)
	assert.Nil(t, mockErr)
}

func TestBuildTableName(t *testing.T) {
	var res string
	res = buildTableName("", "mytable")
	assert.Equal(t, "mytable", res)
	res = buildTableName("go_ciba", "mytable")
	assert.Equal(t, "go_ciba_mytable", res)
}
