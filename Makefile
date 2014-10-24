all: morganfield

morganfield:
	go get github.com/blackjack/syslog
	go build -o morganfield github.com/mikkolehtisalo/morganfield/main

install:
	mkdir -p /opt/morganfield/bin
	mkdir -p /opt/morganfield/etc
	cp morganfield /opt/morganfield/bin/
	chmod 744 /opt/morganfield/bin/morganfield
	cp morganfield.service /etc/systemd/system/
	chmod 755 /etc/systemd/system/morganfield.service

selinux:
	selinux/morganfield.sh

clean:
	rm morganfield

test:
	go test .

.PHONY: selinux
