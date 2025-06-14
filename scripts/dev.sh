set -o allexport; source .env; \
go run cmd/ufl_event_processor/main.go; \
set +o allexport