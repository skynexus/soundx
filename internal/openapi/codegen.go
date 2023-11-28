//go:build codegen

package openapi

// Note: you might see a warning for this file - this is because the
// the Go language server does not handle build tags well:
// https://github.com/golang/go/issues/29202

// Force import of oapi-codegen package to ensure its dependencies are
// available in go.sum when oapi-codegen is executed using 'go run'.
//
// We use the oapi-codegen tool via 'make api' to generate code
// from the available openapi specification.
//
// Note that we also include the 'codegen' build tag (also known as a
// build constraint) to ensure this dependency is not included when
// running a standard build of the soundx system.

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
)
