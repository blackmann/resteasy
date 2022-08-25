package rest

import "golang.org/x/exp/slices"

const (
	MethodGet = iota
	MethodFind
	MethodCreate
	MethodPatch
	MethodUpdate
	MethodRemove
)

var AllMethods = []int{MethodGet, MethodFind, MethodCreate, MethodPatch, MethodUpdate, MethodRemove}

func AllMethodsExcept(methods ...int) []int {
	var res []int

	for _, method := range AllMethods {
		if !slices.Contains(methods, method) {
			res = append(res, method)
		}
	}

	return res
}
