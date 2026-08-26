package handler

import "net/http"

func MethodValue(r *http.Request) string { return r.Header.Get("X-Cairn-Method") }
func MethodReady(r *http.Request) bool   { return r.Method != "" }
