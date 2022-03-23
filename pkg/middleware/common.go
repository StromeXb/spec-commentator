package middleware

import (
	"reflect"

	"github.com/go-kit/kit/endpoint"
)

func SetMiddleware(mw func(field reflect.StructField) endpoint.Middleware, s reflect.Value) {
	if s.Kind() == reflect.Struct {
		for i := 0; i < s.NumField(); i++ {
			valueField := s.Field(i)
			typeField := s.Type().Field(i)
			if valueField.Kind() == reflect.Struct {
				SetMiddleware(mw, valueField)
				continue
			}
			e, ok := valueField.Interface().(endpoint.Endpoint)
			if !ok {
				continue
			}
			e = mw(typeField)(e)
			s.Field(i).Set(reflect.ValueOf(&e).Elem())
		}
	}
}
