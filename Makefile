build:
	docker build -t st0cky-bot .
run:
	docker run -it --rm --name st0cky-bot st0cky-bot
local-run:
	db="host=127.0.0.1 port=5432 user=postgres dbname=postgres password=password sslmode=disable" go run main.go
