package main

import (
	"time"

	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type Report struct {
	ClientAddr tcpasm.Endpoint
	ServerAddr tcpasm.Endpoint
	Username   string

	StartTime    time.Time
	UpdatedRows  uint64
	ResponseSize uint64
	ElapsedTime  time.Duration
	QueryString  string
	QueryParams  string
}

func (r *Report) Reset() {
	r.StartTime = time.Time{}
	r.UpdatedRows = 0
	r.ResponseSize = 0
	r.ElapsedTime = 0
	r.QueryString = ""
	r.QueryParams = ""
}

func (r *Report) StartQuery(s string) {
	r.StartTime = time.Now()
	r.QueryString = s
}

func (r *Report) Querying() bool {
	return !r.StartTime.IsZero()
}

func (r *Report) FinishQuery() {
	r.ElapsedTime = time.Since(r.StartTime)
}
