build:
		-rm dist/posts/*
		go run generate.go
		cp -R static dist/static

serve:
		cd ./dist && python -m SimpleHTTPServer

dev:
		watchexec -r make all

all: build serve