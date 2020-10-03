PREFIX=$(PWD)
BINDIR=$(PREFIX)/bin

COMPILER=go

GOFILES=$(wildcard ./*.go ./**/*.go)

#-------------------------------------------------------------------------------

APP_SRV=isuumo.go.service

GIT_ROOT=.
GIT_EMAIL=syu.takayama@gmail.com
GIT_NAME=Shu Takayama

NGINX_LOG=/tmp/access.log

KATARIBE_CFG=/etc/kataribe.toml

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=isucon
DB_PASS=isucon
DB_NAME=isuumo

SLOW_LOG=/tmp/slow-query.log

SLACKCAT_CNL=isucon

LOGS_DIR=/etc/logs

#-------------------------------------------------------------------------------

PPROF=go tool pprof

KATARIBE=kataribe -f $(KATARIBE_CFG)

MYSQL=sudo mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -p$(DB_PASS) $(DB_NAME)

SLACKCAT=slackcat --channel $(SLACKCAT_CNL)

WHEN:=$(shell date +%H:%M:%S)

#-------------------------------------------------------------------------------

.PHONY: build
build: $(GOFILES) ## Build executable files
	@GOBIN=$(BINDIR) $(COMPILER) install ./cmd/...  

.PHONY: clean
clean: ## Cleanup files
	@$(RM) -r $(BINDIR)

#-------------------------------------------------------------------------------

.PHONY: prebench-dev
prebench-dev: pull clean build restart slow-on 

.PHONY: prebench-prod
prebench-prod: pull clean build restart slow-off

.PHONY: pull
pull:
	cd $(GIT_ROOT) && \
		git pull

.PHONY: slow-on
slow-on:
	$(MYSQL) -e "set global slow_query_log_file = '$(SLOW_LOG)'; set global long_query_time = 0; set global slow_query_log = ON;"

.PHONY: slow-off
slow-off:
	$(MYSQL) -e "set global slow_query_log = OFF;"

.PHONY: restart
restart: restart-nginx restart-mysql restart-app 

.PHONY: restart-app
restart-app:
	sudo systemctl restart $(APP_SRV)

.PHONY: restart-nginx
restart-nginx:
	@if [ -f $(NGINX_LOG) ]; then \
		sudo mkdir -p $(LOGS_DIR)/$(WHEN) ; \
		sudo mv -f $(NGINX_LOG) $(LOGS_DIR)/$(WHEN)/ ; \
	fi
	sudo systemctl restart nginx

.PHONY: restart-mysql
restart-mysql:
	@if [ -f $(SLOW_LOG) ]; then \
		sudo mkdir -p $(LOGS_DIR)/$(WHEN) ; \
		sudo mv -f $(SLOW_LOG) $(LOGS_DIR)/$(WHEN)/ ; \
	fi
	sudo systemctl restart mysql

#-------------------------------------------------------------------------------

.PHONY: analyze
analyze: kataribe slow

.PHONY: pprof
pprof:
	$(PPROF) -png -output pprof1.png http://localhost:6060/debug/pprof/profile?seconds=60 &
	$(PPROF) -png -output pprof2.png http://localhost:6061/debug/pprof/profile?seconds=60 &
	$(SLACKCAT) -n pprof1.png pprof1.png
	$(SLACKCAT) -n pprof2.png pprof2.png

.PHONY: kataribe
kataribe:
	sudo cat $(NGINX_LOG) | $(KATARIBE) | $(SLACKCAT) --tee

.PHONY: slow
slow: 
	sudo cat $(SLOW_LOG) | pt-query-digest | $(SLACKCAT) --tee

#-------------------------------------------------------------------------------

.PHONY: install-tools
install-tools: install-git install-unzip install-kataribe install-myprofiler install-pt install-dstat install-slackcat

.PHONY: install-git
install-git:
	sudo apt install -y git # TODO
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_NAME)"

.PHONY: install-unzip
install-unzip:
	sudo apt install -y unzip # TODO

.PHONY: install-kataribe
install-kataribe:
	wget https://github.com/matsuu/kataribe/releases/download/v0.4.1/kataribe-v0.4.1_linux_amd64.zip -O kataribe.zip
	mkdir -p tmp_kataribe
	unzip -o kataribe.zip -d tmp_kataribe
	rm -f kataribe.zip
	sudo cp tmp_kataribe/kataribe /usr/local/bin/
	rm -rf tmp_kataribe
	sudo chmod +x /usr/local/bin/kataribe
	kataribe -generate
	sudo cp kataribe.toml $(KATARIBE_CFG)
	rm -f kataribe.toml

.PHONY: install-myprofiler
install-myprofiler:
	wget https://github.com/KLab/myprofiler/releases/latest/download/myprofiler.linux_amd64.tar.gz -O myprofiler.tar.gz
	tar xf myprofiler.tar.gz
	rm -f myprofiler.tar.gz
	sudo cp myprofiler /usr/local/bin/
	rm -f myprofiler
	sudo chmod +x /usr/local/bin/myprofiler

.PHONY: install-pt
install-pt:
	sudo apt install -y percona-toolkit # TODO

.PHONY: install-dstat
install-dstat:
	sudo apt install -y dstat # TODO

.PHONY: install-slackcat
install-slackcat:
	wget https://github.com/bcicen/slackcat/releases/download/v1.6/slackcat-1.6-linux-amd64 -O slackcat
	sudo cp slackcat /usr/local/bin/
	rm -f slackcat
	sudo chmod +x /usr/local/bin/slackcat
	slackcat --configure