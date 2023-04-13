package ercache

import (
	"net"
	"net/http"

	log "github.com/ericluj/elog"
)

type Server struct {
	Opts  *options
	Cache *cache
	Raft  *raftNode
}

func NewServer() *Server {
	s := &Server{
		Opts:  newOptions(),
		Cache: newCache(),
	}

	// http
	l, err := net.Listen("tcp", s.Opts.httpAddress)
	if err != nil {
		log.Infof("http server error:%v", err)
	}
	log.Infof("http server listen:%s", l.Addr())
	httpServer := newHTTPServer(s)
	go func() {
		http.Serve(l, httpServer)
	}()

	// raft
	raft, err := newRaftNode(s.Opts, s)
	if err != nil {
		log.Infof("newRaftNode error:%v", err)
	}
	s.Raft = raft

	if s.Opts.joinAddress != "" {
		err = joinRaftCluster(s.Opts)
		if err != nil {
			log.Infof("join raft cluster failed:%v", err)
		}
	}

	// monitor leadership
	go func() {
		for {
			select {
			case leader := <-s.Raft.leaderNotifyCh:
				if leader {
					log.Infof("become leader, enable write api")
				} else {
					log.Infof("become follower, close write api")
				}
			}
		}
	}()

	return s
}
