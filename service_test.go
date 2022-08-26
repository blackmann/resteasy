package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"golang.org/x/exp/maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestService_Builder(t *testing.T) {
	getHandler := func(id string, p Params) (interface{}, *ServiceError) {
		return map[string]string{"id": id, "detail": "Mock detail"}, nil
	}

	findHandler := func(p Params) (interface{}, *ServiceError) {
		return []map[string]string{{"id": "1", "detail": "Mock item 1"}, {"id": "2", "detail": "Mock item 2"}}, nil
	}

	patchHandler := func(id string, data interface{}, p Params) (interface{}, *ServiceError) {
		original := map[string]string{"id": id, "detail": "Mock original"}

		maps.Copy(original, data.(map[string]string))

		return original, nil
	}

	createHandler := func(data interface{}, p Params) (interface{}, *ServiceError) {
		return data, nil
	}

	updateHandler := func(id string, data interface{}, p Params) (interface{}, *ServiceError) {
		return data, nil
	}

	removeHandler := func(id string, p Params) (interface{}, *ServiceError) {
		return nil, nil
	}

	service, allowedMethods := NewService().
		Create(createHandler).
		Get(getHandler).
		Find(findHandler).
		Patch(patchHandler).
		Update(updateHandler).
		Remove(removeHandler).
		Service()

	router := gin.Default()
	With(service, allowedMethods...).Register(router.Group("/demo"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/demo", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(
		t,
		"[{\"detail\":\"Mock item 1\",\"id\":\"1\"},{\"detail\":\"Mock item 2\",\"id\":\"2\"}]",
		recorder.Body.String(),
	)

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/demo", strings.NewReader("{\"detail\":\"new item\"}"))

	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	assert.Equal(t, "{\"detail\":\"new item\"}", recorder.Body.String())
}
