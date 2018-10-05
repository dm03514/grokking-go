BINARY_NAME=deposits-service
BIN_DIR=bin/
TEST_SQS_ENDPOINT_URL=http://localhost:4576

AWS_CONFIG_FILE=$(shell pwd)/config/aws.config
AWS_SHARED_CREDENTIALS_FILE=$(shell pwd)/config/aws.credentials
export AWS_CONFIG_FILE
export AWS_SHARED_CREDENTIALS_FILE

build:
	go build -o $(BIN_DIR)$(BINARY_NAME) -v

start-service-stack:
	docker-compose down && SERVICES=sqs docker-compose up -d
	wait-for --poll-interval 1s postgres \
		--connection-string="postgresql://root:root@localhost/deposits?sslmode=disable"
	wait-for net --address="localhost:4576"

start-service:
	aws --endpoint-url=$(TEST_SQS_ENDPOINT_URL) sqs delete-queue --queue-url $(TEST_SQS_ENDPOINT_URL)/queue/deposits-in || true
	aws --endpoint-url=$(TEST_SQS_ENDPOINT_URL) sqs create-queue --queue-name deposits-in

	aws --endpoint-url=$(TEST_SQS_ENDPOINT_URL) sqs delete-queue --queue-url $(TEST_SQS_ENDPOINT_URL)/queue/deposits-out || true
	aws --endpoint-url=$(TEST_SQS_ENDPOINT_URL) sqs create-queue --queue-name deposits-out

	AWS_SECRET_ACCESS_KEY=x \
	AWS_ACCESS_KEY_ID=x \
	AWS_REGION=us-west-2 \
	./$(BIN_DIR)$(BINARY_NAME) \
		--sqs-endpoint-url=$(TEST_SQS_ENDPOINT_URL) \
		--db-connection-string="postgresql://root:root@localhost/deposits?sslmode=disable" \
		--sqs-input-url=$(TEST_SQS_ENDPOINT_URL)/queue/deposits-in \
		--sqs-output-url=$(TEST_SQS_ENDPOINT_URL)/queue/deposits-out

test-service:
	AWS_SECRET_ACCESS_KEY=x \
	AWS_ACCESS_KEY_ID=x \
	AWS_REGION=us-west-2 \
	go test -tags=service ./... -v


.PHONY: build start-service-stack test-service