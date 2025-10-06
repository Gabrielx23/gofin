package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gofin/internal/container"
	"gofin/pkg/session"
	"gofin/web"
	"gofin/web/components"
	"gofin/cmd/web/handlers"
	"gofin/cmd/web/middleware"
)

func NewRouter(container *container.Container, mux *http.ServeMux) (*chi.Mux, error) {	
	loginComponent, err := components.NewLoginComponent(container)
	if err != nil {
		return nil, fmt.Errorf("failed to create login component: %w", err)
	}
	
	dashboardComponent, err := components.NewDashboardComponent(container)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard component: %w", err)
	}
	
	sessionManager := session.NewSessionManager()
	
	router := chi.NewRouter()
	router.Handle(web.RouteStatic, http.StripPrefix("/static/", http.FileServer(http.Dir(web.StaticDir+"/"))))
	router.Route("/{projectSlug}", func(chiRouter chi.Router) {
		chiRouter.Use(middleware.ProjectBased(container))
		chiRouter.Get("/", handlers.NewMainHandler().Handle)
		chiRouter.Get(web.RouteLogin, handlers.NewLoginHandler(container, loginComponent, sessionManager).Handle)
		chiRouter.Post(web.RouteLogin, handlers.NewLoginHandler(container, loginComponent, sessionManager).Handle)
		chiRouter.Get(web.RouteLogout, handlers.NewLogoutHandler().Handle)
		chiRouter.Get(web.RouteDashboard, middleware.AuthRequired(container, sessionManager)(handlers.NewDashboardHandler(container, dashboardComponent).Handle))
	})
	
	mux.Handle("/", router)
	
	return router, nil
}

func Start(addr string, mux *http.ServeMux) {
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
