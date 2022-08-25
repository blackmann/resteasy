# rest

This framework is based on [gin](https://github.com/gin-gonic/gin) and inspired
by [FeathersJS](https://github.com/feathersjs/feathers). That is, this library/framework allows you to develop REST
services with a rapid workflow. The API is easy to follow/learn.

## Examples

```go
package demo

import (
	"github.com/blackmann/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	group := router.Group("/demo")

	findHandler := func(p rest.Param) (interface{}, *rest.ServiceError) {
		return []int{1, 2, 3}, nil
	}

	getHandler := func(id string, p rest.Param) (interface{}, *rest.ServiceError) {
		return map[string]string{"title": "Hello world"}, nil
	}

	service, allowedMethods := rest.NewService().
		Find(findHandler). // Resolves on /demo
		Get(getHandler). // Resolves /demo/123
		Service()

	rest.With(service, allowedMethods...).Register(group)

	_ = router.Run()
}

```
