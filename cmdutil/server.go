package cmdutil

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type TempServer struct {
	Port        string
	Server      *http.Server
	RouteName   string
	HandlerFunc func(http.ResponseWriter, *http.Request)
}

func CreateTempServer(t TempServer) *TempServer {
	router := mux.NewRouter().StrictSlash(true)
	handler := cors.Default().Handler(router)
	router.HandleFunc(t.RouteName, t.HandlerFunc)
	t.Server = &http.Server{Addr: t.Port, Handler: handler}
	return &t
}

func (t *TempServer) CloseServer() {
	t.Server.Shutdown(context.Background())
}
