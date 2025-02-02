.PHONY: setup start stop logs sqlc prefetch clean

SCRIPTS_DIR := ./scripts
BACKEND_CONTAINER := $(shell docker ps | grep reppy-backend | awk '{print $$1}')


# Main commands
setup:
	$(SCRIPTS_DIR)/prefetch-images.sh
	docker-compose up --build

start:
	$(SCRIPTS_DIR)/start-app.sh

stop:
	docker-compose down -v


# Logs and development commands
logs:
	docker logs -f $(BACKEND_CONTAINER)

sqlc:
	docker exec -it $(BACKEND_CONTAINER) sqlc generate


# Utility commands
prefetch:
	$(SCRIPTS_DIR)/prefetch-images.sh

clean: stop
	docker system prune -f