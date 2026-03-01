.PHONY: up up-logs down down-volumes stop test-up test-up-logs test-down test-down-volumes test-stop pg-migrate pg-info pg-validate pg-clean pg-repair pg-baseline

## docker commands
up:
	docker-compose -f docker-compose-dev.yml up -d
up-logs:
	docker-compose -f docker-compose-dev.yml up
down:
	docker-compose -f docker-compose-dev.yml down
down-volumes:
	docker-compose -f docker-compose-dev.yml down --volumes --remove-orphans
stop:
	docker-compose -f docker-compose-dev.yml stop

## test docker commands
test-up:
	docker-compose -f docker-compose-test.yml up -d
test-up-logs:
	docker-compose -f docker-compose-test.yml up
test-down:
	docker-compose -f docker-compose-test.yml down
test-down-volumes:
	docker-compose -f docker-compose-test.yml down --volumes --remove-orphans
test-stop:
	docker-compose -f docker-compose-test.yml stop

## PostgreSQL commands
pg-migrate: 
	docker compose -f docker-compose-dev.yml run --rm flyway migrate

pg-info: 
	docker compose -f docker-compose-dev.yml run --rm flyway info

pg-validate: 
	docker compose -f docker-compose-dev.yml run --rm flyway validate

pg-clean:
	@echo "WARNING: This will delete ALL data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose -f docker-compose-dev.yml run --rm flyway clean; \
	fi

pg-repair: 
	docker compose -f docker-compose-dev.yml run --rm flyway repair

pg-baseline: 
	docker compose -f docker-compose-dev.yml run --rm flyway baseline