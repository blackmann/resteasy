package resteasy

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerBuilder struct {
	service        Service
	allowedMethods []int
	beforeHooks    []gin.HandlerFunc
	afterHooks     []gin.HandlerFunc
}

func GetParams(ctx *gin.Context) Params {
	params, exists := ctx.Get("params")
	if !exists {
		params = Params{}
		ctx.Set("params", params)
	}

	return params.(Params)
}

func getData(ctx *gin.Context) interface{} {
	var data, exists = ctx.Get("data")

	if !exists {
		panic("`data` does not exist in the context. Make sure a middleware is setting this value")
	}

	return data
}

func (b HandlerBuilder) find(route *gin.RouterGroup) {
	route.GET("", func(ctx *gin.Context) {
		params := GetParams(ctx)
		response, err := b.service.Find(params)
		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b HandlerBuilder) get(route *gin.RouterGroup) {
	route.GET("/:id", func(ctx *gin.Context) {
		params := GetParams(ctx)
		id := ctx.Param("id")
		response, err := b.service.Get(id, params)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b HandlerBuilder) create(route *gin.RouterGroup) {
	route.POST("", func(ctx *gin.Context) {
		data := getData(ctx)
		params := GetParams(ctx)
		response, err := b.service.Create(data, params)
		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusCreated, response)
	})
}

func (b HandlerBuilder) patch(route *gin.RouterGroup) {
	route.PATCH("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		params := GetParams(ctx)

		data := getData(ctx)
		response, err := b.service.Patch(id, data, params)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b HandlerBuilder) remove(route *gin.RouterGroup) {
	route.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		params := GetParams(ctx)
		response, err := b.service.Remove(id, params)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b HandlerBuilder) update(route *gin.RouterGroup) {
	route.PUT("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		data := getData(ctx)
		params := GetParams(ctx)
		response, err := b.service.Update(id, data, params)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b HandlerBuilder) Register(route *gin.RouterGroup) {
	route.Use(b.service.Prepare)

	route.Use(func(ctx *gin.Context) {
		for _, hook := range b.beforeHooks {
			if !ctx.IsAborted() {
				hook(ctx)
			}
		}

		ctx.Next()

		for _, hook := range b.afterHooks {
			if !ctx.IsAborted() {
				hook(ctx)
			}
		}
	})

	var methodMap = map[int]func(r *gin.RouterGroup){
		MethodGet:    b.get,
		MethodFind:   b.find,
		MethodCreate: b.create,
		MethodPatch:  b.patch,
		MethodRemove: b.remove,
		MethodUpdate: b.update,
	}

	for _, method := range b.allowedMethods {
		methodMap[method](route)
	}
}

func (b HandlerBuilder) Before(hooks ...gin.HandlerFunc) HandlerBuilder {
	b.beforeHooks = append(b.beforeHooks, hooks...)
	return b
}

func (b HandlerBuilder) After(hooks ...gin.HandlerFunc) HandlerBuilder {
	b.afterHooks = append(b.afterHooks, hooks...)
	return b
}

func With(service Service, allowedMethods ...int) HandlerBuilder {
	if len(allowedMethods) == 0 {
		allowedMethods = AllMethods
	}

	return HandlerBuilder{service: service, allowedMethods: allowedMethods}
}
