package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"

	"github.com/go-eagle/eagle-layout/internal/routers"
)

// NewHTTPServer creates a HTTP server
// grpc -> if open http by protocol, then add second param: userSvc userv1.UserServiceHTTPServer
func NewHTTPServer(c *app.Config) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router

	return srv
}
