package handler

import "net/http"

func HeaderValue(r *http.Request) string { return r.Header.Get("X-Cairn-Header") }
func HeaderReady(r *http.Request) bool   { return r.Method != "" }
