-include .env

.PHONY: up up-logs down down-volumes stop test-up test-up-logs test-up-tls test-up-tls-logs test-down test-down-volumes test-stop test-stop-tls pg-migrate pg-info pg-validate pg-clean pg-repair pg-baseline build-push build-experimentation-service build-assignment-service build-clickhouse-writer build-events-service certs-host certs-ip certs-clean certs-trust certs-untrust

CADDY_CERT_DIR ?= ./caddy/certs
PRISM_PUBLIC_HOSTNAME ?= localhost
UNAME_S := $(shell uname -s)

## docker commands
up:
	docker compose -f docker-compose-dev.yml up -d
up-logs:
	docker compose -f docker-compose-dev.yml up
down:
	docker compose -f docker-compose-dev.yml down
down-volumes:
	docker compose -f docker-compose-dev.yml down --volumes --remove-orphans
stop:
	docker compose -f docker-compose-dev.yml stop

## test docker commands
test-up:
	docker compose --env-file .env.test -f docker-compose-test.yml up -d
test-up-logs:
	docker compose --env-file .env.test -f docker-compose-test.yml up
test-up-tls:
	docker compose --env-file .env.test -f docker-compose-test.yml --profile tls up -d
test-up-tls-logs:
	docker compose --env-file .env.test -f docker-compose-test.yml --profile tls up
test-down:
	docker compose --env-file .env.test -f docker-compose-test.yml down
test-down-volumes:
	docker compose --env-file .env.test -f docker-compose-test.yml down --volumes --remove-orphans
test-stop:
	docker compose --env-file .env.test -f docker-compose-test.yml stop
test-stop-tls:
	docker compose --env-file .env.test -f docker-compose-test.yml --profile tls stop

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

certs-host:
	@mkdir -p $(CADDY_CERT_DIR)
	openssl req -x509 -newkey rsa:2048 -nodes \
		-keyout $(CADDY_CERT_DIR)/tls.key \
		-out $(CADDY_CERT_DIR)/tls.crt \
		-days 365 -subj "/CN=$(PRISM_PUBLIC_HOSTNAME)" \
		-addext "subjectAltName=DNS:$(PRISM_PUBLIC_HOSTNAME)"
	@echo "Generated self-signed cert for DNS:$(PRISM_PUBLIC_HOSTNAME) in $(CADDY_CERT_DIR)"

certs-ip:
	@mkdir -p $(CADDY_CERT_DIR)
	openssl req -x509 -newkey rsa:2048 -nodes \
		-keyout $(CADDY_CERT_DIR)/tls.key \
		-out $(CADDY_CERT_DIR)/tls.crt \
		-days 365 -subj "/CN=$(PRISM_PUBLIC_HOSTNAME)" \
		-addext "subjectAltName=IP:$(PRISM_PUBLIC_HOSTNAME)"
	@echo "Generated self-signed cert for IP:$(PRISM_PUBLIC_HOSTNAME) in $(CADDY_CERT_DIR)"

certs-clean:
	rm -f $(CADDY_CERT_DIR)/tls.crt $(CADDY_CERT_DIR)/tls.key

certs-trust:
ifeq ($(UNAME_S),Darwin)
	sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $(CADDY_CERT_DIR)/tls.crt
	@echo "Added $(CADDY_CERT_DIR)/tls.crt to the System keychain"
else
	sudo cp $(CADDY_CERT_DIR)/tls.crt /usr/local/share/ca-certificates/prism-caddy.crt
	sudo update-ca-certificates
	@echo "Added cert to /usr/local/share/ca-certificates and refreshed the trust store"
endif

certs-untrust:
ifeq ($(UNAME_S),Darwin)
	sudo security remove-trusted-cert -d $(CADDY_CERT_DIR)/tls.crt
	@echo "Removed cert from the System keychain"
else
	sudo rm -f /usr/local/share/ca-certificates/prism-caddy.crt
	sudo update-ca-certificates --fresh
	@echo "Removed cert from /usr/local/share/ca-certificates and refreshed the trust store"
endif
