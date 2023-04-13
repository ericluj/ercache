package ercache

import "flag"

type options struct {
	dataDir        string // 数据存储路径
	httpAddress    string // http地址
	raftTCPAddress string // raft地址
	boostrap       bool   // 是否master
	joinAddress    string // 加入的peer地址
}

func newOptions() *options {
	opts := &options{}

	httpAddress := flag.String("http", "127.0.0.1:6000", "http address")
	raftTCPAddress := flag.String("raft", "127.0.0.1:7000", "raft tcp address")
	node := flag.String("node", "node1", "raft node name")
	bootstrap := flag.Bool("bootstrap", false, "start as raft cluster")
	joinAddress := flag.String("join", "", "join address for raft cluster")
	flag.Parse()

	opts.dataDir = "./" + *node
	opts.httpAddress = *httpAddress
	opts.raftTCPAddress = *raftTCPAddress
	opts.boostrap = *bootstrap
	opts.joinAddress = *joinAddress
	return opts
}
