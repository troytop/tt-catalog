default: all
all: build run

build:
	godep go build

run: 
	PORT=8888 ./stackato-mysql

deploy:
	cf push

clean:
	rm -f stackato-mysql
