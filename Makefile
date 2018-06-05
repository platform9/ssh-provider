# Copyright 2018 Platform 9 Systems, Inc.

.PHONY: image push dev_image dev_push

PREFIX = "platform9"
NAME = "ssh-provider"
TAG ?= $(shell git describe --tags)

image:
	docker build -t "$(PREFIX)/$(NAME):$(TAG)" -f ./Dockerfile .

push: image
	docker push "$(PREFIX)/$(NAME):$(TAG)"

dev_image:
	docker build -t "$(PREFIX)/$(NAME):$(TAG)-dev" -f ./Dockerfile .

dev_push: dev_image
	docker push "$(PREFIX)/$(NAME):$(TAG)-dev"