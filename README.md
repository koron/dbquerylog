# MySQL Packets Sniffer

How to use.

    $ vagrant up
    $ GOOS=linux go build parsepacket.go
    $ vagrant ssh
    $ sudo ./bin/tcpdump-mysql | /vagrant/parsepacket

## Example

```console
$ make run
go build github.com/koron/mysql-packet-sniffer
sudo ./bin/tcpdump-mysql | ./mysql-packet-sniffer
tcpdump: listening on lo, link-type EN10MB (Ethernet), capture size 262144 bytes
127.0.0.1:60720 query "show databases" finished in 14.918µs
127.0.0.1:60720 query "show tables" finished in 13.349µs
127.0.0.1:60720 query "select @@version_comment limit 1" finished in 18.019µs
2017/07/26 02:51:45 stream closed for 127.0.0.1:60720
2017/07/26 02:51:45 stream closed for 127.0.0.1:60720
```

## Report format

Report format is based on TSV (tab separated values).
Each rows represent database queries.
Each columns are consisted by below:

*   start time
*   client address (IP address and port)
*   server address (IP address and port)
*   username of database
*   total response size in byte
*   column number in response
*   row numbers which updated
*   elapsed time
*   query
*   parameters (available for prepared statement only)
