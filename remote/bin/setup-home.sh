#!/bin/sh

touch /home/vagrant/.hushlogin

(
cd /vagrant/remote
for entry in * ; do
  if [ ! -e "/home/vagrant/$entry" ] ; then
    ln -v -s "/vagrant/remote/$entry" "/home/vagrant/$entry"
  fi
done
)
