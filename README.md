# MySQL Packets Sniffer

How to use.

    $ vagrant up
    $ GOOS=linux go build parsepacket.go
    $ vagrant ssh
    $ sudo ./bin/tcpdump-mysql | /vagrant/parsepacket

