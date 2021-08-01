package domain

type Scope struct {
	Name string `db:"name" json:"name"`
}

type Claim struct {
	Name      string `db:"name" json:"name"`
	ScopeName string
}

func (c *Claim) Str() string {
	return c.Name
}
