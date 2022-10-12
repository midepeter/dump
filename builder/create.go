package builder

type createBuilder struct {
	table       string
	columns     []string
	ifNotExists bool
}

func NewCreateBuilder() Generator {
	return &insertBuilder{
		values: make([]string, 0),
	}
}

func (c *createBuilder) Table(t string) *createBuilder {
	c.table = t
	return c
}

func (c *createBuilder) Columns(column ...string) *createBuilder {
	c.columns = append(c.columns, column...)
	return c
}

func (c *createBuilder) IfNotExists(exist bool) *createBuilder {
	if exist {
		c.ifNotExists = true
	}
	return c
}

func (c *createBuilder) SQL() (string, interface{}) {
	join := NewJoiner([]Ge{}, " ", "", "", true)
	return "", nil
}

/* 
*	CREATE TABLE user IF NOT EXISTS (
	*		name VARCHAR(50)
	*		email VARCHAR(50)
	*		password VARCHAR(50)
	*	)
**/
