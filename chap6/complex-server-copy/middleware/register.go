package middleware

import (
	"github.com/alphajin1/PracticalGo/complex-server-copy/config"
	"net/http"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, c), c)
}
