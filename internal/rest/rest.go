package rest

import (
	"context"
	"fmt"
	"github.com/def-stocks-t/gateway/internal/config"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type ContextKey string

type Rest struct {
	config config.Config

	httpServer *http.Server
	lock       sync.Mutex

	logger *log.Logger
}

func NewRestService(c config.Config, l *log.Logger) *Rest {
	return &Rest{
		config: c,
		logger: l,
	}
}

func (s *Rest) ServiceProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}

func (s *Rest) Run() {
	s.logger.Infof("Starting HTTP rest server on port %s", s.config.Port)

	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(s.config.Port, s.routes())
	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()
	s.logger.Infof("HTTP rest server terminated, %s", err)
}

func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))

	router.Route("/api/v1", func(rfront chi.Router) {
		remote, err := url.Parse(fmt.Sprintf("http://%s:%s/", s.config.Services.Storage.Host, s.config.Services.Storage.Port))
		if err != nil {
			log.WithError(err).Error("failed to parse storage url")
			return
		}
		storageProxy := httputil.NewSingleHostReverseProxy(remote)

		rfront.Group(func(routerPrivate chi.Router) {
			routerPrivate.Group(func(routerProxy chi.Router) {
				routerProxy.Get("/storage/*", s.ServiceProxyHandler(storageProxy))
				routerProxy.Post("/storage/*", s.ServiceProxyHandler(storageProxy))
			})
		})
	})

	return router
}

// Shutdown rest http server
func (s *Rest) Shutdown() {
	s.logger.Infof("Shutdown HTTP rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Warnf("HTTP rest shutdown error")
		}
		s.logger.Infof("Shutdown HTTP rest server completed")
	}
	s.lock.Unlock()
}

func (s *Rest) makeHTTPServer(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
