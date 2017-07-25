build: packet-parser mysql-packet-sniffer

run: mysql-packet-sniffer
	sudo ./bin/tcpdump-mysql | ./mysql-packet-sniffer

run-parser: packet-parser
	sudo ./bin/tcpdump-mysql | ./packet-parser

mysql-packet-sniffer:
	go build github.com/koron/mysql-packet-sniffer

packet-parser:
	go build github.com/koron/mysql-packet-sniffer/cmd/packet-parser

.PHONY: mysql-packet-sniffer packet-parser
