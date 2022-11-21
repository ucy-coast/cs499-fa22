.PHONY: proto run push

proto:
	for f in internal/*/proto/*.proto; do \
		protoc --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done

run:
	docker-compose -f docker-compose.yml build
	docker-compose -f docker-compose.yml up --remove-orphans

push:
	docker-compose -f docker-compose.yml build
	docker-compose -f docker-compose.yml push
