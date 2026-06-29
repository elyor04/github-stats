package server

import (
	"net/http"

	"github-stats/internal/stats"
	"github-stats/internal/svg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	statsSvc *stats.Service
	log      *zap.Logger
	engine   *gin.Engine
}

func New(statsSvc *stats.Service, log *zap.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())

	s := &Server{statsSvc: statsSvc, log: log, engine: engine}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func (s *Server) registerRoutes() {
	s.engine.GET("/stats", s.handleStatsSVG)
	s.engine.GET("/languages", s.handleLanguagesSVG)
	s.engine.GET("/api/stats", s.handleStatsJSON)
	s.engine.GET("/api/languages", s.handleLanguagesJSON)
	s.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})
}

func queryParams(c *gin.Context) (username string, includePrivate bool) {
	username = c.DefaultQuery("username", "elyor04")
	includePrivate = c.Query("private") == "true"
	return
}

func (s *Server) handleStatsSVG(c *gin.Context) {
	username, includePrivate := queryParams(c)
	result, err := s.statsSvc.GetUserStats(username, includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, "image/svg+xml", []byte(svg.GenerateStatsSVG(result)))
}

func (s *Server) handleLanguagesSVG(c *gin.Context) {
	username, includePrivate := queryParams(c)
	result, err := s.statsSvc.GetLanguageStats(username, includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=3600")
	c.Data(http.StatusOK, "image/svg+xml", []byte(svg.GenerateLanguagesSVG(result)))
}

func (s *Server) handleStatsJSON(c *gin.Context) {
	username, includePrivate := queryParams(c)
	result, err := s.statsSvc.GetUserStats(username, includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleLanguagesJSON(c *gin.Context) {
	username, includePrivate := queryParams(c)
	result, err := s.statsSvc.GetLanguageStats(username, includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, result)
}
