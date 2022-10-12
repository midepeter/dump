package builder

type insertBuilder struct {
	table string
	values []string
}

func NewInsertBuilder () Generator {
	return &insertBuilder{
		values: make([]string, 0)
	}
}

j := NewJoiner(, "", "", "", true)

func (i insertBuilder) SQL() (string, interface{}) {
	return "", nil
}
