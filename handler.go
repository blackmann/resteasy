package rest

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

func getParam(ctx *gin.Context) Param {
	p, exists := ctx.Get("params")
	if !exists {
		p = Param{}
	}

	return p.(Param)
}

func (b *HandlerBuilder) find(route *gin.RouterGroup) {
	route.GET("", func(ctx *gin.Context) {
		p := getParam(ctx)
		response, err := b.service.Find(p)
		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b *HandlerBuilder) get(route *gin.RouterGroup) {
	route.GET("/:id", func(ctx *gin.Context) {
		p := getParam(ctx)
		id := ctx.Query("id")
		response, err := b.service.Get(id, p)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b *HandlerBuilder) create(route *gin.RouterGroup) {
	route.POST("", func(ctx *gin.Context) {
		var data, exists = ctx.Get("data")

		if !exists {
			panic("`data` does not exist in the context. Make sure a middleware is setting this value")
		}

		p := getParam(ctx)
		response, err := b.service.Create(data, p)
		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusCreated, response)
	})
}

func (b *HandlerBuilder) patch(route *gin.RouterGroup) {
	route.PATCH("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		param := getParam(ctx)
		response, err := b.service.Patch(id, param)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b *HandlerBuilder) remove(route *gin.RouterGroup) {
	route.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Query("id")
		param := getParam(ctx)
		response, err := b.service.Remove(id, param)

		if err != nil {
			ctx.JSON(err.Code, err.Detail)
			return
		}

		ctx.JSON(http.StatusOK, response)
	})
}

func (b *HandlerBuilder) update(route *gin.RouterGroup) {
	route.PUT("/:id", func(context *gin.Context) {

	})
}

func (b *HandlerBuilder) Register(route *gin.RouterGroup) {
	route.Use(b.service.Prepare)

	route.Use(func(ctx *gin.Context) {
		for _, hook := range b.beforeHooks {
			hook(ctx)
		}

		ctx.Next()

		for _, hook := range b.afterHooks {
			hook(ctx)
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

func (b *HandlerBuilder) Before(hooks ...gin.HandlerFunc) *HandlerBuilder {
	b.beforeHooks = append(b.beforeHooks, hooks...)
	return b
}

func With(service Service, allowedMethods ...int) *HandlerBuilder {
	if len(allowedMethods) == 0 {
		allowedMethods = AllMethods
	}

	return &HandlerBuilder{service: service, allowedMethods: allowedMethods}
}
