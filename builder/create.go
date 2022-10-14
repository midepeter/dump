package builder

type createBuilder struct {
	table       string
	columns     []string
	ifNotExists bool
}

func NewCreateBuilder() QueryBuilder {
	return &createBuilder{
		columns: make([]string, 0),
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

func (c *createBuilder) ToSQL() (string, interface{}) {
	return "", nil
}

/*
*	CREATE TABLE user IF NOT EXISTS (
	*		name VARCHAR(50)
	*		email VARCHAR(50)
	*		password VARCHAR(50)
	*	)
**/
