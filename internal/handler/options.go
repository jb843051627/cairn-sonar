package handler

import "net/http"

func allowMethods(methods ...string) map[string]bool {
	out := make(map[string]bool, len(methods))
	for _, method := range methods {
		out[method] = true
	}
	return out
}
func methodAllowed(r *http.Request, methods map[string]bool) bool { return methods[r.Method] }
