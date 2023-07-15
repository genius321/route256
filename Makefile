build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build
	cd notifications && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker compose up --force-recreate --build
    # docker-compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

down:
	docker compose down
  	# docker-compose down

clean-data:
	sudo rm -rf ./checkout/.pgdata
	sudo rm -rf ./loms/.pgdata
	sudo rm -rf ./notifications/.pgdata

migration-up:
	goose -dir ./loms/migrations postgres "postgres://user:password@localhost:5434/loms?sslmode=disable" up
	goose -dir ./checkout/migrations postgres "postgres://user:password@localhost:5433/checkout?sslmode=disable" up
	goose -dir ./notifications/migrations postgres "postgres://user:password@localhost:5435/notifications?sslmode=disable" up

migration-down:
	goose -dir ./loms/migrations postgres "postgres://user:password@localhost:5434/loms?sslmode=disable" down
	goose -dir ./checkout/migrations postgres "postgres://user:password@localhost:5433/checkout?sslmode=disable" down
	goose -dir ./notifications/migrations postgres "postgres://user:password@localhost:5435/notifications?sslmode=disable" down
