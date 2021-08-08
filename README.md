# go-ciba

[![adisazhar123](https://circleci.com/gh/adisazhar123/go-ciba/tree/development.svg?style=svg&circle-token=7ce9bd81b2ad605ae129d574f18fe39aee783bab)](https://app.circleci.com/pipelines/github/adisazhar123/go-ciba?branch=development)

[![codecov](https://codecov.io/gh/adisazhar123/go-ciba/branch/development/graph/badge.svg?token=D9648NK5HO)](https://codecov.io/gh/adisazhar123/go-ciba)
## Description

An in-depth paragraph about your project and overview of use.

## Getting Started



### Dependencies

* Go 1.13 or newer

### Step-By-Step Walkthrough

The following instructions will provide you to get this library up and running.

#### Initialize your project

#### Define your schema
The `go-ciba` library is storage agnostic which means that it's not tied to a vendor specific database. It can be plug and played with any database by implementing the interfaces in `repository/repo.go`. 

As of now, it comes with a prebuilt SQL and Redis implementation. For the sake of getting it up and running, we'll use the SQL implementation. 

Use the following schema to create the database.
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

```shell
set client_application:b4620189-c368-43ed-b2b4-2186a61fa664 "{\r\n \"id\": \"b4620189-c368-43ed-b2b4-2186a61fa664\",\r\n  \"secret\": \"83e34759-314e-45ec-8211-c6869e053187\",\r\n  \"name\": \"My First Client\",\r\n  \"scope\": \"openid\",\r\n  \"token_mode\": \"poll\",\r\n  \"client_notification_endpoint\": \"\",\r\n  \"authentication_request_signing_alg\": \"\",\r\n  \"user_code_parameter_supported\": false,\r\n  \"redirect_uri\": \"\",\r\n  \"token_endpoint_auth_method\":\"client_secret_basic\",\r\n  \"token_endpoint_auth_signing_alg\": \"\",\r\n  \"grant_types\": \"urn:openid:params:grant-type:ciba\",\r\n  \"public_key_uri\": \"\"\r\n}"

set user_account:f24e0c6d-dbf0-4753-87ad-b554aab423a5 "{\r\n    \"id\": \"f24e0c6d-dbf0-4753-87ad-b554aab423a5\",\r\n    \"name\": \"Joe Foo\",\r\n    \"email\": \"joe@foo.com\",\r\n    \"password\": \"secret\",\r\n    \"user_code\": \"\",\r\n    \"created_at\": \"2021-01-01\",\r\n    \"updated_at\": \"2021-01-01\"\r\n  }"

set oauth_key:2b075e6c-790c-4ea1-a697-52a382bec9b7 "{\r\n    \"id\": \"2b075e6c-790c-4ea1-a697-52a382bec9b7\",\r\n    \"client_id\": \"b4620189-c368-43ed-b2b4-2186a61fa664\",\r\n    \"alg\":\"RSA256\",\r\n    \"public\": \"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqplqy+c2NbSGMuIRU8t8\\nsaD\/rpnWPw2JGf7RCw9PYqXK1AIiGbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58\\nDzycJnO8Q0D8ViwQ8d5rE3SIoJdFoL\/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI\\n3WXyK8W\/jIYdnPFqi39\/bUXUYBWKmzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJ\\nQgGUpiuSSPHeXCDQiBvdaOLoPr4jR3F6exz39AByK5OkVwKENe9J\/tfZSVxkrG81\\nUd56\/Oal1jWJIiQHqCt7s1hMKInjKFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt\\n+QIDAQAB\\n-----END PUBLIC KEY-----\\n\",\r\n    \"private\": \"-----BEGIN RSA PRIVATE KEY-----\\nMIIEowIBAAKCAQEAqplqy+c2NbSGMuIRU8t8saD\/rpnWPw2JGf7RCw9PYqXK1AIi\\nGbIqN1Gqx6XUNr+xKm0kHc9j4XggDfmCRL58DzycJnO8Q0D8ViwQ8d5rE3SIoJdF\\noL\/0dK+YoxMVwCt+kqZLgq5ZDBj521SADaeI3WXyK8W\/jIYdnPFqi39\/bUXUYBWK\\nmzA2FfA9ucM9idnxKPrXInjelmXd4VnUcXsJQgGUpiuSSPHeXCDQiBvdaOLoPr4j\\nR3F6exz39AByK5OkVwKENe9J\/tfZSVxkrG81Ud56\/Oal1jWJIiQHqCt7s1hMKInj\\nKFLvQIIdWMchpmfB+Gr67pTthCsFAWMDavKt+QIDAQABAoIBABVxv3juEWRi0tOm\\nkyMDWyNA56Lc949pdihsXX6UaBgwWvSXaA3u1VuqylraP3i6U9zPZ1DP9vAql2zq\\nRjO59gI8TiyPM8UIcC+szlx45uDFLz9whHIWbvYT9I3bIkrLrNdmS+ubWtoocY\/e\\naVJOEugxnmVeMBvL6AEIX6o1VqE3h1BrAwLbDdP7T+muxJJC3wiXiSxqRe868AzV\\nc1eKQJjq+BTdV09bcfMTIZ7aNGgI6F1oZ\/NLI3UlwnOiLaWCb8aQyaLENwD0AGDw\\nX5k91OIvogM9cAhlXwidvyW99SuLuWdf\/n+FeXueqIf\/gnHu7BoFVi2uc3p9r4xK\\nF6dUJ+ECgYEA13fXqkoCoQZ2sAwM0gTRIXj0XomTmLBGtpil1P+1l7woUKw+sCxO\\n3oEjEovazEwyOY5bDYPMFqFtR0rNxp2YNPpo8k16XKI+p312p2BAtpbqLEImN\/tY\\n8idZfeClB1XFN8VSAC8OM3w30BHI+aHHMw32gI29ygApvxOWLKiTlWUCgYEAyrDa\\n4y4fc+ba0rKKCUoCHvJihiXPIuxfMPyVagCrRgr9WB5NNj4c1kKgHwDrwaz2xrtN\\nOlGAwn9X0i3e6bcfoQ0nRJslQbn66qqfjeNpqK6CEVv9DIgYVajKLDIkWhdQ6+Si\\nqn5Vq2MU6NIly34XBFoYfQTmRH0R7azdD+TfBwUCgYEArznt8LXJl4x7H0Zdcrqq\\nHI+SJAO8PZM1nq9bRXJDCsfg\/WJmhL0z0q2wiRelczmQKtCDaeVCJzFWfoDuAdUO\\nAB+ZE1xA426qh2l4Ajw7xIHMpPuSuzo0JpIrrDvx2Zo+DdHxkuaxpNsjRJoCGEkh\\nh3qWegtLSiiByruyCFV72CUCgYBgepxGBN9N0PYZ0ogn8cVeq6tABWE6U17gN2p7\\ngYQFHBgJSKsiBaC+UApdl5egoc75O5CAEOmEKw9HaTQw9Uyl4VfurRan2XnZF4xJ\\nApV5iE87KhkiTOmgZG6PaPKqu2x2TGctVmM66De8tsLswMD9\/lCnuZxNv2a4Rk8X\\nUK7kbQKBgECFgpN3sDketKz3DUa2oH6eLKHl1c1VWdnCKs7EJySn6p9nTqPg7B3J\\n8xe6VGuuU0vo2MqHuZJ+Oudpbz9iXpcyij6OcqCxgy8BV+yPV3WZ\/LNQ2fJbsQdS\\nBE6vbTy42rJAxgWLTkaJuDo7UFIpAw361R59n5nTIk5Mxtq3kIxa\\n-----END RSA PRIVATE KEY-----\\n\"\r\n  }"

lpush scope:openid id

```
#### Boostrap the CIBA server
We need create and configure our authorization server.

`main.go`
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
// e.g. if we created the tables with 'my_app' prefix => my_app_access_tokens, my_app_user_accounts etc we can pass in 'my_app' as the third parameter

ds := go_ciba.NewSQLDataStore(db, "postgres", "")

```


## Help

Any advise for common problems or issues.
```
command to run if program contains helper info
```

## Authors 
- [Adis Azhar](https://github.com/adisazhar123)

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the [NAME HERE] License - see the LICENSE.md file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [bshaffer/oauth2-server-php](https://github.com/bshaffer/oauth2-server-php)
* [ory/fosite](https://github.com/ory/fosite)

