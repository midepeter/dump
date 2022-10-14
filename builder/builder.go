package builder

type QueryBuilder interface {
	ToSQL() (string, interface{})
}
