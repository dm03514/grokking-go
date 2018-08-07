// build +service,!unit

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/dm03514/grokking-go/service-level-tests-in-action/deposits"
	sqsin "github.com/dm03514/grokking-go/service-level-tests-in-action/sqs"
	"testing"
	"time"
)

type DepositsTest struct {
	svc      *sqs.SQS
	messages []*sqs.Message
}

// Provision Purges all messages in the Queues and truncates all data from
// DB tables
func (dt *DepositsTest) InitializeState() error {
	if err := sqsin.PurgeQueues(dt.svc, []string{
		"http://localhost:4576/queue/deposits-in",
		"http://localhost:4576/queue/deposits-out",
	}); err != nil {
		return err
	}

	return nil
}

func (dt *DepositsTest) ApplyInput() error {
	inQueueUrl := "http://localhost:4576/queue/deposits-in"

	bs, err := json.Marshal(deposits.Deposit{
		TransactionID: 1,
		AccountNumber: 1,
		AmountCents:   1,
	})
	if err != nil {
		return err
	}

	sendResult, err := dt.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(string(bs)),
		QueueUrl:     &inQueueUrl,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Sent Message (%s)->%s, %+v\n",
		string(bs), inQueueUrl, sendResult)

	return nil
}

func NewDepositsTest() (*DepositsTest, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://localhost:4576"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &DepositsTest{
		svc: sqs.New(sess),
	}, nil
}

func Test_Service(t *testing.T) {
	test, err := NewDepositsTest()
	if err != nil {
		t.Error(err)
	}

	test.InitializeState()
	test.ApplyInput()
	messageChan := sqsin.Messages(
		context.Background(),
		test.svc,
		"http://localhost:4576/queue/deposits-out")

	// we can read after sending messages becaues we have the SQS to buffer
	// us.  Other service test types (request/response) may necessitate setting up
	// receiver handlers first

	var msg *sqs.Message
receiveloop:
	for {
		select {
		case msg = <-messageChan:
			fmt.Printf("Received Message, %+v\n", msg)
			break receiveloop

		case <-time.After(5 * time.Second):
			t.Errorf("Timeout %s reached", time.Duration(30*time.Second))
			break receiveloop
		}
	}

	var d deposits.Deposit
	if err := json.Unmarshal([]byte(*msg.Body), &d); err != nil {
		t.Error(err)
	}

	// proxy presence of postgres primary auto incremented ID for postgres
	if d.ID == 0 {
		t.Errorf("expected postgres ID received: %d. %+v", d.ID, d)
	}
}
