build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build
	cd notifications && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker compose up --force-recreate --build
    # docker-compose up --force-recreate --build

down:
	docker compose down
  	# docker-compose down

clean-data:
	sudo rm -rf ./checkout/pgdata
	sudo rm -rf ./loms/pgdata

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit