



Executing

- Install python requirements

```
$ pip install -r requirements.txt
```

- Install [wait-for](https://github.com/dm03514/wait-for)
```
$ git clone git@github.com:dm03514/wait-for.git
$ cd wait-for && GO111MODULE=on go install
```

- Build

```
$ make build
go build -o bin/deposits-service -v
```

- Start Test Dependencies

```
$ make start-service-stack
make: Warning: File 'Makefile' has modification time 231 s in the future
docker-compose down && SERVICES=sqs docker-compose up -d
Stopping service-level-tests-in-action_localstack_1 ... done
Stopping service-level-tests-in-action_postgres_1   ... done
Removing service-level-tests-in-action_localstack_1 ... done
Removing service-level-tests-in-action_postgres_1   ... done
Removing network service-level-tests-in-action_default
WARNING: The Docker Engine you're using is running in swarm mode.

Compose does not use swarm mode to deploy services to multiple nodes in a swarm. All containers will be scheduled on the current node.

To deploy your application across the swarm, use `docker stack deploy`.

Creating network "service-level-tests-in-action_default" with the default driver

Creating service-level-tests-in-action_postgres_1   ... done
Creating service-level-tests-in-action_localstack_1 ... done
wait-for --poll-interval 1s postgres \
        --connection-string="postgresql://root:root@localhost/deposits?sslmode=disable"
{"level":"info","msg":"polling","time":"2018-09-01T01:25:41Z"}
{"err":"read tcp 127.0.0.1:34240-\u003e127.0.0.1:5432: read: connection reset by peer","level":"debug","msg":"poll_result","ready":false,"time":"2018-09-01T01:25:41Z"}
{"level":"info","msg":"polling","time":"2018-09-01T01:25:42Z"}
{"err":null,"level":"debug","msg":"poll_result","ready":true,"time":"2018-09-01T01:25:42Z"}
wait-for net --address="localhost:4576"
{"level":"info","msg":"polling","time":"2018-09-01T01:25:42Z"}
{"error":null,"level":"debug","local_addr":{"IP":"127.0.0.1","Port":60140,"Zone":""},"module":"poller.Net","msg":"conn_response","remote_addr":{"IP":"127.0.0.1","Port":4576,"Zone":""},"time":"2018-09-01T01:25
:42Z"}
{"err":null,"level":"debug","msg":"poll_result","ready":true,"time":"2018-09-01T01:25:42Z"}
```

- Start Service 

```
$ make start-service
  make: Warning: File 'Makefile' has modification time 165 s in the future
  aws --endpoint-url=http://localhost:4576 sqs delete-queue --queue-url http://localhost:4576/queue/deposits-in || true
  
  An error occurred (AWS.SimpleQueueService.NonExistentQueue) when calling the DeleteQueue operation: AWS.SimpleQueueService.NonExistentQueue; see the SQS docs.
  aws --endpoint-url=http://localhost:4576 sqs create-queue --queue-name deposits-in
  {
      "QueueUrl": "http://localhost:4576/queue/deposits-in"
  }
  aws --endpoint-url=http://localhost:4576 sqs delete-queue --queue-url http://localhost:4576/queue/deposits-out || true
  
  An error occurred (AWS.SimpleQueueService.NonExistentQueue) when calling the DeleteQueue operation: AWS.SimpleQueueService.NonExistentQueue; see the SQS docs.
  aws --endpoint-url=http://localhost:4576 sqs create-queue --queue-name deposits-out
  {
      "QueueUrl": "http://localhost:4576/queue/deposits-out"
  }
  AWS_SECRET_ACCESS_KEY=x \
  AWS_ACCESS_KEY_ID=x \
  AWS_REGION=us-west-2 \
  ./bin/deposits-service \
          --sqs-endpoint-url=http://localhost:4576 \
          --db-connection-string="postgresql://root:root@localhost/deposits?sslmode=disable" \ 
          --sqs-output-url=http://localhost:4576/queue/deposits-out
          
  {"db_connection_string":"postgresql://root:root@localhost/deposits?sslmode=disable","level":"info","msg":"initializing","sqs_endpoint_url":"http://localhost:4576","sqs_in_queue_url":"http://localhost:4576/que
  ue/deposits-in","sqs_out_queue_url":"http://localhost:4576/queue/deposits-out","time":"2018-09-01T01:26:47Z"}
  {"level":"debug","msg":"Entering Receive Message Loop","time":"2018-09-01T01:26:47Z"}
  ...
  {"level":"debug","messages":{"Messages":null},"msg":"receiving_messages","time":"2018-09-01T01:34:19Z"}
  ...
```

- Execute Test

```
 make test-service                                                                [11/2061]
AWS_SECRET_ACCESS_KEY=x \
AWS_ACCESS_KEY_ID=x \
AWS_REGION=us-west-2 \
go test -tags=service ./... -v
=== RUN   Test_Service
Purged Queue (http://localhost:4576/queue/deposits-in): {

}
Purged Queue (http://localhost:4576/queue/deposits-out): {

}
Sent Message ({"transaction_id":1,"account_number":1,"amount_cents":1})->http://localhost:4576/queue/deposits-in, {
  MD5OfMessageAttributes: "d41d8cd98f00b204e9800998ecf8427e",
  MD5OfMessageBody: "eed62631cdf7fb97f1bdc09937f7d1aa",
  MessageId: "25db2c51-dc3c-4b2b-b84c-cbba180de7e0"
}
{"level":"debug","messages":{"Messages":[{"Attributes":{"SentTimestamp":"1535765666884"},"Body":"{\"id\":1,\"transaction_id\":1,\"account_number\":1,\"amount_cents\":1}","MD5OfBody":"693bcb12275ce7875ef923263
6fd6f5f","MD5OfMessageAttributes":null,"MessageAttributes":null,"MessageId":"1e562db5-7556-40bd-9503-9947c1beb98d","ReceiptHandle":"1e562db5-7556-40bd-9503-9947c1beb98d#1f9059f4-a585-41ce-9f9d-815c5a2f8f4c"}]
},"msg":"receiving_messages","time":"2018-09-01T01:34:26Z"}
Received Message, {
  Attributes: {
    SentTimestamp: "1535765666884"
  },
  Body: "{\"id\":1,\"transaction_id\":1,\"account_number\":1,\"amount_cents\":1}",
  MD5OfBody: "693bcb12275ce7875ef9232636fd6f5f",
  MessageId: "1e562db5-7556-40bd-9503-9947c1beb98d",
  ReceiptHandle: "1e562db5-7556-40bd-9503-9947c1beb98d#1f9059f4-a585-41ce-9f9d-815c5a2f8f4c"
}
--- PASS: Test_Service (0.24s)
PASS
ok      github.com/dm03514/grokking-go/service-level-tests-in-action    0.243s
?       github.com/dm03514/grokking-go/service-level-tests-in-action/deposits   [no test files]
?       github.com/dm03514/grokking-go/service-level-tests-in-action/sqs        [no test files]
```

- Service Logs during Test

```
{"level":"debug","messages":{"Messages":null},"msg":"receiving_messages","time":"2018-09-01T01:34:26Z"}
{"level":"debug","messages":{"Messages":[{"Attributes":{"SentTimestamp":"1535765666795"},"Body":"{\"transaction_id\":1,\"account_number\":1,\"amount_cents\":1}","MD5OfBody":"eed62631cdf7fb97f1bdc09937f7d1aa",
"MD5OfMessageAttributes":null,"MessageAttributes":null,"MessageId":"25db2c51-dc3c-4b2b-b84c-cbba180de7e0","ReceiptHandle":"25db2c51-dc3c-4b2b-b84c-cbba180de7e0#d8c706b7-4dae-45b0-8b8e-2c248499c763"}]},"msg":"
receiving_messages","time":"2018-09-01T01:34:26Z"}
{"level":"debug","msg":"received_message","sqs.message":{"Attributes":{"SentTimestamp":"1535765666795"},"Body":"{\"transaction_id\":1,\"account_number\":1,\"amount_cents\":1}","MD5OfBody":"eed62631cdf7fb97f1b
dc09937f7d1aa","MD5OfMessageAttributes":null,"MessageAttributes":null,"MessageId":"25db2c51-dc3c-4b2b-b84c-cbba180de7e0","ReceiptHandle":"25db2c51-dc3c-4b2b-b84c-cbba180de7e0#d8c706b7-4dae-45b0-8b8e-2c248499c
763"},"time":"2018-09-01T01:34:26Z"}
{"deposit":{"transaction_id":1,"account_number":1,"amount_cents":1},"level":"debug","msg":"deserialized_deposit","time":"2018-09-01T01:34:26Z"}
{"deposit":{"id":1,"transaction_id":1,"account_number":1,"amount_cents":1},"level":"debug","msg":"inserted_deposit","time":"2018-09-01T01:34:26Z"}
{"deposit_sent":{"id":1,"transaction_id":1,"account_number":1,"amount_cents":1},"level":"debug","msg":"sent_message","output_queue":"http://localhost:4576/queue/deposits-out","sqs.message":"1e562db5-7556-40bd
-9503-9947c1beb98d","time":"2018-09-01T01:34:26Z"}
{"level":"debug","msg":"message_deleted","sqs.message":{},"time":"2018-09-01T01:34:26Z"}
{"level":"debug","messages":{"Messages":null},"msg":"receiving_messages","time":"2018-09-01T01:34:27Z"}
```