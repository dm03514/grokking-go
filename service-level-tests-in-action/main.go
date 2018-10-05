package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/dm03514/grokking-go/service-level-tests-in-action/deposits"
	sqsin "github.com/dm03514/grokking-go/service-level-tests-in-action/sqs"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr // Can be any io.Writer, see below for File example log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

const postgresDuplicateKeyErrorCode = "23505"

type MessageProcessor struct {
	svc         *sqs.SQS
	db          *sql.DB
	inQueueURL  string
	outQueueURL string
}

func (p MessageProcessor) Process(msg *sqs.Message) error {
	log.WithFields(log.Fields{
		"sqs.message": msg,
	}).Debug("received_message")

	// do some work
	// ie write to database
	var d deposits.Deposit
	if err := json.Unmarshal([]byte(*msg.Body), &d); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"deposit": d,
	}).Debug("deserialized_deposit")

	err := p.db.QueryRow(
		`INSERT 
					INTO deposits(transaction_id, account_number, amount_cents)
					VALUES($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`,
		d.TransactionID, d.AccountNumber, d.AmountCents).Scan(&d.ID)

	log.WithFields(log.Fields{
		"deposit": d,
	}).Debug("inserted_deposit")

	if err, ok := err.(*pq.Error); ok && err.Code != postgresDuplicateKeyErrorCode {
		return err
	} else if err != nil {
		return err
	}

	bs, err := json.Marshal(d)
	if err != nil {
		return err
	}

	sendResult, err := p.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(string(bs)),
		QueueUrl:     aws.String(p.outQueueURL),
		DelaySeconds: aws.Int64(0),
	})

	if err != nil {
		return err
	}

	// write output message
	log.WithFields(log.Fields{
		"deposit_sent": d,
		"sqs.message":  *sendResult.MessageId,
		"output_queue": p.outQueueURL,
	}).Debug("sent_message")

	// delete message
	resultDelete, err := p.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(p.inQueueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"sqs.message": resultDelete,
	}).Debug("message_deleted")

	return nil
}

func main() {
	var sqsEndpointURL = flag.String("sqs-endpoint-url", "", "")
	var inQueueURL = flag.String("sqs-input-url", "", "")
	var outQueueURL = flag.String("sqs-output-url", "", "")
	var dbConnectionString = flag.String("db-connection-string", "", "")
	flag.Parse()

	log.WithFields(log.Fields{
		"sqs_endpoint_url":     sqsEndpointURL,
		"sqs_in_queue_url":     inQueueURL,
		"sqs_out_queue_url":    outQueueURL,
		"db_connection_string": dbConnectionString,
	}).Info("initializing")

	// connect to postgres
	db, err := sql.Open("postgres", *dbConnectionString)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	// connect to sqs
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(*sqsEndpointURL),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	mp := MessageProcessor{
		svc:         svc,
		db:          db,
		inQueueURL:  *inQueueURL,
		outQueueURL: *outQueueURL,
	}

	// read a message
	inmessages := sqsin.Messages(context.Background(), svc, *inQueueURL)

	log.Debugf("Entering Receive Message Loop")
	for {
		select {
		case msg := <-inmessages:
			if err := mp.Process(msg); err != nil {
				panic(err)
			}
		}
	}
}
