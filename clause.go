package dblite

type On struct {
	On            string
	UpsertColumns []string
	Arguments     []any
}

type WhereClause struct {
	Where     string
	Arguments []any
}
