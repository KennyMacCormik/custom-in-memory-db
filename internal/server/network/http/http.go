package http

import (
	"custom-in-memory-db/internal/server/network"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type payload struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type errMsg struct {
	Error string `json:"error"`
}

type Server struct {
	initDone bool

	lg     *slog.Logger
	router *gin.Engine
}

func (s *Server) Listen(f network.Handler) {
	s.initHandlers(f)
	_ = s.router.Run()
}

func (s *Server) New(lg *slog.Logger) {
	if !s.initDone {
		s.initDone = true
		s.lg = lg
		s.router = gin.Default()
	}
}

func (s *Server) initHandlers(clientHandler network.Handler) {
	s.cmdHandlers(clientHandler)
}

func (s *Server) cmdHandlers(clientHandler network.Handler) {
	s.router.GET("/cmd/:key", func(c *gin.Context) {
		key := c.Param("key")
		result, err := clientHandler(strings.NewReader(strings.Join([]string{"GET", key, "\n"}, " ")), s.lg)
		if err != nil {
			c.JSON(http.StatusBadRequest, errMsg{err.Error()})
			return
		}
		c.JSON(http.StatusOK, payload{
			Key:   key,
			Value: result,
		})
	})
	s.router.DELETE("/cmd/:key", func(c *gin.Context) {
		key := c.Param("key")
		_, err := clientHandler(strings.NewReader(strings.Join([]string{"DEL", key, "\n"}, " ")), s.lg)
		if err != nil {
			c.JSON(http.StatusBadRequest, errMsg{err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	f := func(c *gin.Context) {
		var body payload
		err := c.BindJSON(&body)
		if err == nil {
			_, err = clientHandler(strings.NewReader(strings.Join([]string{"SET", body.Key, body.Value, "\n"}, " ")), s.lg)
			if err != nil {
				c.JSON(http.StatusBadRequest, errMsg{err.Error()})
			}
			c.Status(http.StatusOK)
		}
	}
	s.router.POST("/cmd", f)
	s.router.PUT("/cmd", f)
}
