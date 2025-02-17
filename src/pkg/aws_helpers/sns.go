package aws_helpers

import (
	"context"
	// "encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// SnsActions encapsulates the Amazon Simple Notification Service (Amazon SNS) actions
// used in the examples.
type SnsActions struct {
	SnsClient *sns.Client
}

// Publish publishes a message to an Amazon SNS topic. The message is then sent to all
// subscribers. When the topic is a FIFO topic, the message must also contain a group ID
// and, when ID-based deduplication is used, a deduplication ID. An optional key-value
// filter attribute can be specified so that the message can be filtered according to
// a filter policy.
func (actor SnsActions) Publish(ctx context.Context, topicArn string, message string, groupId string, dedupId string, filterKey string, filterValue string) error {
	publishInput := sns.PublishInput{TopicArn: aws.String(topicArn), Message: aws.String(message)}
	if groupId != "" {
		publishInput.MessageGroupId = aws.String(groupId)
	}
	if dedupId != "" {
		publishInput.MessageDeduplicationId = aws.String(dedupId)
	}
	if filterKey != "" && filterValue != "" {

		publishInput.MessageAttributes = map[string]types.MessageAttributeValue{
			filterKey: {DataType: aws.String("String"), StringValue: aws.String(filterValue)},
		}
	}
	_, err := actor.SnsClient.Publish(ctx, &publishInput)
	if err != nil {
		log.Printf("Couldn't publish message to topic %v. Here's why: %v", topicArn, err)
	}
	return err
}
