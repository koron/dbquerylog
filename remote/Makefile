build: mysql-packet-sniffer

run: build
	sudo ./bin/tcpdump-mysql | ./mysql-packet-sniffer

mysql-packet-sniffer:
	go build github.com/koron/mysql-packet-sniffer

.PHONY: mysql-packet-sniffer
