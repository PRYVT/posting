package httphandler

import (
	"context"
	"net/http"

	"github.com/PRYVT/posting/pkg/query/httphandler/controller"
	"github.com/PRYVT/utils/pkg/auth"
	ws "github.com/PRYVT/utils/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpHandler struct {
	httpServer     *http.Server
	router         *gin.Engine
	postController *controller.PostController
	authMiddleware *auth.AuthMiddleware
	wsController   *ws.WSController
}

func NewHttpHandler(c *controller.PostController, am *auth.AuthMiddleware, wsController *ws.WSController) *HttpHandler {
	r := gin.Default()
	srv := &http.Server{
		Addr:    "0.0.0.0" + ":" + "5520",
		Handler: r,
	}
	handler := &HttpHandler{
		router:         r,
		httpServer:     srv,
		postController: c,
		authMiddleware: am,
		wsController:   wsController,
	}
	handler.RegisterRoutes()
	return handler
}

func (h *HttpHandler) RegisterRoutes() {
	h.router.Use(auth.CORSMiddleware())
	h.router.GET("/ws", h.wsController.OnRequest)
	h.router.Use(h.authMiddleware.AuthenticateMiddleware)
	{
		h.router.GET("posts/:postId", h.postController.GetPost)
		h.router.GET("posts/", h.postController.GetPosts)
	}
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
