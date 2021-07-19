# go-ciba


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

