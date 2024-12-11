# name app
APP_NAME=server
DB_URL=postgresql://root:secret@localhost:5432/goBackend?sslmode=disable

.PHONY: postgres createdb dropdb migrateup migratedown sqlc mockgen db

run: 
	go run ./cmd/${APP_NAME}/

#Định nghĩa target cho Postgress
postgres:
	docker run --name pg16_go -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

# Định nghĩa target để tạo database
createdb:
	docker exec -it pg16_go createdb --username=root --owner=root goBackend

# Định nghĩa target để xóa database
dropdb:
	docker exec -it pg16_go dropdb goBackend

# Định nghĩa target để thực hiện migration lên
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

# Định nghĩa target để thực hiện migration xuống
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
# Định nghĩa sqlc generate
sqlc: 
	sqlc generate

# Định nghĩa mockgen
mockgen:
	mockgen -package mockdb -destination db/mock/store.go hieupc05.github/backend-server/db/sqlc Store

# start db
db:
	docker start pg16_go redis_shopdev