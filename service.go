package rest

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(data interface{}, p Param) (interface{}, *ServiceError)
	Find(p Param) (interface{}, *ServiceError)
	Get(id string, p Param) (interface{}, *ServiceError)
	Patch(id string, p Param) (interface{}, *ServiceError)
	Remove(id string, p Param) (interface{}, *ServiceError)
	Update(id string, p Param) (interface{}, *ServiceError)

	Prepare(ctx *gin.Context)
}

type CreateHandler func(data interface{}, p Param) (interface{}, *ServiceError)
type FindHandler func(p Param) (interface{}, *ServiceError)
type GetHandler func(id string, p Param) (interface{}, *ServiceError)
type PatchHandler func(id string, p Param) (interface{}, *ServiceError)
type RemoveHandler func(id string, p Param) (interface{}, *ServiceError)
type UpdateHandler func(id string, p Param) (interface{}, *ServiceError)

type ServiceBuilder struct {
	create CreateHandler
	find   FindHandler
	get    GetHandler
	patch  PatchHandler
	remove RemoveHandler
	update UpdateHandler
}

func (b ServiceBuilder) Create(handle CreateHandler) ServiceBuilder {
	b.create = handle
	return b
}

func (b ServiceBuilder) Find(handle FindHandler) ServiceBuilder {
	b.find = handle
	return b
}

func (b ServiceBuilder) Service() (service Service, allowedMethods []int) {
	service = AnonymousService{b}
	allowedMethods = []int{}

	if b.get != nil {
		allowedMethods = append(allowedMethods, MethodGet)
	}

	if b.find != nil {
		allowedMethods = append(allowedMethods, MethodFind)
	}

	if b.create != nil {
		allowedMethods = append(allowedMethods, MethodCreate)
	}

	if b.remove != nil {
		allowedMethods = append(allowedMethods, MethodRemove)
	}

	if b.patch != nil {
		allowedMethods = append(allowedMethods, MethodPatch)
	}

	if b.update != nil {
		allowedMethods = append(allowedMethods, MethodUpdate)
	}

	return
}

func NewService() *ServiceBuilder {
	return &ServiceBuilder{}
}

type AnonymousService struct {
	serviceBuilder ServiceBuilder
}

func (a AnonymousService) Patch(id string, p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.patch != nil {
		return a.serviceBuilder.patch(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Remove(id string, p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.remove != nil {
		return a.serviceBuilder.remove(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Update(id string, p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.update != nil {
		return a.serviceBuilder.update(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Create(data interface{}, p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.create != nil {
		return a.serviceBuilder.create(data, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Find(p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.find != nil {
		return a.serviceBuilder.find(p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Get(id string, p Param) (interface{}, *ServiceError) {
	if a.serviceBuilder.get != nil {
		return a.serviceBuilder.get(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Prepare(ctx *gin.Context) {}
