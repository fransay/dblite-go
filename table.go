package dblite

type ITable[T any] interface {
	Clone() T
	TableName() string
	Fields() []string
	FilterFieldReferences(fields []string) ([]string, []any, error)
}
