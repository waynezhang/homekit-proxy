package homekit

import (
	"net/http"
	"strings"
)

func (m *HMManager) startUIHandler() {
	m.server.ServeMux().HandleFunc("/views/components/{name}", func(res http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/views/components/")
		name := strings.SplitN(path, "/", 2)[0]
		http.ServeFile(res, req, "views/components/"+name)
	})
	m.server.ServeMux().HandleFunc("/manifest.json", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "views/layouts/manifest.json")
	})
	m.server.ServeMux().HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "views/layouts/main.html")
	})
}
