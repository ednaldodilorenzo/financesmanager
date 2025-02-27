migrate:
	migrate create -ext=sql -dir=migrations -seq init

migrate-db:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/dbname?sslmode=disable" -path migrations up

