LOCAL_BIN ?= $(CURDIR)/bin
PG_TOOL_BIN := $(shell command -v oh-my-pg-tool 2>/dev/null)

MIGRATION_FOLDER ?= apps/api/scripts/migrations
MIGRATE_DSN ?= postgres://pulse:pulse@localhost:5432/running_club?sslmode=disable
MIGRATE_DSN_TEST ?= postgres://pulse:pulse@localhost:5432/running_club_test?sslmode=disable

.pg_tool_bin:
	@if [ -z "$(PG_TOOL_BIN)" ]; then \
		echo "oh-my-pg-tool not found in PATH. Install it and retry."; \
		exit 1; \
	fi

migration: .pg_tool_bin
	@(printf "Enter name of migration: "; read arg; $(PG_TOOL_BIN) goose create $$arg -d $(MIGRATION_FOLDER))

migrate-up: .pg_tool_bin
	$(PG_TOOL_BIN) local -c "$(MIGRATE_DSN)" -d $(MIGRATION_FOLDER)

migrate-up-test: .pg_tool_bin
	$(PG_TOOL_BIN) local -c "$(MIGRATE_DSN_TEST)" -d $(MIGRATION_FOLDER)

migrate-up-all: migrate-up migrate-up-test

migrate-down: .pg_tool_bin
	$(PG_TOOL_BIN) local -c "$(MIGRATE_DSN)" -d $(MIGRATION_FOLDER) --command=down

migrate-down-test: .pg_tool_bin
	$(PG_TOOL_BIN) local -c "$(MIGRATE_DSN_TEST)" -d $(MIGRATION_FOLDER) --command=down

migrate-down-all: migrate-down migrate-down-test

migrate-status: .pg_tool_bin
	$(PG_TOOL_BIN) local -c "$(MIGRATE_DSN)" -d $(MIGRATION_FOLDER) --command=status

.PHONY: migration migrate-up migrate-up-test migrate-up-all migrate-down migrate-down-test migrate-down-all migrate-status .pg_tool_bin
