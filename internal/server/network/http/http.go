package http

import (
	"custom-in-memory-db/internal/server/network"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		s.initGin()
	}
}

func (s *Server) initGin() {
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	_ = s.router.SetTrustedProxies(nil)
}

func (s *Server) initHandlers(clientHandler network.Handler) {
	s.cmdHandlers(clientHandler)
}

func (s *Server) connLog(c *gin.Context) *slog.Logger {
	uid := uuid.New()
	lg := s.lg.With("ID", uid,
		"ClientIP", c.ClientIP(),
		"Method", c.Request.Method,
		"Path", c.Request.URL.Path,
		"Proto", c.Request.Proto,
		"Headers", c.Request.Header,
	)
	lg.Info("connection accepted")
	return lg
}

func (s *Server) cmdHandlers(clientHandler network.Handler) {
	s.router.GET("/cmd/:key", func(c *gin.Context) {
		key := c.Param("key")
		result, err := clientHandler(strings.NewReader(strings.Join([]string{"GET", key, "\n"}, " ")), s.connLog(c))
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
		_, err := clientHandler(strings.NewReader(strings.Join([]string{"DEL", key, "\n"}, " ")), s.connLog(c))
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
			_, err = clientHandler(strings.NewReader(strings.Join([]string{"SET", body.Key, body.Value, "\n"}, " ")), s.connLog(c))
			if err != nil {
				c.JSON(http.StatusBadRequest, errMsg{err.Error()})
			}
			c.Status(http.StatusOK)
		}
	}
	s.router.POST("/cmd", f)
	s.router.PUT("/cmd", f)
}
