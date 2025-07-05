package homekit

import (
	"net/http"
)

func (m *HMManager) startUIHandler() {
	m.server.ServeMux().HandleFunc("/manifest.json", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		http.ServeFile(res, req, "web/manifest.json")
	})
	m.server.ServeMux().HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(res, req, "web/index.html")
	})
}
