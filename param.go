package rest

type Param map[string]interface{}

func (p Param) Get(key string) interface{} {
	return p[key]
}

func (p Param) Set(key string, value interface{}) {
	p[key] = value
}
