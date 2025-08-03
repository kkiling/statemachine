TEST_DB_NAME:=./statemachine.db


.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0
	go install github.com/pressly/goose/v3/cmd/goose@latest

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
	@echo "\n --- 🖲️ Migrate test sqlite database --- \n"
	rm -f ${TEST_DB_NAME}
	@echo "\n --- 🖲️ statemachine sqlite migrations --- \n"
	goose -dir=migrations/sqlite sqlite3 ${TEST_DB_NAME} up
	@echo "\n --- 🖲️ Creating .testenv file --- \n"
	rm -f ./testenv
	echo "SQLITE_DSN=$(CURDIR)/$(notdir ${TEST_DB_NAME})" > ./.testenv