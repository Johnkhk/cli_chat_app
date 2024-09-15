# Makefile for running MySQL migrations with .env support

DB_USER=cli_chat_dev
DB_NAME=cli_chat_app
UP_MIGRATION=db/migrations/up.sql
DOWN_MIGRATION=db/migrations/down.sql
UI_TEST=db/migrations/test.sql

.PHONY: up down

# Target for running the up migration
up:
	export MYSQL_PASSWORD=$$(grep MYSQL_PASSWORD .env | cut -d '=' -f2) && mysql -u $(DB_USER) -p$$MYSQL_PASSWORD $(DB_NAME) < $(UP_MIGRATION)

# Target for running the down migration
down:
	export MYSQL_PASSWORD=$$(grep MYSQL_PASSWORD .env | cut -d '=' -f2) && mysql -u $(DB_USER) -p$$MYSQL_PASSWORD $(DB_NAME) < $(DOWN_MIGRATION)

ui_test:
	export MYSQL_PASSWORD=$$(grep MYSQL_PASSWORD .env | cut -d '=' -f2) && mysql -u $(DB_USER) -p$$MYSQL_PASSWORD $(DB_NAME) < $(UI_TEST)
# Target for cleaning JWT tokens
clean:
	rm -f $(HOME)/$(APP_DIR_NAME)/jwt_tokens

# Target for tailing the client debug log
tail-log:
	tail -f $(HOME)/$(APP_DIR_NAME)/debug.log

# Target to run main.go with a numbered APP_DIR_NAME
# e.g make run-user USER_NUM=1
run-user:
	@if [ -z "$(USER_NUM)" ]; then \
		echo "Please provide a USER_NUM, e.g., 'make run-user USER_NUM=1'"; \
	else \
		APP_DIR_NAME=".cli_chat_app_user$(USER_NUM)" go run cmd/client/main.go; \
	fi

# Target to tail the user log
# e.g make tail-user-log USER_NUM=1
tail-user-log:
	@if [ -z "$(USER_NUM)" ]; then \
		echo "Please provide a USER_NUM, e.g., 'make tail-user-log USER_NUM=1'"; \
	else \
		tail -f $(HOME)/.cli_chat_app_user$(USER_NUM)/debug.log; \
	fi


# Target to clean user JWT tokens
# e.g make clean-user-token USER_NUM=1
clean-user-token:
	@if [ -z "$(USER_NUM)" ]; then \
		echo "Please provide a USER_NUM, e.g., 'make clean-user USER_NUM=1'"; \
	else \
		rm -f $(HOME)/.cli_chat_app_user$(USER_NUM)/jwt_tokens; \
	fi