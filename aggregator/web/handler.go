package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	JsonData []byte
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	// Добавьте обработчик CORS
	router.Use(cors.Default())
	router.GET("", h.Print)
	router.GET("/api", h.To)
	return router
}

func (h *Handler) Print(c *gin.Context) {
	c.String(http.StatusOK, "Make a GET request on the path /api")
}

func (h *Handler) To(c *gin.Context) {
	// Отправляем JsonData в ответ на GET /api
	logrus.Infof("Send JsonData /api")
	c.Data(http.StatusOK, "application/json", h.JsonData)
}
