package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/koron/mysql-packet-sniffer/tcpasm"
)

type Report struct {
	ClientAddr tcpasm.Endpoint
	ServerAddr tcpasm.Endpoint
	Username   string

	StartTime    time.Time
	ResponseSize uint64
	ColumnNum    uint64
	UpdatedRows  uint64
	ElapsedTime  time.Duration
	QueryString  string
	QueryParams  string
}

func (r *Report) Reset() {
	r.StartTime = time.Time{}
	r.ResponseSize = 0
	r.ColumnNum = 0
	r.UpdatedRows = 0
	r.ElapsedTime = 0
	r.QueryString = ""
	r.QueryParams = ""
}

func (r *Report) StartQuery(s string, args ...interface{}) {
	r.StartTime = time.Now()
	r.QueryString = s
	if len(args) > 0 {
		b := new(bytes.Buffer)
		for i, arg := range args {
			if i != 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "%#v", arg)
		}
		r.QueryParams = b.String()
	}
}

func (r *Report) Querying() bool {
	return !r.StartTime.IsZero()
}

func (r *Report) FinishQuery() {
	r.ElapsedTime = time.Since(r.StartTime)
}
