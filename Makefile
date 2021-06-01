export TELOXIDE_TOKEN := 1820131780:AAHYQEZynCIIyvw_nT09DRAvJUVU9zLVmr0
export DATABASE_URL := postgres://postgres:password@localhost
build:
	docker build -t st0cky-bot .
run:
	cargo run
watch:
	cargo watch -x 'run'
local-run:
	db="host=127.0.0.1 port=5432 user=postgres dbname=postgres password=password sslmode=disable" go run main.go



