build: packet-parser dbquerylog

run: dbquerylog
	sudo ./bin/tcpdump-mysql | ./dbquerylog

run-debug: dbquerylog
	sudo ./bin/tcpdump-mysql | ./dbquerylog -debug

run-parser: packet-parser
	sudo ./bin/tcpdump-mysql | ./packet-parser

dbquerylog:
	go build github.com/koron/dbquerylog

packet-parser:
	go build github.com/koron/dbquerylog/cmd/packet-parser

.PHONY: dbquerylog packet-parser
