package jinx

type JinxEntry interface{}

type databaseEntry struct {
	Value          interface{}
	expirationDate int64
}
