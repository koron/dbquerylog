#!/bin/sh

(
cd /home/vagrant

cat << __EOF__ >> .profile
PATH="\$HOME/go/bin:\$PATH"
__EOF__

cat << __EOF__ >> .bashrc
PS1='\[\e]0;\u@\h:\w\a\]\n\[\e[32m\]\u@\h \[\e[33m\]\w\[\e[0m\]\n\$ '
__EOF__

cat << __EOF__ > .bash_alias
alias ls='ls -CF --color=auto --show-control-chars -N'
alias la='ls -A'
alias ll='ls -l'
alias lla='ls -lA'
alias dirs='dirs -v'
alias popd='popd_v'
alias pushd='pushd_v'
pushd_v () {
  "pushd" "$@" > /dev/null && "dirs" -v
}
popd_v () {
  "popd" "$@" > /dev/null && "dirs" -v
}
__EOF__

cat << __EOF__ >> .inputrc
"\C-w": backward-kill-word
"\C-p": history-search-backward
"\C-g": menu-complete
"\C-d": delete-char-or-list
__EOF__

cat << __EOF__ > .screenrc
startup_message off
escape ^Tt
vbell off
hardstatus string "%?%H %?[screen %n%?: %t%?] %h"
__EOF__
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
export PATH="$PATH:/usr/local/go/bin"
mkdir -p go/src/github.com/koron
cd go/src/github.com/koron
ln -s /vagrant dbquerylog
cd dbquerylog
go get -v ./...
)

touch /home/vagrant/.hushlogin
