package rest

import "net/http"

type ServiceError struct {
	Code   int         `json:"code"`
	Detail interface{} `json:"detail"`
}

func MethodNotAllowed() *ServiceError {
	return &ServiceError{
		Code:   http.StatusMethodNotAllowed,
		Detail: map[string]string{"detail": "Method not allowed"},
	}
}
