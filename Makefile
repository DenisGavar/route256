
build-all:
	cd checkout && GOOS=linux make build
	cd loms && GOOS=linux make build
	cd notifications && GOOS=linux make build

run-all: build-all
	sudo docker compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

install-go-deps-all:
	cd checkout && make install-go-deps
	cd loms && make install-go-deps

get-go-deps-all:
	cd checkout && make get-go-deps
	cd loms && make get-go-deps

vendor-proto-all:
	cd checkout && make vendor-proto
	cd loms && make vendor-proto

generate-all:
	cd checkout && make generate
	cd loms && make generate	