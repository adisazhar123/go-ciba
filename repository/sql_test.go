package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"
)

type anyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a anyTime) Match(v driver.Value) bool {
	fmt.Print(v)
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

func TestBuildTableName(t *testing.T) {
	var res string
	res = buildTableName("", "mytable")
	assert.Equal(t, "mytable", res)
	res = buildTableName("go_ciba", "mytable")
	assert.Equal(t, "go_ciba_mytable", res)
}
