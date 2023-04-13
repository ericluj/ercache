package ercache

import (
	"encoding/json"
	"io"

	log "github.com/ericluj/elog"
	"github.com/hashicorp/raft"
)

type FSM struct {
	server *Server
}

func newFSM(server *Server) *FSM {
	return &FSM{
		server: server,
	}
}

type logEntryData struct {
	Key   string
	Value string
}

func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	e := logEntryData{}
	if err := json.Unmarshal(logEntry.Data, &e); err != nil {
		panic("Failed unmarshaling Raft log entry. This is a bug.")
	}
	ret := f.server.Cache.Set(e.Key, e.Value)
	log.Infof("fms.Apply(), logEntry:%s, ret:%v", logEntry.Data, ret)
	return ret
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{cache: f.server.Cache}, nil
}

func (f *FSM) Restore(serialized io.ReadCloser) error {
	return f.server.Cache.Unmarshal(serialized)
}
