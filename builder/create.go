package builder

import (
	"strconv"
	"strings"
)

type CreateBuilder struct {
	table       string
	columns     []string
	values      []string
	ifNotExists bool
}

func NewCreateBuilder() *CreateBuilder {
	return &CreateBuilder{
		columns: make([]string, 0),
	}
}

func (c *CreateBuilder) Table(t string) *CreateBuilder {
	c.table = t
	return c
}

func (c *CreateBuilder) Columns(column ...string) *CreateBuilder {
	c.columns = append(c.columns, column...)
	return c
}

func (c *CreateBuilder) IfNotExists(exist bool) *CreateBuilder {
	if exist {
		c.ifNotExists = true
	}
	return c
}

//first line of value to determine the type
func (c *CreateBuilder) Values(values ...string) *CreateBuilder {
	c.values = append(c.values, values...)
	return c
}

func (c *CreateBuilder) ToSQL() (string, interface{}) {
	var q strings.Builder

	q.WriteString("CREATE TABLE ")
	q.WriteString(c.table)
	if !c.ifNotExists {
		q.WriteString("IF NOT EXISTS (")
	}

	for i := 0; i < len(c.columns); i++ {
		q.WriteString(c.columns[i])

		switch c.values[i] {
		case "true", "false":
			q.WriteString("BIT ")
		default:
			_, err := strconv.Atoi(c.values[i])
			if err == nil {
				q.WriteString("INT")
				break
			}
			q.WriteString("VARCHAR(50)")
		}
	}

	q.WriteString(")")
	return q.String(), nil
}

/*
*	CREATE TABLE user IF NOT EXISTS (
	*		name VARCHAR(50)
	*		email VARCHAR(50)
	*		password VARCHAR(50)
	*	)
**/
