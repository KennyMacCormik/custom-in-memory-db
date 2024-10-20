package http

import (
	"context"
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage/wal"
	"custom-in-memory-db/internal/server/network"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type payload struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type errMsg struct {
	Error string `json:"error"`
}

type Server struct {
	timeout time.Duration
	addr    string

	lg     *slog.Logger
	router *gin.Engine
	server *http.Server
}

func (s *Server) New(conf cmd.Config, lg *slog.Logger) {

	s.addr = strings.Join([]string{conf.Network.Host, strconv.Itoa(conf.Network.Port)}, ":")
	s.timeout = conf.Network.Timeout

	s.lg = lg
	s.initGin(conf)
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) Listen(f network.Handler) {
	s.initHandlers(f)
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  s.timeout,
		WriteTimeout: s.timeout,
		IdleTimeout:  s.timeout,
	}
	_ = s.server.ListenAndServe()
}

func (s *Server) initGin(conf cmd.Config) {
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.Use(clientConnLimiter(conf))
	_ = s.router.SetTrustedProxies(nil)
}

// clientConnLimiter limits the number of goroutines actually doing the job.
// Neither gin nor http.Server allows to prevent goroutines from spawning, but we can hold them.
func clientConnLimiter(conf cmd.Config) func(c *gin.Context) {
	limiter := make(chan struct{}, conf.Network.MaxConn)
	return func(c *gin.Context) {
		for {
			select {
			case limiter <- struct{}{}:
				c.Next()
				<-limiter
				return
			default:
				runtime.Gosched()
			}
		}
	}
}

func (s *Server) initHandlers(clientHandler network.Handler) {
	s.cmdHandlers(clientHandler)
}

// connLog inits logger for each request
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

// cmdHandlers inits handlers for the /cmd path
func (s *Server) cmdHandlers(clientHandler network.Handler) {
	isError := func(c *gin.Context, err error) bool {
		if err != nil {
			if err == wal.ErrWalWriteFailed {
				c.JSON(http.StatusInternalServerError, errMsg{err.Error()})
				return true
			}
			c.JSON(http.StatusBadRequest, errMsg{err.Error()})
			return true
		}
		return false
	}

	s.router.GET("/cmd/:key", func(c *gin.Context) {
		key := c.Param("key")
		result, err := clientHandler(strings.NewReader(strings.Join([]string{"GET", key, "\n"}, " ")), s.connLog(c))
		if isError(c, err) {
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
		if isError(c, err) {
			return
		}
		c.Status(http.StatusOK)
	})

	f := func(c *gin.Context) {
		var body payload
		err := c.BindJSON(&body)
		if err == nil {
			_, err = clientHandler(strings.NewReader(strings.Join([]string{"SET", body.Key, body.Value, "\n"}, " ")), s.connLog(c))
			if isError(c, err) {
				return
			}
			c.Status(http.StatusOK)
		}
	}
	s.router.POST("/cmd", f)
	s.router.PUT("/cmd", f)
}
