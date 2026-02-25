package http

import "embed"

// SpecFS embeds the OpenAPI specification file.
//
//go:embed openapi.yaml
var SpecFS embed.FS

// SpecPath is the path to the OpenAPI spec within the embedded filesystem.
const SpecPath = "openapi.yaml"
