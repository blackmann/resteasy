package resteasy

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"io"
)

type Service interface {
	Create(data interface{}, p Params) (interface{}, *ServiceError)
	Find(p Params) (interface{}, *ServiceError)
	Get(id string, p Params) (interface{}, *ServiceError)
	Patch(id string, data interface{}, p Params) (interface{}, *ServiceError)
	Remove(id string, p Params) (interface{}, *ServiceError)
	Update(id string, data interface{}, p Params) (interface{}, *ServiceError)

	Prepare(ctx *gin.Context)
}

type CreateHandler func(data interface{}, p Params) (interface{}, *ServiceError)
type FindHandler func(p Params) (interface{}, *ServiceError)
type GetHandler func(id string, p Params) (interface{}, *ServiceError)
type PatchHandler func(id string, data interface{}, p Params) (interface{}, *ServiceError)
type RemoveHandler func(id string, p Params) (interface{}, *ServiceError)
type UpdateHandler func(id string, data interface{}, p Params) (interface{}, *ServiceError)

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

func (b ServiceBuilder) Get(handle GetHandler) ServiceBuilder {
	b.get = handle
	return b
}

func (b ServiceBuilder) Patch(handle PatchHandler) ServiceBuilder {
	b.patch = handle
	return b
}

func (b ServiceBuilder) Remove(handle RemoveHandler) ServiceBuilder {
	b.remove = handle
	return b
}

func (b ServiceBuilder) Update(handle UpdateHandler) ServiceBuilder {
	b.update = handle
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

func NewService() ServiceBuilder {
	return ServiceBuilder{}
}

type AnonymousService struct {
	serviceBuilder ServiceBuilder
}

func (a AnonymousService) Patch(id string, data interface{}, p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.patch != nil {
		return a.serviceBuilder.patch(id, data, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Remove(id string, p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.remove != nil {
		return a.serviceBuilder.remove(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Update(id string, data interface{}, p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.update != nil {
		return a.serviceBuilder.update(id, data, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Create(data interface{}, p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.create != nil {
		return a.serviceBuilder.create(data, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Find(p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.find != nil {
		return a.serviceBuilder.find(p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Get(id string, p Params) (interface{}, *ServiceError) {
	if a.serviceBuilder.get != nil {
		return a.serviceBuilder.get(id, p)
	}

	return nil, MethodNotAllowed()
}

func (a AnonymousService) Prepare(ctx *gin.Context) {
	if slices.Contains([]string{"POST", "PATCH", "PUT"}, ctx.Request.Method) {
		bytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}

		var data interface{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			panic(err)
		}

		ctx.Set("data", data)
	}
}
