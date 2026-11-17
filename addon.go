package main

import (
	"net/http"

	"github.com/antresol-app/antresol/addonkit"
)

func newHandler(env addonkit.Env, core *addonkit.Client, ready *addonkit.Ready) http.Handler {
	mux := http.NewServeMux()
	// the core's health check: {"ok": true, "ready": …}
	mux.HandleFunc("GET /health", ready.Health)
	// GET /api/addons/<id>/x/hello in the app; X-Antresol-User says who asks
	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		u, ok := addonkit.UserFrom(r)
		if !ok {
			addonkit.WriteError(w, http.StatusUnauthorized, "no user")
			return
		}
		addonkit.WriteJSON(w, http.StatusOK, map[string]string{"name": u.Name})
	})
	// /api/admin/… (x/admin/… in the app) is for admins only: the core
	// refuses it to anyone else, RequireAdmin checks again
	mux.Handle("GET /api/admin/hello", addonkit.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := addonkit.UserFrom(r)
		addonkit.WriteJSON(w, http.StatusOK, map[string]string{"admin": u.Name})
	})))
	// the hooks the core calls (POST /hooks/<name>): handle exactly those
	// the manifest declares in "hooks"; this addon declares none
	hooks := &addonkit.HookMux{}
	mux.Handle("POST /hooks/", hooks)
	return mux
}
