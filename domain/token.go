package domain

type AccessToken struct {
	value     string
	tokenType string
	expiresIn string
}

type IdToken struct {
	value string
}
