package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

func Messages(ctx context.Context, svc *sqs.SQS, queueURL string) <-chan *sqs.Message {
	received := make(chan *sqs.Message)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
					QueueUrl: aws.String(queueURL),
					AttributeNames: aws.StringSlice([]string{
						"SentTimestamp",
					}),
					MaxNumberOfMessages:   aws.Int64(1),
					MessageAttributeNames: []*string{aws.String("All")},
					WaitTimeSeconds:       aws.Int64(1),
				})
				if err != nil {
					log.Errorf("Unable to receive message from queue %q, %v.", queueURL, err)
					return
				}

				log.WithFields(log.Fields{
					"messages": result,
				}).Debug("receiving_messages")

				for _, m := range result.Messages {
					received <- m
				}

			}
		}
	}()
	return received
}
