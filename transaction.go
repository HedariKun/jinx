package jinx

type TransactionHandler func(tx *JinxTransaction) error

type JinxTransaction struct {
	localDB  *JinxDatabase
	parentDB *JinxDatabase
}

func (jt *JinxTransaction) Set(key string, value interface{}) error {
	return jt.localDB.Set(key, value)
}

func (jt *JinxTransaction) SetExpire(key string, value interface{}, expire int) error {
	return jt.localDB.SetExpire(key, value, expire)
}

func (jt *JinxTransaction) SetMap(key string, m map[string]string) error {
	return jt.localDB.SetMap(key, m)
}

func (jt *JinxTransaction) Get(key string) interface{} {
	element := jt.localDB.Get(key)
	if element != nil {
		return element
	}
	return jt.parentDB.Get(key)
}

func (jt *JinxTransaction) Delete(key string) {
	jt.localDB.Delete(key)
}

func (jt *JinxTransaction) Exists(key string) bool {
	return jt.localDB.Exists(key)
}
