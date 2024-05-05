package main

import (
	"fmt"
	"net/http"
	"os"
	handler "practice/pkg/api"
	"practice/pkg/middleware"
	"practice/pkg/token"
	"strings"

	"practice/pkg/service"

	"github.com/joho/godotenv"
)

type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (h *maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "practice/import") {
		r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	} else {
		r.Body = http.MaxBytesReader(w, r.Body, h.n)
	}

	h.h.ServeHTTP(w, r)
}

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	asvc := service.NewAccountService()
	h := handler.NewHTTPHandler(asvc)
	h = middleware.WithAuth(token.NewJwTVerifier(asvc))(h)

	fmt.Printf("Sub is running on port 8080")
	if err := http.ListenAndServe(":8080", &maxBytesHandler{h: h, n: 4 << 20}); err != nil {
		fmt.Println("Starting http server error: ", err)
		os.Exit(1)
	}
}
