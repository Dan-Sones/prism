.PHONY: up-session down-session stop-session up-application down-application stop-application

up-session:
	docker-compose -f docker-compose-session.yml up -d

down-session:
	docker-compose -f docker-compose-session.yml down --volumes --remove-orphans

stop-session:
	docker-compose -f docker-compose-session.yml stop

up-application:
	docker-compose -f docker-compose-application.yml up -d

down-application:
	docker-compose -f docker-compose-application.yml down --volumes --remove-orphans

stop-application:	
	docker-compose -f docker-compose-application.yml down --volumes --remove-orphans