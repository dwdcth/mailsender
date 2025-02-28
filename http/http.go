package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dwdcth/mailsender/g"
)

// AuthMiddleware 验证 Bearer Token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := g.GetConfig().Http.Token
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 检查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// 验证 token
		if parts[1] != token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func Start() {
	go startHttpServer()
}

// add url mapping config
func addRoutes() {
	configCommonRoutes()
	configMailSenderApiRoutes()
	configProcHttpRoutes()
}

func startHttpServer() {
	if !g.GetConfig().Http.Enable {
		return
	}

	addr := g.GetConfig().Http.Listen
	if addr == "" {
		return
	}

	addRoutes()
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.startHttpServer ok, listening", addr)
	log.Fatalln(s.ListenAndServe())
}

// interfaces
type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, Dto{Msg: "success", Data: data})
}

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}
