package builder

import "strings"

type InsertBuilder struct {
	table   string
	columns []string
	values  []string
}

func NewInsertBuilder() *InsertBuilder {
	return &InsertBuilder{
		values: make([]string, 0),
	}
}

func (i *InsertBuilder) Table(name string) *InsertBuilder {
	i.table = name
	return i
}

func (i *InsertBuilder) Values(values ...string) *InsertBuilder {
	i.values = append(i.values, values...)
	return i
}

func (i *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	i.columns = append(i.columns, columns...)
	return i
}

func (i *InsertBuilder) ToSQL() (string, interface{}) {
	var q strings.Builder

	q.WriteString("INSERT INTO ")
	q.WriteString(i.table)

	q.WriteString("(")
	for k := range i.columns {
		q.WriteString(i.columns[k])
		if k != len(i.values)-1 {
			q.WriteString(",")
		}
	}

	q.WriteString(")")
	q.WriteString("VALUES (")
	for k := range i.values {
		q.WriteString(i.values[k])
		if k != len(i.values)-1 {
			q.WriteString(",")
		}
	}
	q.WriteString(")")
	return q.String(), nil
}

/*
*	INSERT INTO table (columsns ....)
*	VALUES (data for the columns...)
* */
