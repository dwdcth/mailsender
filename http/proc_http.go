package http

import (
	"net/http"

	"github.com/dwdcth/mailsender/proc"
)

func configProcHttpRoutes() {
	// counter
	http.HandleFunc("/counter/all", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, proc.GetAll())
	})

	// trace
	http.HandleFunc("/trace/", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, "not supported right now")
	})
}
