default: all
all: build run

build:
	godep go build

run: 
	PORT=8888 ./stackato-mongo

deploy:
	cf push

clean:
	rm -f stackato-mongo
