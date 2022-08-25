package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestService_Find(t *testing.T) {
	router := gin.Default()
	group := router.Group("/demo")

	findHandler := func(p Param) (interface{}, *ServiceError) {
		return []int{1, 2, 3}, nil
	}

	service, allowedMethods := NewService().
		Find(findHandler).
		Service()

	With(service, allowedMethods...).Register(group)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/demo", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Body.String(), "[1,2,3]")
}

func TestService_Raw(t *testing.T) {
	router := gin.Default()

	router.GET("/demo", func(context *gin.Context) {
		context.JSON(http.StatusOK, []int{1, 2, 3})
	})

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/demo", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Body.String(), "[1,2,3]")
}

func BenchmarkService_Find(b *testing.B) {
	f, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0755)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	group := router.Group("/demo")

	findHandler := func(p Param) (interface{}, *ServiceError) {
		return []int{1, 2, 3}, nil
	}

	service, allowedMethods := NewService().
		Find(findHandler).
		Service()

	With(service, allowedMethods...).Register(group)

	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/demo", nil)
		router.ServeHTTP(recorder, req)
	}
}

func BenchmarkService_Raw(b *testing.B) {
	f, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0755)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	router.GET("/demo", func(context *gin.Context) {
		context.JSON(http.StatusOK, []int{1, 2, 3})
	})

	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/demo", nil)
		router.ServeHTTP(recorder, req)
	}
}
