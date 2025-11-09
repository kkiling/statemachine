TEST_DB_NAME:=statemachine_test
TEST_DB_USERNAME:=root
TEST_DB_PASSWORD:=root
DB_HOST:=localhost
DB_PORT:=5432

.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: mocks
mocks:
	@echo "\n --- 🤡 Create Mocks --- \n"
	go generate ./...

.PHONY: test
test:
	@echo "\n --- 🧪 Run project tests --- \n"
	go test ./...

.PHONY: format
format:
	@echo "\n --- 🚀 Start format imports --- \n"
	smartimports -local "github.com/kkiling/statemachine/" -path . -exclude "*_mock.go"

.PHONY: test-db
test-db:
	@echo "\n --- 🖲️ Migrate test postgresql database --- \n"
	@echo "\n --- 🖲️ Dropping and creating test database --- \n"
	# Удаляем и создаем базу данных используя правильные учетные данные
	psql postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres -c "DROP DATABASE IF EXISTS ${TEST_DB_NAME};"
	psql postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres -c "CREATE DATABASE ${TEST_DB_NAME};"
	@echo "\n --- 🖲️ Applying postgresql migrations --- \n"
	# Накатываем миграции на тестовую базу
	goose -dir=./migrations/postgresql postgres "postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" up

	@echo "\n --- 🖲️ Creating .testenv file --- \n"
	rm -f ./.testenv
	# Создаем файл с строкой подключения для PostgreSQL
	echo "POSTGRES_CONN_STRING=postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" > ./.testenv

	@echo "\n --- ✅ Test database setup completed --- \n"

.PHONY: schema
schema:
	pg_dump "postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments -s >  schema.sql
	sed -i '/^\\restrict/d' schema.sql
	sed -i '/^\\unrestrict/d' schema.sql
	sqlc generate