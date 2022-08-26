package rest

type Params map[string]interface{}

func (p Params) Get(key string) interface{} {
	return p[key]
}

func (p Params) Set(key string, value interface{}) {
	p[key] = value
}
