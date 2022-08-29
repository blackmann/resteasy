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
