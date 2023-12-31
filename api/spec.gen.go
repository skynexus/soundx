// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
