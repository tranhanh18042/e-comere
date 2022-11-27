projectDir:=$(shell pwd)
buildDir:=$(projectDir)/build
servicesDir:=$(projectDir)/services

# change docker compose command based on the OS
os:=$(shell uname -a | egrep Darwin)
ifeq ($(shell uname -a | egrep Darwin),) # MacOS
	dkpcmnd:="docker compose"
else # Linux
	dkpcmnd:="docker-compose"
endif

tidy-packages:
	cd $(servicesDir) && go mod tidy

build-app:
	cd $(servicesDir) && go build -o $(buildDir)/service . && chmod +x $(buildDir)/service

run-app:
	cd $(buildDir) && ./service;

build-and-run: build-app run-app

build-and-run-docker:
	$(dkpcmnd) up --build

start-all-docker:
	echo "$(os)"
	$(dkpcmnd) up

stop-all-docker:
	$(dkpcmnd) stop

start-service-docker: build-app
	$(dkpcmnd) up --build $(SVC)

start-service: build-app
	ECOM_SERVICE=$(SVC) ./build/service