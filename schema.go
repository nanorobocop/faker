package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

func loadSchema(path string) (schema *openapi3.T, err error) {
	loader := openapi3.NewLoader()
	schemaURI, err := url.Parse(path)
	if err != nil {
		schema, err = loader.LoadFromFile(path)
	} else {
		schema, err = loader.LoadFromURI(schemaURI)
	}

	return schema, err
}

func (a *app) getSchemaRouter() http.HandlerFunc {
	schema, err := loadSchema(*a.schemaPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load schema: %+v", err))
	}

	if err := schema.Validate(context.Background()); err != nil {
		panic(fmt.Sprintf("Schema validation failed: %+v", err))
	}

	router, err := gorillamux.NewRouter(schema)
	if err != nil {
		panic(fmt.Sprintf("Failed to create router: %+v", err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		route, params, err := router.FindRoute(r)
		if err != nil {
			resp := fmt.Sprintf("Failed to find route: %+v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(resp))
			return
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: params,
			Route:      route,
		}
		if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
			resp := fmt.Sprintf("Request validation failed: %+v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(resp))
			return
		}

		resp := fmt.Sprintf(`Route: %s %s
Params: %+v
`, route.Method, route.Path, params)
		w.Write([]byte(resp))
	}
}
