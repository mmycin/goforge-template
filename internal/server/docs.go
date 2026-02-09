package server

import (
	"github.com/danielgtaylor/huma/v2"
)

// NewHumaConfig creates a standard Huma configuration with App Key security configured
func NewHumaConfig(title, version, docsPath string) huma.Config {
	config := huma.DefaultConfig(title, version)
	config.DocsPath = docsPath
	config.OpenAPIPath = docsPath + ".json"
	config.SchemasPath = docsPath + "/schemas"

	// Define the App Key security scheme
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"AppKey": {
			Type: "apiKey",
			Name: "X-App-Key",
			In:   "header",
		},
	}

	// Apply it globally to all operations in this API
	config.Security = []map[string][]string{
		{"AppKey": {}},
	}

	return config
}
