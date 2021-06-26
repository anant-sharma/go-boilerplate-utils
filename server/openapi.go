package server

import (
	"net/http"
)

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
