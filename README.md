# resteasy

This framework is based on [gin](https://github.com/gin-gonic/gin) and inspired
by [FeathersJS](https://github.com/feathersjs/feathers). That is, this library/framework allows you to develop REST
services with a rapid workflow. The API is easy to follow/learn.

## Examples

```go
package demo

import (
	rest "github.com/blackmann/resteasy"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	group := router.Group("/demo")

	findHandler := func(p rest.Params) (interface{}, *rest.ServiceError) {
		return []int{1, 2, 3}, nil
	}

	getHandler := func(id string, p rest.Params) (interface{}, *rest.ServiceError) {
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

That was a quick and simple demonstration on how to use `resteasy`. Look inside the [examples/](/examples) folder for
more.

## Concepts

### Methods

Following similar convention from FeathersJS, methods are nicknamed as the following:

| Nickname | Method/Path        |                                             |
|----------|--------------------|:--------------------------------------------|
| find     | `GET` /demo        | The index path, where you return all items. |
| get      | `GET` /demo/:id    | Get a single item                           |
| create   | `POST` /demo/      |                                             |
| patch    | `PATCH` /demo/:id  | Patch a single item                         |
| update   | `PUT` /demo/:id    | Replace a single item                       |
| remove   | `DELETE` /demo/:id |                                             |

These may feel foreign at first sight, but they make it easy for you to implement services that are not directly coupled
to HTTP request methods. This also allows to test your service implementations in isolation from HTTP request
processing.

### Hooks

Hooks help perform some actions before or after a request is handled by the _service_. For example, you may want to
authenticate, check authorization, [populate params](#params) or perform some side effects. You can add as many hooks
for `.After()` or `.Before()` as you needed. The hook is passed the gin request context.

If you don't want to continue with the request, (eg. user not authorized), you can call `ctx.Abort()`. It's not
recommended, however, to call `ctx.Next()` in a hook as it may mess up the hooks flow.

```go
package demo

import (
	rest "github.com/blackmann/resteasy"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	getHandler := func(id string, p rest.Params) (interface{}, *rest.ServiceError) {
		return map[string]string{}, nil
	}

	checkAuthorization := func(ctx *gin.Context) {
		if ctx.Param("id") == "4" { // nobody is allowed to see the resource with id = 4
			ctx.JSON(http.StatusNotFound, nil)
			ctx.Abort() // do this so the request chain does not proceed
			return
		}
	}

	updateLastViewed := func(ctx *gin.Context) {
		// call some API to update the last viewed for this document
	}

	notifyChannels := func(ctx *gin.Context) {
		// send to the websockets listening to updates on this document
	}

	service, allowedMethods := rest.NewService().Get(getHandler).Service()
	rest.With(service, allowedMethods...).
		Before(checkAuthorization).
		After(updateLastViewed, notifyChannels). // notice more than one hook
		Register(router.Group("/documents"))
}
```

### Params

Params is used to hold data needed to be passed to the service. It's simply a `map[string]interface{}`. Below is an
example on how to interact with the params of a request [context]:

```go
package demo

import (
	rest "github.com/blackmann/resteasy"
	"github.com/gin-gonic/gin"
)

type query struct {
	name string
	age  int
}

func parseQuery(ctx *gin.Context) {
	params := rest.GetParams(ctx)
	params.Set("query", query{name: "Hello", age: 23})
}

func getQuery(ctx *gin.Context) query {
	return rest.GetParams(ctx).Get("query").(query)
}
```

## Constraints

This library does not intend to be an omnipotent library for developing rest services. Therefore, features will be
limited to very few options. Since this library is based on `gin`, 1. you can implement custom behavior
with [hooks](#hooks) or 2. implement your service from scratch. 
