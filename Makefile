include Makefile.config
include Makefile.docker
include Makefile.test
include Makefile.dev
include Makefile.docs
include Makefile.clean

.PHONY: help status

help:
	@echo "$(APP_NAME) v$(VERSION) - Available commands:"
	@echo ""
	@echo "BUILD & RUN:"
	@echo "  make build        - Build project"
	@echo "  make run          - Run project"
	@echo "  make run-detached - Run in background"
	@echo "  make stop         - Stop project"
	@echo "  make logs         - View logs"
	@echo "  make health       - Health check"
	@echo "  make quick-run    - Quick run"
	@echo ""
	@echo "TESTING:"
	@echo "  make test         - All tests"
	@echo "  make test-service - Service tests"
	@echo "  make test-repo    - Repository tests"
	@echo "  make test-cover   - Tests with coverage"
	@echo "  make test-team    - TeamService tests"
	@echo "  make test-pr      - PRService tests"
	@echo "  make load-test    - Load testing"
	@echo "  make quick-test   - Quick tests"
	@echo ""
	@echo "DEVELOPMENT:"
	@echo "  make dev          - Local development"
	@echo "  make run-local    - Local run"
	@echo "  make migrate      - Run migrations"
	@echo "  make lint         - Code linting"
	@echo "  make fmt          - Code formatting"
	@echo "  make tidy         - Dependency analysis"
	@echo ""
	@echo "DOCUMENTATION:"
	@echo "  make swagger      - Generate Swagger"
	@echo "  make run-swagger  - Run with Swagger"
	@echo "  make dev-swagger  - Dev with Swagger"
	@echo ""
	@echo "CLEANUP:"
	@echo "  make clean        - Full cleanup"
	@echo "  make fresh        - Full rebuild"
	@echo ""
	@echo "STATUS:"
	@echo "  make status       - Service status"

status:
	@echo "Service status:"
	docker-compose ps
	@echo ""
	@echo "API check:"
	@curl -s http://localhost:$(PORT)/health | jq . 2>/dev/null || curl -s http://localhost:$(PORT)/health