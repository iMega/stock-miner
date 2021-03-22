package helpers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/gomega"
)

var (
	MockHTTPServer = New()

	currentHandlerFunc http.HandlerFunc
)

type mockHTTPServer struct {
	server   *http.Server
	handlers map[string]http.HandlerFunc
}

func New() *mockHTTPServer {
	router := httprouter.New()
	m := &mockHTTPServer{
		server: &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: router,
		},
		handlers: make(map[string]http.HandlerFunc),
	}
	router.NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			currentHandlerFunc(w, r)
		},
	)

	return m
}

func (m *mockHTTPServer) BeforeSuite() {
	go func() {
		if err := m.server.ListenAndServe(); err != nil {
			switch err {
			case http.ErrServerClosed:
				return
			default:
				fmt.Printf("failed to serve http: %s", err)
			}
		}
	}()
}

func (m *mockHTTPServer) AfterSuite() {
	err := m.server.Shutdown(context.Background())
	Expect(err).NotTo(HaveOccurred())
}

func (m *mockHTTPServer) AddHandler(hf http.HandlerFunc) {
	currentHandlerFunc = hf
}
