package homekit

import (
	"net/http"
)

func (m *HMManager) startUIHandler() {
	m.server.ServeMux().HandleFunc("/ui", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "views/layouts/main.html")
	})
}
