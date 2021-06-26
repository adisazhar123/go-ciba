# go-ciba


## Description

An in-depth paragraph about your project and overview of use.

## Getting Started



### Dependencies

* Describe any prerequisites, libraries, OS version, etc., needed before installing program.
* ex. Windows 10

### Installing

* How/where to download your program
* Any modifications needed to be made to files/folders

### Executing program

* How to run the program
* Step-by-step bullets
```
code blocks for commands
```

## Help

Any advise for common problems or issues.
```
command to run if program contains helper info
```

## Authors

Contributors names and contact info

ex. Dominique Pizzie  
ex. [@DomPizzie](https://twitter.com/dompizzie)

## Version History

* 0.2
    * Various bug fixes and optimizations
    * See [commit change]() or See [release history]()
* 0.1
    * Initial Release

## License

This project is licensed under the [NAME HERE] License - see the LICENSE.md file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [bshaffer/oauth2-server-php](https://github.com/bshaffer/oauth2-server-php)
* [ory/fosite](https://github.com/ory/fosite)


### SQL Table Structure
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
    Token_endpoint_auth_signing_alg VARCHAR(10),
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
```