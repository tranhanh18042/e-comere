projectDir:=$(shell pwd)
buildDir:=$(projectDir)/build
servicesDir:=$(projectDir)/services

tidy-packages:
	cd $(servicesDir) && go mod tidy

build-app:
	cd $(servicesDir) && go build -o $(buildDir)/service . && chmod +x $(buildDir)/service

run-app:
	cd $(buildDir) && ./service;

build-and-run: build-app run-app

build-and-run-docker:
	docker-compose up --build

start-all-docker:
	docker-compose up

stop-all-docker:
	docker-compose stop

start-service-docker: build-app
	docker-compose up --build $(SVC)

start-service: build-app
	ECOM_SERVICE=$(SVC) ./build/service
