.PHONY: seed-db

seed-db:
	docker exec -i prism_postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < pg-test-data.sql
