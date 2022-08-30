package resteasy

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"golang.org/x/exp/maps"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	f, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0755)
	gin.DefaultWriter = io.MultiWriter(f)
	os.Exit(m.Run())
}

func TestService_Builder(t *testing.T) {
	getHandler := func(id string, p Params) (interface{}, *ServiceError) {
		return map[string]string{"id": id, "detail": "Mock detail"}, nil
	}

	findHandler := func(p Params) (interface{}, *ServiceError) {
		return []map[string]string{{"id": "1", "detail": "Mock item 1"}, {"id": "2", "detail": "Mock item 2"}}, nil
	}

	patchHandler := func(id string, data interface{}, p Params) (interface{}, *ServiceError) {
		original := map[string]interface{}{"id": id, "detail": "Mock original"}

		maps.Copy(original, data.(map[string]interface{}))

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
	req, _ = http.NewRequest("GET", "/demo/1", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(t, "{\"detail\":\"Mock detail\",\"id\":\"1\"}", recorder.Body.String())

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/demo", strings.NewReader("{\"detail\":\"new item\"}"))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, "{\"detail\":\"new item\"}", recorder.Body.String())

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/demo/1", strings.NewReader("{\"detail\":\"new replacement\"}"))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, "{\"detail\":\"new replacement\"}", recorder.Body.String())

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/demo/1", strings.NewReader("{\"detail\":\"patching\"}"))

	router.ServeHTTP(recorder, req)

	assert.Equal(t, "{\"detail\":\"patching\",\"id\":\"1\"}", recorder.Body.String())

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/demo/1", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(t, "null", recorder.Body.String())
}

func TestService_HooksFullWalk(t *testing.T) {
	router := gin.Default()
	hookWalk := 0

	setCountParam := func(ctx *gin.Context) {
		GetParams(ctx).Set("count", 3)
		hookWalk += 1
	}

	multiplyCountParam := func(ctx *gin.Context) {
		params := GetParams(ctx)
		params.Set("count", params.Get("count").(int)*5)
		hookWalk += 1
	}

	doExtra := func(ctx *gin.Context) {
		hookWalk += 1
	}

	findHandler := func(p Params) (interface{}, *ServiceError) {
		return map[string]int{"count": p.Get("count").(int)}, nil
	}

	service, allowedMethods := NewService().
		Find(findHandler).
		Service()

	With(service, allowedMethods...).
		Before(setCountParam, multiplyCountParam).
		After(doExtra).
		Register(router.Group("/count"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/count", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Body.String(), "{\"count\":15}")
	assert.Equal(t, hookWalk, 3)
}

func TestService_HooksAbortPremature(t *testing.T) {
	router := gin.Default()
	hookWalk := 0

	setCountParam := func(ctx *gin.Context) {
		GetParams(ctx).Set("count", 3)
		hookWalk += 1
	}

	multiplyCountParam := func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "error"})
	}

	addOneMore := func(ctx *gin.Context) {
		hookWalk += 1
	}

	doExtra := func(ctx *gin.Context) {
		hookWalk += 1
	}

	findHandler := func(p Params) (interface{}, *ServiceError) {
		return map[string]int{"count": p.Get("count").(int)}, nil
	}

	service, allowedMethods := NewService().
		Find(findHandler).
		Service()

	With(service, allowedMethods...).
		Before(setCountParam, multiplyCountParam, addOneMore).
		After(doExtra).
		Register(router.Group("/count"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/count", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Body.String(), "{\"detail\":\"error\"}")
	assert.Equal(t, hookWalk, 1)
}
