LOAD_TEST_RATE=500
LOAD_TEST_TARGET=http://localhost:8080

stack:
	docker-compose down && docker-compose up

load-test:
	echo "GET $(LOAD_TEST_TARGET)/slow" | vegeta attack -rate=$(LOAD_TEST_RATE) -duration=0 | tee results.bin | vegeta report

.PHONY: stack load-test
