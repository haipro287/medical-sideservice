dev-recreate:
	@docker-compose --project-name=medical-chain-dev -f deploy/dev/docker-compose-dev.yaml up -d --build --force-recreate

#dev-migration-up:
#	NETWORK_NAME=medical-chain-server-dev docker-compose --project-name=medical-chain-dev -f deploy/dev/docker-compose-migration-tool.yaml up up
#
#dev-migration-down:
#	NETWORK_NAME=medical-chain-server-dev docker-compose --project-name=medical-chain-dev -f deploy/dev/docker-compose-migration-tool.yaml up down
build-and-push-image: build-image push-image

build-image:
	@docker build . --target=release -t supermedicalchain/side-service:pre-release

push-image:
	@docker tag supermedicalchain/side-service:pre-release supermedicalchain/side-service${TAG}
	@docker push supermedicalchain/side-service${TAG}

build:
	@go build -o ./dist/server ./src

serve:
	@./dist/server serve

dev:
	@./dist/server  --log-format plain --log-level debug --disable-profiler --allow-kill serve

cleanDB:
#	@./dist/server clean
	@echo "Hello"

seed:
	@#./dist/server seed-data --clean
	@echo "Hello"

test:
	#go test ./src/cockroach/... -v -check.f "CockroachDbGraphTestSuite.*"
	@go test ./... -v

test-prepare-up:
	@#docker exec  up -f deploy/dev/docker-compose.yaml side-cdb -d
	@echo "Hello"

test-prepare-down:
	@#docker-compose down -f deploy/dev/docker-compose.yaml side-cdb
	@echo "Hello"
grpc-client:
	@grpc-client-cli localhost:${GRPC_PORT}

kill:
	@(echo '{}' | grpc-client-cli -service CommonService -method Kill localhost:${GRPC_PORT}) > /nil 2> /nil || return 0

proto-gen:
	@./genproto.sh
