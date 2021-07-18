package domain

type Scope struct {
	Name string `db:"name"`
}

type Claim struct {
	Name string `db:"name"`
	ScopeName string
}

func (c *Claim) Str() string {
	return c.Name
}