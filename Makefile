.PHONY: setup start stop logs sqlc prefetch clean

# Get the backend container ID
BACKEND_CONTAINER := $(shell docker ps | grep fitstat-backend | awk '{print $$1}')

# Main commands
setup:
	./prefetch-images.sh
	docker-compose up --build

start:
	./../start-app.sh

stop:
	docker-compose down -v

# Logs and development commands
logs:
	docker logs -f $$(docker ps | grep fitstat-backend | awk '{print $$1}')

sqlc:
	docker exec -it $$(docker ps | grep fitstat-backend | awk '{print $$1}') sqlc generate

# Utility commands
prefetch:
	./prefetch-images.sh

clean: stop
	docker system prune -f