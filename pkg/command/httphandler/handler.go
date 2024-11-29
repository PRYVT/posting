package httphandler

import (
	"context"
	"net/http"

	"github.com/PRYVT/posting/pkg/command/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	postController *controller.PostController
	authMiddleware *auth.AuthMiddleware
}

func NewHttpHandler(c *controller.PostController, m *auth.AuthMiddleware) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5519",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		postController: c,
		authMiddleware: m,
	}

	handler.RegisterRoutes()

	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.Use(h.authMiddleware.AuthenticateMiddleware)
	h.router.POST("posts/", h.postController.CreatePost)
}

func (h *HttpHandler) Start() error {
	return h.httpServer.ListenAndServe()
}

func (h *HttpHandler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Error during reading response body")
	}
}
