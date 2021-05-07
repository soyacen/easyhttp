package easyhttp

import "net/http"

// isBodySupported get head options 不支持body
func isBodySupported(m string) bool {
	return !(m == http.MethodHead || m == http.MethodOptions || m == http.MethodGet)
}
