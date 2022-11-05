projectDir:=$(shell pwd)
buildDir:=$(projectDir)/build
servicesDir:=$(projectDir)/services

tidy-packages:
	cd $(servicesDir)/item && go mod tidy && \
	cd $(servicesDir)/order && go mod tidy && \
	cd $(servicesDir)/user && go mod tidy

build-packages:
	cd $(servicesDir)/item && go build -o $(buildDir)/item . && chmod +x $(buildDir)/item && \
	cd $(servicesDir)/order && go build -o $(buildDir)/order . && chmod +x $(buildDir)/item && \
	cd $(servicesDir)/user && go build -o $(buildDir)/user . && chmod +x $(buildDir)/item

run-services: build-packages
	cd $(buildDir) && ./item;
	cd $(buildDir) && ./order;
	cd $(buildDir) && ./user;
