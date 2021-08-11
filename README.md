# go-ciba

[![adisazhar123](https://circleci.com/gh/adisazhar123/go-ciba/tree/development.svg?style=svg&circle-token=7ce9bd81b2ad605ae129d574f18fe39aee783bab)](https://app.circleci.com/pipelines/github/adisazhar123/go-ciba?branch=development)

## Description

`go-ciba` is a server side Software Development Kit (SDK) which attempts to implement the [OpenID Connect Client-Initiated Backchannel Authentication Flow - Core 1.0](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html). This is merely a proof of concept and should not be relied upon in production.

Please feel free to use this for studying purposes.

## Getting Started

### Dependencies

* Go 1.14 or newer
* Firebase Cloud Messaging

### Step-By-Step Walkthrough

The following instructions will provide you to get this library up and running.

#### Initialize your project

You can see [this project](https://github.com/adisazhar123/go-ciba-demo) for a working demo.

#### Define your schema
The `go-ciba` library is storage agnostic which means that it's not tied to a vendor specific database. It can be plug and played with any database by implementing the interfaces in `repository/repo.go`. 

As of now, it comes with a prebuilt SQL and Redis implementation. For the sake of getting it up and running, we'll use the SQL implementation. 

Use the following schema to create the database.

**SQL**

```sql
CREATE TABLE ciba_sessions (
    auth_req_id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    hint VARCHAR(255),
    binding_message VARCHAR(10),
    client_notification_token VARCHAR(255),
    expires_in INT NOT NULL,
    interval INT,
    valid BOOLEAN,
    id_token VARCHAR(2000),
    consented BOOLEAN,
    scope VARCHAR(4000),
    latest_token_requested_at INT,
    created_at TIMESTAMP
);

CREATE TABLE client_applications (
    id VARCHAR(255) PRIMARY KEY,
    secret VARCHAR(255),
    name VARCHAR(255),
    scope VARCHAR(4000),
    token_mode VARCHAR(255),
    client_notification_endpoint VARCHAR(2000),
    authentication_request_signing_alg VARCHAR(10),
    user_code_parameter_supported BOOLEAN,
    redirect_uri VARCHAR(2000),
    token_endpoint_auth_method VARCHAR(20),
    token_endpoint_auth_signing_alg VARCHAR(10),
    grant_types VARCHAR(255),
    public_key_uri VARCHAR(2000)
);

CREATE TABLE keys (
    id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255),
    alg VARCHAR(10),
    public TEXT,
    private TEXT
);

CREATE TABLE access_tokens (
    access_token VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255),
    expires TIMESTAMP,
    user_id VARCHAR(255),
    scope VARCHAR(4000)
);

CREATE TABLE user_accounts (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    user_code VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE scopes (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE claims (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE scope_claims (
    scope_id VARCHAR(255) PRIMARY KEY,
    claim_id  VARCHAR(255)
);
```
Do not use the values below in production. This is merely for example purposes and proof of concept. I do not claim responsibility should a security breach happen.

```sql
INSERT INTO client_applications (id, secret, name, scope, token_mode, client_notification_endpoint, authentication_request_signing_alg, user_code_parameter_supported, redirect_uri, token_endpoint_auth_method, token_endpoint_auth_signing_alg, grant_types, public_key_uri) VALUES ('2a8c10ed-ca2d-42c6-830a-062b379f5e28', 'cb56645e-a250-4bc9-a716-107347929391', 'Client App 1', 'openid bio timestamp.read', 'poll', '', '', false, '', 'client_secret_basic', '', 'urn:openid:params:grant-type:ciba', '');

insert into keys (id, client_id, alg, public, private) values ('e2557d15-6f75-449d-a4f5-357f6e294d87', '2a8c10ed-ca2d-42c6-830a-062b379f5e28', 'RS256', '-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqplqy+c2NbSGMuIRU8t8
saD/rpnWPw2JGf7RCw9PYqXK1AIiGbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58
DzycJnO8Q0D8ViwQ8d5rE3SIoJdFoL/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI
3WXyK8W/jIYdnPFqi39/bUXUYBWKmzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJ
QgGUpiuSSPHeXCDQiBvdaOLoPr4jR3F6exz39AByK5OkVwKENe9J/tfZSVxkrG81
Ud56/Oal1jWJIiQHqCt7s1hMKInjKFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt
+QIDAQAB
-----END PUBLIC KEY-----
', '-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqplqy+c2NbSGMuIRU8t8saD/rpnWPw2JGf7RCw9PYqXK1AIi
GbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58DzycJnO8Q0D8ViwQ8d5rE3SIoJdF
oL/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI3WXyK8W/jIYdnPFqi39/bUXUYBWK
mzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJQgGUpiuSSPHeXCDQiBvdaOLoPr4j
R3F6exz39AByK5OkVwKENe9J/tfZSVxkrG81Ud56/Oal1jWJIiQHqCt7s1hMKInj
KFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt+QIDAQABAoIBABVxv3juEWRi0tOm
kyMDWyNA56Lc949pdihsXX6UaBgwWvSXaA3u1VuqylraP3i6U9zPZ1DP9vAql2zq
RjO59gI8TiyPM8UIcC+szlx45uDFLz9whHIWbvYT9I3bIkrLrNdmS+ubWtoocY/e
aVJOEugxnmVeMBvL6AEIX6o1VqE3h1BrAwLbDdP7T+muxJJC3wiXiSxqRe868AzV
c1eKQJjq+BTdV09bcfMTIZ7aNGgI6F1oZ/NLI3UlwnOiLaWCb8aQyaLENwD0AGDw
X5k91OIvogM9cAhlXwidvyW99SuLuWdf/n+FeXueqIf/gnHu7BoFVi2uc3p9r4xK
F6dUJ+ECgYEA13fXqkoCoQZ2sAwM0gTRIXj0XomTmLBGtpil1P+1l7woUKw+sCxO
3oEjEovazEwyOY5bDYPMFqFtR0rNxp2YNPpo8k16XKI+p312p2BAtpbqLEImN/tY
8idZfeClB1XFN8VSAC8OM3w30BHI+aHHMw32gI29ygApvxOWLKiTlWUCgYEAyrDa
4y4fc+ba0rKKCUoCHvJihiXPIuxfMPyVagCrRgr9WB5NNj4c1kKgHwDrwaz2xrtN
OlGAwn9X0i3e6bcfoQ0nRJslQbn66qqfjeNpqK6CEVv9DIgYVajKLDIkWhdQ6+Si
qn5Vq2MU6NIly34XBFoYfQTmRH0R7azdD+TfBwUCgYEArznt8LXJl4x7H0Zdcrqq
HI+SJAO8PZM1nq9bRXJDCsfg/WJmhL0z0q2wiRelczmQKtCDaeVCJzFWfoDuAdUO
AB+ZE1xA426qh2l4Ajw7xIHMpPuSuzo0JpIrrDvx2Zo+DdHxkuaxpNsjRJoCGEkh
h3qWegtLSiiByruyCFV72CUCgYBgepxGBN9N0PYZ0ogn8cVeq6tABWE6U17gN2p7
gYQFHBgJSKsiBaC+UApdl5egoc75O5CAEOmEKw9HaTQw9Uyl4VfurRan2XnZF4xJ
ApV5iE87KhkiTOmgZG6PaPKqu2x2TGctVmM66De8tsLswMD9/lCnuZxNv2a4Rk8X
UK7kbQKBgECFgpN3sDketKz3DUa2oH6eLKHl1c1VWdnCKs7EJySn6p9nTqPg7B3J
8xe6VGuuU0vo2MqHuZJ+Oudpbz9iXpcyij6OcqCxgy8BV+yPV3WZ/LNQ2fJbsQdS
BE6vbTy42rJAxgWLTkaJuDo7UFIpAw361R59n5nTIk5Mxtq3kIxa
-----END RSA PRIVATE KEY-----
');

INSERT INTO user_accounts (id, name, email, password, user_code, created_at, updated_at) VALUES ('133d0f1e-0256-4616-989c-fa569c217355', 'User 123', 'user123.example@email.com', 'password', '12345', now(), now());

INSERT INTO scopes (id, name) values ('81a10de4-d4ff-4c15-b867-e766c9167a94', 'openid');
INSERT INTO scopes (id, name) values ('ec32bda1-2d18-407e-af55-6f7b5bb7f1fa', 'timestamp.read');

INSERT INTO claims(id, name) VALUES ('fc204a63-be8b-463b-81d6-959f4dc0c1df', 'id');
INSERT INTO claims(id, name) VALUES ('10770265-802d-444c-a980-72d228069c20', 'created_at');
INSERT INTO claims(id, name) VALUES ('37b73cb2-e133-421f-b13b-bba1885c64d6', 'updated_at');

INSERT INTO scope_claims(scope_id, claim_id)
VALUES ('ec32bda1-2d18-407e-af55-6f7b5bb7f1fa', '10770265-802d-444c-a980-72d228069c20'),
       ('ec32bda1-2d18-407e-af55-6f7b5bb7f1fa', '37b73cb2-e133-421f-b13b-bba1885c64d6'),
       ('81a10de4-d4ff-4c15-b867-e766c9167a94', 'fc204a63-be8b-463b-81d6-959f4dc0c1df');
```

**Redis**

```shell
set client_application:b4620189-c368-43ed-b2b4-2186a61fa664 "{\r\n \"id\": \"b4620189-c368-43ed-b2b4-2186a61fa664\",\r\n  \"secret\": \"83e34759-314e-45ec-8211-c6869e053187\",\r\n  \"name\": \"My First Client\",\r\n  \"scope\": \"openid\",\r\n  \"token_mode\": \"poll\",\r\n  \"client_notification_endpoint\": \"\",\r\n  \"authentication_request_signing_alg\": \"\",\r\n  \"user_code_parameter_supported\": false,\r\n  \"redirect_uri\": \"\",\r\n  \"token_endpoint_auth_method\":\"client_secret_basic\",\r\n  \"token_endpoint_auth_signing_alg\": \"\",\r\n  \"grant_types\": \"urn:openid:params:grant-type:ciba\",\r\n  \"public_key_uri\": \"\"\r\n}"

set user_account:f24e0c6d-dbf0-4753-87ad-b554aab423a5 "{\r\n    \"id\": \"f24e0c6d-dbf0-4753-87ad-b554aab423a5\",\r\n    \"name\": \"Joe Foo\",\r\n    \"email\": \"joe@foo.com\",\r\n    \"password\": \"secret\",\r\n    \"user_code\": \"\",\r\n    \"created_at\": \"2021-08-08T19:28:03.700800474+07:00\",\r\n    \"updated_at\": \"2021-08-08T19:28:03.700800474+07:00\"\r\n  }"

set oauth_key:b4620189-c368-43ed-b2b4-2186a61fa664 "{\r\n    \"id\": \"2b075e6c-790c-4ea1-a697-52a382bec9b7\",\r\n    \"client_id\": \"b4620189-c368-43ed-b2b4-2186a61fa664\",\r\n    \"alg\":\"RS256\",\r\n    \"public\": \"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqplqy+c2NbSGMuIRU8t8\\nsaD\/rpnWPw2JGf7RCw9PYqXK1AIiGbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58\\nDzycJnO8Q0D8ViwQ8d5rE3SIoJdFoL\/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI\\n3WXyK8W\/jIYdnPFqi39\/bUXUYBWKmzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJ\\nQgGUpiuSSPHeXCDQiBvdaOLoPr4jR3F6exz39AByK5OkVwKENe9J\/tfZSVxkrG81\\nUd56\/Oal1jWJIiQHqCt7s1hMKInjKFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt\\n+QIDAQAB\\n-----END PUBLIC KEY-----\\n\",\r\n    \"private\": \"-----BEGIN RSA PRIVATE KEY-----\\nMIIEowIBAAKCAQEAqplqy+c2NbSGMuIRU8t8saD\/rpnWPw2JGf7RCw9PYqXK1AIi\\nGbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58DzycJnO8Q0D8ViwQ8d5rE3SIoJdF\\noL\/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI3WXyK8W\/jIYdnPFqi39\/bUXUYBWK\\nmzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJQgGUpiuSSPHeXCDQiBvdaOLoPr4j\\nR3F6exz39AByK5OkVwKENe9J\/tfZSVxkrG81Ud56\/Oal1jWJIiQHqCt7s1hMKInj\\nKFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt+QIDAQABAoIBABVxv3juEWRi0tOm\\nkyMDWyNA56Lc949pdihsXX6UaBgwWvSXaA3u1VuqylraP3i6U9zPZ1DP9vAql2zq\\nRjO59gI8TiyPM8UIcC+szlx45uDFLz9whHIWbvYT9I3bIkrLrNdmS+ubWtoocY\/e\\naVJOEugxnmVeMBvL6AEIX6o1VqE3h1BrAwLbDdP7T+muxJJC3wiXiSxqRe868AzV\\nc1eKQJjq+BTdV09bcfMTIZ7aNGgI6F1oZ\/NLI3UlwnOiLaWCb8aQyaLENwD0AGDw\\nX5k91OIvogM9cAhlXwidvyW99SuLuWdf\/n+FeXueqIf\/gnHu7BoFVi2uc3p9r4xK\\nF6dUJ+ECgYEA13fXqkoCoQZ2sAwM0gTRIXj0XomTmLBGtpil1P+1l7woUKw+sCxO\\n3oEjEovazEwyOY5bDYPMFqFtR0rNxp2YNPpo8k16XKI+p312p2BAtpbqLEImN\/tY\\n8idZfeClB1XFN8VSAC8OM3w30BHI+aHHMw32gI29ygApvxOWLKiTlWUCgYEAyrDa\\n4y4fc+ba0rKKCUoCHvJihiXPIuxfMPyVagCrRgr9WB5NNj4c1kKgHwDrwaz2xrtN\\nOlGAwn9X0i3e6bcfoQ0nRJslQbn66qqfjeNpqK6CEVv9DIgYVajKLDIkWhdQ6+Si\\nqn5Vq2MU6NIly34XBFoYfQTmRH0R7azdD+TfBwUCgYEArznt8LXJl4x7H0Zdcrqq\\nHI+SJAO8PZM1nq9bRXJDCsfg\/WJmhL0z0q2wiRelczmQKtCDaeVCJzFWfoDuAdUO\\nAB+ZE1xA426qh2l4Ajw7xIHMpPuSuzo0JpIrrDvx2Zo+DdHxkuaxpNsjRJoCGEkh\\nh3qWegtLSiiByruyCFV72CUCgYBgepxGBN9N0PYZ0ogn8cVeq6tABWE6U17gN2p7\\ngYQFHBgJSKsiBaC+UApdl5egoc75O5CAEOmEKw9HaTQw9Uyl4VfurRan2XnZF4xJ\\nApV5iE87KhkiTOmgZG6PaPKqu2x2TGctVmM66De8tsLswMD9\/lCnuZxNv2a4Rk8X\\nUK7kbQKBgECFgpN3sDketKz3DUa2oH6eLKHl1c1VWdnCKs7EJySn6p9nTqPg7B3J\\n8xe6VGuuU0vo2MqHuZJ+Oudpbz9iXpcyij6OcqCxgy8BV+yPV3WZ\/LNQ2fJbsQdS\\nBE6vbTy42rJAxgWLTkaJuDo7UFIpAw361R59n5nTIk5Mxtq3kIxa\\n-----END RSA PRIVATE KEY-----\\n\"\r\n  }"

lpush scope:openid id

```
#### Boostrap the CIBA server - Create the datastore
The CIBA server will need a persistence layer. `go-ciba`  provides what's called a datastore, an object that holds each repository respective of their vendor. Since this library has SQL and Redis out of the box, there will be a *SQLDataStore* and *RedisDataStore*. The naming convention has the vendor prefixed to *DataStore*.

Let's create a SQL datastore object.

```go
// replace this with your own credentials
db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=123 dbname=ciba sslmode=disable")
if err != nil {
    panic(err)
}
defer db.Close()

// third parameter is the prefix of the tables created
// since we didn't give it a prefix, we can pass in an
// empty string
// e.g. if we created the tables with 'my_app' prefix => my_app_access_tokens, 
// my_app_user_accounts etc we can pass in 'my_app' as the third parameter

ds := go_ciba.NewSQLDataStore(db, "postgres", "")

```

#### List of datastore initializers

Datastore objects must implement `DataStoreInterface` which is essentially a getter abstractions for the repositories it holds.

**Method: NewSQLDataStore**

| Parameters        | Description                                               |
|-------------------|-----------------------------------------------------------|
| defaultDb *sql.DB | The database connection                                   |
| driverName string | Driver of the database                                    |
| prefix string     | Prefix name of the tables, leave as empty string for none |

| Return type   | Description                                           |
|---------------|-------------------------------------------------------|
| *SQLDataStore | SQL datastore object which holds all the repositories |

**Method: NewRedisDataStore**

| Parameters           | Description          |
|----------------------|----------------------|
| client *redis.Client | The Redis connection |

| Return type     | Description                                             |
|-----------------|---------------------------------------------------------|
| *RedisDataStore | Redis datastore object which holds all the repositories |



#### Boostrap the CIBA server - Create the server objects

Once we have the datastore initialized, it can be used by the server objects. The server will use the repositories to gain access to the datalayer.

Let's create the CIBA server configuration.

```go
cibaGrant := grant.NewCustomCibaGrant(&grant.GrantConfig{
    Issuer:                       "auth.ciba.com",
    IdTokenLifetimeInSeconds:     3600,
    AccessTokenLifetimeInSeconds: 3600,
    PollingIntervalInSeconds:     &pollIntervalInSeconds,
    AuthReqIdLifetimeInSeconds:   120,
    TokenEndpointUrl:             "/token",
})
```

**Method: NewCustomCibaGrant**

| Parameters   | Description            |
|--------------|------------------------|
| *GrantConfig | The CIBA configuration |

**Properties in GrantConfig**

| Properties                         | Description                                                                                                                                                                                                                         |
|------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Issuer string                      | The identifier of the authorization server. It can be a URI. It will be the value of `iss` claim in the ID token                                                                                                                    |
| IdTokenLifetimeInSeconds int64     | The ID token lifetime in seconds until it expires                                                                                                                                                                                   |
| AccessTokenLifetimeInSeconds int64 | The access token lifetime in seconds until it expires                                                                                                                                                                               |
| PollingIntervalInSeconds *int64    |  The polling interval in seconds that the server will accept in `poll` mode. Clients polling faster than the specified amount will get the `slow_down` error. This parameter should be non null if the server supports `poll` mode. |
| AuthReqIdLifetimeInSeconds int64   | The authentication request ID lifetime in seconds until it expires                                                                                                                                                                  |
| TokenEndpointUrl string            | The URI of the token endpoint. This will be used in authenticating clients in `client_secret_jwt` method. Currently, `client_secret_jwt` method is not yet supported.                                                               |

----

Let's create the CIBA service object. The CIBA service will hold the logic to perform tasks such as handling authentication and consent requests. As you can see, we're passing in the repositories from the datastore we made earlier.

This library uses Firebase Cloud Messaging (FCM) to send notifications to Authentication Devices, a decoupled device possessed by the end-user to *give consent*. The way FCM is leveraged is by publishing to a topic with the user identifier. Therefore, our server must also register the topic of each user. This is implementation specific, but it can be done on each user login / registration.

```go
cibaService := gocibaService.NewCibaService(
    dataStore.GetClientApplicationRepository(),
    dataStore.GetUserAccountRepository(),
    dataStore.GetCibaSessionRepository(),
    dataStore.GetKeyRepository(),
    dataStore.GetUserClaimRepository(),
    gocibaTransport.NewFirebaseCloudMessaging(fcmServerKey),
    cibaGrant,
    func(token string) bool {
        return token != ""
    },
)

authorizationServer := gociba.NewAuthorizationServer(dataStore)
authorizationServer.AddService(cibaService)
```

**Method: NewCibaService**

| Parameters                                                    | Description                                                                                                                                                                                        |
|---------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| clientAppRepo ClientApplicationRepositoryInterface            | Client application repository                                                                                                                                                                      |
| userAccountRepo UserAccountRepositoryInterface                | User account repository                                                                                                                                                                            |
| cibaSessionRepo CibaSessionRepositoryInterface                | CIBA session repository                                                                                                                                                                            |
| keyRepo KeyRepositoryInterface                                | Key repository                                                                                                                                                                                     |
| userClaimRepo UserClaimRepositoryInterface                    | User claim repository                                                                                                                                                                              |
| notificationClient NotificationInterface                      | HTTP client to send notification to Authentication Device                                                                                                                                          |
| cibaGrant *CibaGrant                                          | CIBA config                                                                                                                                                                                        |
| validateClientNotificationToken  func ( token  string )  bool | Function to validate the client notification token sent by the client. Clients sends this in `ping` and `push` mode. Return `true` if the token conforms to specification, `false` in the contrary |

----

Let's create the token service object. This will hold logic to handle granting access and ID tokens.

```go
tokenService := gocibaService.NewTokenService(
  dataStore.GetAccessTokenRepository(),
  dataStore.GetClientApplicationRepository(),
  dataStore.GetCibaSessionRepository(),
  dataStore.GetKeyRepository(),
  dataStore.GetUserClaimRepository(),
  cibaGrant,
)

tokenServer := gociba.NewTokenServer(tokenService)
```


**Method: NewTokenService**

| Parameters                                         | Description                   |
|----------------------------------------------------|-------------------------------|
| accessTokenRepo AccessTokenRepositoryInterface     | Access token repository       |
| clientAppRepo ClientApplicationRepositoryInterface | Client application repository |
| cibaSessionRepo CibaSessionRepositoryInterface     | CIBA session repository       |
| keyRepo KeyRepositoryInterface                     | Key repository                |
| userClaimRepo UserClaimRepositoryInterface         | User claim repository         |
| grant *CibaGrant                                   | CIBA config                   |

---

Let's create the resource server. This will hold logic to protect non-public resources by the scope it was assigned to.


```go
resourceServer := gociba.NewResourceServer(dataStore.GetAccessTokenRepository())
```


#### Putting everything together

Once we have the building blocks done, we can use it in our HTTP handlers. We'll be using the gin library as an example, but it can be used in any HTTP router library.

```go
r.POST("/auth", func(context *gin.Context) {
    req := gocibaService.NewAuthenticationRequest(context.Request)
    req.ValidateBindingMessage = func(bindingMessage string) bool {
        return true
    }
    req.ValidateUserCode = func(code, givenCode string) bool {
        return true
    }
    res, err := authorizationServer.HandleCibaRequest(req)
    if err != nil {
        context.JSON(err.Code, err)
        return
    }
    context.JSON(http.StatusOK, res)
})

r.POST("/consent", func(context *gin.Context) {
    authReqId := context.PostForm("auth_req_id")
    consented := context.PostForm("consented") == "true"
    req := gocibaService.NewConsentRequest(authReqId, &consented)

    err := authorizationServer.HandleConsentRequest(req)
    if err != nil {
        context.JSON(err.Code, err)
        return
    }
    context.JSON(http.StatusOK, req)
})

r.POST("/token", func(context *gin.Context) {
    req := gocibaService.NewTokenRequest(context.Request)
    res, err := tokenServer.HandleTokenRequest(req)
    if err != nil {
        context.JSON(err.Code, err)
        return
    }
    context.JSON(http.StatusOK, res)
})

r.POST("/protected", func(c *gin.Context) {
    req := gociba.NewResourceRequest(c.Request)
    err := resourceServer.HandleResourceRequest(req, "timestamp.read")
    if err != nil {
        c.JSON(err.Code, err)
        return
    }
    c.JSON(http.StatusOK, "In protected resource")
})
```


## Authors 
- [Adis Azhar](https://id.linkedin.com/in/adis-azhar-33216a15a)

## Version History

* 0.1
    * Initial Release

## Acknowledgments

Inspiration, code snippets, etc.
* [bshaffer/oauth2-server-php](https://github.com/bshaffer/oauth2-server-php)
* [ory/fosite](https://github.com/ory/fosite)
* [OpenID Connect CIBA draft](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html) - at the time of this implementation, draft 03 was used