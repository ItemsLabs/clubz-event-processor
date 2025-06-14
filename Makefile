dev:
	air
run:
	@set -o allexport; source .env
	go run .
generate:
	@set -o allexport; source .env
	@(PSQL_USER=${DATABASE_USER} \
		PSQL_HOST=${DATABASE_HOST} \
		PSQL_PORT=${DATABASE_PORT} \
		PSQL_DBNAME=${DATABASE_NAME} \
		PSQL_PASS=${DATABASE_PASSWORD} \
		PSQL_SSLMODE=${DATABASE_SSLMODE} \
		go generate ./...)