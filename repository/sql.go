package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/jmoiron/sqlx"
)

type clientApplicationSQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func (c *clientApplicationSQLRepository) Register(ca *domain.ClientApplication) error {
	cmd := c.db.Rebind(fmt.Sprintf("INSERT INTO %s (id, secret, name, scope, token_mode, client_notification_endpoint, authentication_request_signing_alg, user_code_parameter_supported, redirect_uri, token_endpoint_auth_method, token_endpoint_auth_signing_alg, grant_types, public_key_uri) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", c.tableName))
	_, err := c.db.Exec(cmd, ca.Id, ca.Secret, ca.Name, ca.Scope, ca.TokenMode, ca.ClientNotificationEndpoint, ca.AuthenticationRequestSigningAlg, ca.UserCodeParameterSupported, ca.RedirectUri, ca.TokenEndpointAuthMethod, ca.TokenEndpointAuthSigningAlg, ca.GrantTypes, ca.PublicKeyUri)
	return err
}

func (c *clientApplicationSQLRepository) FindById(id string) (*domain.ClientApplication, error) {
	var clientApp domain.ClientApplication
	cmd := c.db.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", c.tableName))
	err := c.db.Get(&clientApp, cmd, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &clientApp, nil
}

type accessTokenSQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func (a *accessTokenSQLRepository) Create(at *domain.AccessToken) error {
	cmd := a.db.Rebind(fmt.Sprintf("INSERT INTO %s (access_token, client_id, expires, user_id, scope) VALUES (?, ?, ?, ?, ?)", a.tableName))
	_, err := a.db.Exec(cmd, at.Value, at.ClientId, at.Expires, at.UserId, at.Scope)
	return err
}

func (a *accessTokenSQLRepository) Find(at string) (*domain.AccessToken, error) {
	var accessToken domain.AccessToken
	cmd := a.db.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE access_token = ?", a.tableName))
	err := a.db.Get(&accessToken, cmd, at)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &accessToken, nil
}

type cibaSessionSQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func (c *cibaSessionSQLRepository) Create(cs *domain.CibaSession) error {
	cmd := c.db.Rebind(fmt.Sprintf("INSERT INTO %s (auth_req_id, client_id, user_id, hint, binding_message, client_notification_token, expires_in, interval, valid, id_token, consented, scope, latest_token_requested_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", c.tableName))
	_, err := c.db.Exec(cmd, cs.AuthReqId, cs.ClientId, cs.UserId, cs.Hint, cs.BindingMessage, cs.ClientNotificationToken, cs.ExpiresIn, cs.Interval, cs.Valid, cs.IdToken, cs.Consented, cs.Scope, cs.LatestTokenRequestedAt, cs.CreatedAt.Format(time.RFC3339))
	return err
}

func (c *cibaSessionSQLRepository) FindById(id string) (*domain.CibaSession, error) {
	var cibaSession domain.CibaSession
	cmd := c.db.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE auth_req_id = ?", c.tableName))
	err := c.db.Get(&cibaSession, cmd, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cibaSession, nil
}

func (c *cibaSessionSQLRepository) Update(cs *domain.CibaSession) error {
	cmd := c.db.Rebind(fmt.Sprintf("UPDATE %s SET auth_req_id = ?, client_id = ?, user_id = ?, hint = ?, binding_message = ?, client_notification_token = ?, expires_in = ?, interval = ?, valid = ?, id_token = ?, consented = ?, scope = ?, latest_token_requested_at = ?", c.tableName))
	_, err := c.db.Exec(cmd, cs.AuthReqId, cs.ClientId, cs.UserId, cs.Hint, cs.BindingMessage, cs.ClientNotificationToken, cs.ExpiresIn, cs.Interval, cs.Valid, cs.IdToken, cs.Consented, cs.Scope, cs.LatestTokenRequestedAt)
	return err
}

type keySQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func (k *keySQLRepository) FindPrivateKeyByClientId(clientId string) (*domain.Key, error) {
	var key domain.Key
	cmd := k.db.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE client_id = ?", k.tableName))
	err := k.db.Get(&key, cmd, clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &key, nil
}

type userAccountSQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func (u *userAccountSQLRepository) FindById(id string) (*domain.UserAccount, error) {
	var user domain.UserAccount
	cmd := u.db.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", u.tableName))
	err := u.db.Get(&user, cmd, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

type SQLDataStore struct {
	accessTokenRepo       *accessTokenSQLRepository
	cibaSessionRepo       *cibaSessionSQLRepository
	clientApplicationRepo *clientApplicationSQLRepository
	keyRepositoryRepo     *keySQLRepository
	userAccountRepo       *userAccountSQLRepository
}

func buildTableName(prefix, tableName string) string {
	if prefix == "" {
		return tableName
	}
	return fmt.Sprintf("%s_%s", prefix, tableName)
}

func NewSQLDataStore(defaultDb *sql.DB, driverName, prefix string) *SQLDataStore {
	db := sqlx.NewDb(defaultDb, driverName)
	return &SQLDataStore{
		accessTokenRepo: &accessTokenSQLRepository{
			db:        db,
			tableName: buildTableName(prefix, "access_tokens"),
		},
		cibaSessionRepo: &cibaSessionSQLRepository{
			db:        db,
			tableName: buildTableName(prefix, "ciba_sessions"),
		},
		clientApplicationRepo: &clientApplicationSQLRepository{
			db:        db,
			tableName: buildTableName(prefix, "client_applications"),
		},
		keyRepositoryRepo: &keySQLRepository{
			db:        db,
			tableName: buildTableName(prefix, "keys"),
		},
		userAccountRepo: &userAccountSQLRepository{
			db:        db,
			tableName: buildTableName(prefix, "user_accounts"),
		},
	}
}

func (s *SQLDataStore) GetAccessTokenRepository() AccessTokenRepositoryInterface {
	return s.accessTokenRepo
}

func (s *SQLDataStore) GetCibaSessionRepository() CibaSessionRepositoryInterface {
	return s.cibaSessionRepo
}

func (s *SQLDataStore) GetClientApplicationRepository() ClientApplicationRepositoryInterface {
	return s.clientApplicationRepo
}

func (s *SQLDataStore) GetKeyRepository() KeyRepositoryInterface {
	return s.keyRepositoryRepo
}

func (s *SQLDataStore) GetUserAccountRepository() UserAccountRepositoryInterface {
	return s.userAccountRepo
}
