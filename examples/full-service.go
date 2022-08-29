package examples

import (
	rest "github.com/blackmann/resteasy"
	"github.com/gin-gonic/gin"
)

func usage() {
	router := gin.Default()
	rest.With(newCustomService()).Register(router.Group("/items"))

	_ = router.Run()
}

type custom struct {
}

func (c custom) Create(data interface{}, p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Find(p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Get(id string, p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Patch(id string, data interface{}, p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Remove(id string, p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Update(id string, data interface{}, p rest.Params) (interface{}, *rest.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (c custom) Prepare(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func newCustomService() custom {
	return custom{}
}

