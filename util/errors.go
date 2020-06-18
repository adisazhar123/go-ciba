package util

import "net/http"

const (
	errAuthorizationPending = "authorization_pending"
	errSlowDown = "slow_down"
	errExpiredToken = "expired_token"
	errAccessDenied = "access_denied"
	errUnauthorizedClient = "unauthorized_client"
	errInvalidRequest = "invalid_request"
	errInvalidGrant = "invalid_grant"
	errTransactionFailed = "transaction_failed"
	errInvalidScope = "invalid_scope"
	errExpiredLoginHintToken = "expired_login_hint_token"
	errUnknownUserId = "unknown_user_id"
	errMissingUserCode = "missing_user_code"
	errInvalidUserCode = "invalid_user_code"
	errInvalidBindingMessage = "invalid_binding_message"
	errInvalidClient = "invalid_client"
)

var (
	ErrAuthorizationPending = &OidcError{
		Error:            errAuthorizationPending,
		ErrorDescription: "The authorization request is still pending as the end-user hasn't yet been authenticated.",
		Code:             http.StatusBadRequest,
	}
	ErrSlowDown = &OidcError{
		Error:            errSlowDown,
		ErrorDescription: "The token request is too fast.",
		Code:             http.StatusBadRequest,
	}
	ErrExpiredToken = &OidcError{
		Error:            errExpiredToken,
		ErrorDescription: "The auth_req_id has expired.",
		Code:             http.StatusUnauthorized,
	}
	ErrAccessDenied = &OidcError{
		Error:            errAccessDenied,
		ErrorDescription: "The end-user denied the authorization request.",
		Code:             http.StatusForbidden,
	}
	ErrInvalidGrant = &OidcError{
		Error:            errInvalidGrant,
		ErrorDescription: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
		Code:             http.StatusBadRequest,
	}
	ErrUnauthorizedClient = &OidcError{
		Error:            errUnauthorizedClient,
		ErrorDescription: "The Client is not authorized to use this authentication flow.",
		Code:             http.StatusBadRequest,
	}
	ErrInvalidRequest = &OidcError{
		Error:            errInvalidRequest,
		ErrorDescription: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, contains more than one of the hints, or is otherwise malformed.",
		Code:             http.StatusBadRequest,
	}
	ErrTransactionFailed = &OidcError{
		Error:            errTransactionFailed,
		ErrorDescription: "The OpenID Provider encountered an unexpected condition that prevented it from successfully completing the transaction.",
		Code:             http.StatusBadRequest,
	}
	ErrInvalidScope = &OidcError{
		Error:            errInvalidScope,
		ErrorDescription: "The requested scope is invalid, unknown, or malformed.",
		Code:             http.StatusBadRequest,
	}
	ErrExpiredLoginHintTOken = &OidcError{
		Error:            errExpiredLoginHintToken,
		ErrorDescription: "The login_hint_token provided in the authentication request is not valid because it has expired.",
		Code:             http.StatusBadRequest,
	}
	ErrUnknownUserId = &OidcError{
		Error:            errUnknownUserId,
		ErrorDescription: "The OpenID Provider is not able to identify which end-user the Client wishes to be authenticated by means of the hint provided in the request (login_hint_token, id_token_hint or login_hint).",
		Code:             http.StatusBadRequest,
	}
	ErrMissingUserCode = &OidcError{
		Error:            errMissingUserCode,
		ErrorDescription: "User code is required but was missing from the request.",
		Code:             http.StatusBadRequest,
	}
	ErrInvalidUserCode = &OidcError{
		Error:            errInvalidUserCode,
		ErrorDescription: "User code was invalid",
		Code:             http.StatusBadRequest,
	}
	ErrInvalidBindingMessage = &OidcError{
		Error:            errInvalidBindingMessage,
		ErrorDescription: "The binding message is invalid or unacceptable for use in the context of the given request.",
		Code:             http.StatusBadRequest,
	}
	ErrInvalidClient = &OidcError{
		Error:            errInvalidClient,
		ErrorDescription: "Client authentication failed (e.g., invalid client credentials, unknown client, no client authentication included, or unsupported authentication method).",
		Code:             http.StatusUnauthorized,
	}
)

type OidcError struct {
	Error string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri string `json:"error_uri,omitempty"`
	Code int `json:"status_code,omitempty"`
}