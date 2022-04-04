run:
	go run main.go

postgres:
	docker run --name=first-task-local-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=postgres -p 5433:5432 -d postgres:14