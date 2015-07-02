package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Config Struct for unmarshalling yaml config file
type Config struct {
	Queue_url          *string
	Message_attributes []*string
	Region             *string
}

// Get Message Function
// Takes SQS Object, Queue_url, and Message Attributes
// Returns a ReceiveMessageOutput Struct and Errors
func getMessage(queue *sqs.SQS, queue_url *string, message_attributes []*string) (sqs.ReceiveMessageOutput, error) {

	// Params object of ReceiveMessageInput Struct
	params := &sqs.ReceiveMessageInput{
		QueueURL:              queue_url,
		MessageAttributeNames: message_attributes,
		MaxNumberOfMessages:   aws.Long(1),
	}
	resp, err := queue.ReceiveMessage(params)
	return *resp, err
}

// Load Config Function
// Takes path as string of config file,
// Unmarshalls the yaml file into the Config type and returns it and any errors
func LoadConfig(path string) (Config, error) {
	c := Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error Occurred: %v", err)
		return c, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	return c, err
}

// Type Method to return Message Attributes
func (c Config) getAttributes() []*string {
	return c.Message_attributes
}

// Type Method to return queue_url
func (c Config) getUrl() *string {
	return c.Queue_url
}

// Type Method to return AWS Region
func (c Config) getRegion() *string {
	return c.Region
}

// catchError Function to correctly parse any AWS errors returned from go-aws-sdk
func catchError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		// Generic AWS Error with Code, Message, and original error (if any)
		log.Fatal(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			// A service error occurred
			log.Fatal(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
		}
	} else {
		log.Fatal(err.Error())
	}
}

func main() {
	// Load Config Path
	config, err := LoadConfig("config.yml")

	// If error, panic
	if err != nil {
		panic(err)
	}

	queue_url := config.getUrl()
	message_attr := config.getAttributes()
	// Need to dereference the region pointer to pass to &aws.Config
	region := *config.getRegion()

	// Create new sqs queue Object with Config Supplied
	queue := sqs.New(&aws.Config{Region: region})

	// Retrieve a message from the queue
	message, err := getMessage(queue, queue_url, message_attr)

	if err != nil {
		catchError(err)
	}

	// Print the Message
	// Would love to abstract this awsutil call, but
	// sqs.ReceiveMessageOutput.String() cannot be found, even though it's here:
	// https://github.com/aws/aws-sdk-go/blob/master/service/sqs/api.go#L1639
	fmt.Println(awsutil.StringValue(message))
}
