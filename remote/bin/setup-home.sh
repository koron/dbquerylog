#!/bin/sh

(
cd /home/vagrant

cat << __EOF__ > .screenrc
startup_message off
escape ^Tt
vbell off
hardstatus string "%?%H %?[screen %n%?: %t%?] %h"
__EOF__

touch .hushlogin
)

(
cd /vagrant/remote
for entry in * ; do
  if [ ! -e "/home/vagrant/$entry" ] ; then
    ln -v -s "/vagrant/remote/$entry" "/home/vagrant/$entry"
  fi
done
)

(
cd /home/vagrant
mkdir -p go/src/github.com/koron
cd go/src/github.com/koron
ln -s /vagrant mysql-packet-sniffer
cd mysql-packet-sniffer
go get -v ./...
)
