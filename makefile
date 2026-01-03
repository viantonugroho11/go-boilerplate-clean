migrate_up:
	migrate -source file://db/migrations -database postgres://postgres:s3cret@127.0.0.1:5432/reverse_etl?sslmode=disable up

migrate_down:
	migrate -source file://db/migrations -database postgres://postgres:s3cret@127.0.0.1:5432/reverse_etl?sslmode=disable down -all

migrate_drop:
	migrate -source file://db/migrations -database postgres://postgres:s3cret@127.0.0.1:5432/reverse_etl?sslmode=disable drop
