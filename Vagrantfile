# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "minimal/xenial64"

  # for MySQL
  config.vm.network :forwarded_port, guest:3306, host:3306

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
    vb.cpus = 2
  end

  config.ssh.insert_key = false

  config.vm.provision "shell", inline: <<-SHELL
    update-locale LANG=C.UTF-8 LANGUAGE=
    sed -i.bak -e 's!http://\\(archive\\|security\\).ubuntu.com/!ftp://ftp.jaist.ac.jp/!g' /etc/apt/sources.list
    apt update

    # MySQL
    apt install -y debconf-utils
    debconf-set-selections <<< 'mysql-server mysql-server/root_password password mysql123'
    debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password mysql123'
    apt install -y mysql-server
    echo "default-character-set=utf8" >> /etc/mysql/conf.d/mysql.cnf
    echo "default-character-set=utf8" >> /etc/mysql/conf.d/mysqldump.cnf
    cat >> /etc/mysql/mysql.conf.d/mysqld.cnf <<__EOS__
bind-address=0.0.0.0
character-set-server=utf8
skip-character-set-client-handshake
default-storage-engine=INNODB
__EOS__
    systemctl restart mysql
    mysql -u root --password=mysql123 <<__EOS__
CREATE DATABASE vagrant;
GRANT ALL ON vagrant.* TO vagrant@"%" IDENTIFIED BY 'db1234';
FLUSH PRIVILEGES;
__EOS__

    # Golang 1.8
    apt install -y g++ gcc libc6-dev pkg-config
    curl -s -o go.tar.gz https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
    tar xzf go.tar.gz -C /usr/local
    rm -r go.tar.gz
    echo 'export PATH="$PATH:/usr/local/go/bin"' > /etc/profile.d/go.sh
    hash -r

    # Others tools
    apt install -y man tcpdump screen

    sudo -i -u vagrant /vagrant/remote/bin/setup-home.sh
  SHELL

end
