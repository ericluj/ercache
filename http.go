package ercache

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/ericluj/elog"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
)

type httpServer struct {
	router http.Handler
	server *Server
}

func newHTTPServer(server *Server) *httpServer {
	s := &httpServer{
		server: server,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("ping", s.ping)
	router.GET("get", s.get)
	router.GET("set", s.set)
	router.GET("join", s.join)
	s.router = router

	return s
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *httpServer) ping(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (s *httpServer) get(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		log.Infof("get error: nil key")
		c.String(http.StatusInternalServerError, "internal error")
		return
	}
	ret := s.server.Cache.Get(key)
	c.String(http.StatusOK, ret)
}

func (s *httpServer) set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == "" {
		log.Infof("set error: nil key or nil value")
		c.String(http.StatusInternalServerError, "internal error")
		return
	}

	event := logEntryData{Key: key, Value: value}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Infof("json.Marshal error:%v", err)
		c.String(http.StatusInternalServerError, "internal error")
		return
	}

	applyFuture := s.server.Raft.raft.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		log.Infof("raft.Apply error:%v", err)
		c.String(http.StatusInternalServerError, "internal error")
		return
	}
	c.String(http.StatusOK, "ok")
}

func (s *httpServer) join(c *gin.Context) {
	peerAddress := c.Query("peerAddress")
	if peerAddress == "" {
		log.Infof("invalid peerAddress")
		c.String(http.StatusInternalServerError, "invalid peerAddress")
		return
	}

	addPeerFuture := s.server.Raft.raft.AddVoter(raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		log.Infof("raft.AddVoter error:%v", err)
		c.String(http.StatusInternalServerError, "internal error")
		return
	}
	c.String(http.StatusOK, "ok")
}
