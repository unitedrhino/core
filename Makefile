# -*- coding:utf-8 -*-
.PHOmakeNY: build

build:build.clean mod cp.etc build.api  build.view  build.sys  build.timedjob build.timedscheduler

buildback: build.clean mod cp.etc build.api  build.view


buildone: buildback moduleupdate build.front


runall:  run.timedjob run.timedscheduler run.sys  run.api run.view

packone:  buildone  toremote

packback:  buildback  toremote


toremote:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>tormote cmd<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rsync -r -v ./cmd/* root@120.79.205.165:/root/git/iThings/core


moduleupdate:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@git submodule update --init --recursive
	@git submodule foreach git checkout master
	@git submodule foreach git pull

build.front:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>$@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@mkdir -p ./cmd/dist/app/core
	@cd module/front/core  && npm install && npm run pro && cp -rf ./dist/* ../../../cmd/dist/app/core
#	@git switch 'origin/dev/dyb/v1.0.4'
	@mkdir -p ./cmd/dist/app/system-manage
	@cd module/front/systemManage  && npm install && npm run pro && cp -rf ./dist/* ../../../cmd/dist/app/system-manage
#	@git switch 'origin/dyb/v1.0.5'

killall:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>killing all<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@killall  apisvr &
	@killall  syssvr &
	@killall  timedjob &
	@killall  timedscheduler &

build.clean:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>clean cmd<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rm -rf ./cmd/*

mod:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>downloading $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go mod tidy


cp.etc:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>copying etc<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@mkdir -p ./cmd/etc/
	@cp -rf ./service/apisvr/etc/* ./cmd/etc/
	@cp -rf ./service/viewsvr/etc/* ./cmd/etc/


build.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -o ./cmd/apisvr ./service/apisvr

build.view:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build -o ./cmd/viewsvr ./service/viewsvr

build.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/syssvr ./service/syssvr

build.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/timedjobsvr ./service/timed/timedjobsvr

build.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go build  -o ./cmd/timedschedulersvr ./service/timed/timedschedulersvr


run.api:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./apisvr &  cd ..

run.view:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./viewsvr &  cd ..


run.sys:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./syssvr &  cd ..

run.ud:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cmd && nohup ./udsvr &  cd ..

run.timedjob:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd && nohup ./timedjobsvr &  cd ..

run.timedscheduler:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>run $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd cd cmd && nohup ./timedschedulersvr &  cd ..
