default: all
all: build run

build:
	godep go build

run: 
	PORT=8888 ./stackato-redis

deploy:
	cf push

clean:
	rm -f stackato-redis
