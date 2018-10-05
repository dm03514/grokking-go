package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// PurgeQueues uses the aws client library to purge all messages
func PurgeQueues(svc *sqs.SQS, queueURLs []string) error {
	for _, qURL := range queueURLs {
		req, resp := svc.PurgeQueueRequest(&sqs.PurgeQueueInput{
			QueueUrl: aws.String(qURL),
		})

		err := req.Send()
		if err != nil { // resp is now filled
			return err
		}
		fmt.Printf("Purged Queue (%s): %+v\n", qURL, resp)
	}
	return nil
}
