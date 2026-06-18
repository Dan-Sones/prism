.PHONY: run-prism run-powell run-api run-website install stop-prism

install:
	cd powell-motors-website && npm install

run-api:
	cd powell-motors-api-springboot && ./mvnw spring-boot:run

run-website:
	cd powell-motors-website && npm run dev

run-prism:
	docker compose up -d

stop-prism:
	docker compose down

run-powell:
	@echo "Starting powell-motors-api-springboot and powell-motors-website..."
	@trap 'kill 0' INT TERM EXIT; \
		$(MAKE) run-api & \
		$(MAKE) run-website & \
		wait




