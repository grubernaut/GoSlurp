package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Config Struct for unmarshalling yaml config file
type Config struct {
	Queue_url          *string
	Message_attributes []*string
	Region             *string
	Outputs            []string
}

// Get Message Function
// Takes SQS Object, Queue_url, and Message Attributes
// Returns a ReceiveMessageOutput Struct and Errors
func getMessage(queue *sqs.SQS, queue_url *string, message_attributes []*string) (sqs.Message, error) {

	// Params object of ReceiveMessageInput Struct
	params := &sqs.ReceiveMessageInput{
		QueueURL:              queue_url,
		MessageAttributeNames: message_attributes,
		MaxNumberOfMessages:   aws.Long(1),
		VisibilityTimeout:     aws.Long(1),
		WaitTimeSeconds:       aws.Long(1),
	}
	resp, err := queue.ReceiveMessage(params)

	//	fmt.Println(reflect.ValueOf(resp.Messages).Kind())
	message := *resp.Messages[0]
	return message, err
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

// Type Method to return Specified Outputs
func (c Config) getOutputs() []string {
	return c.Outputs
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

// Function to Check if slice of strings contains a string
func contains(s []string, c string) bool {
	for _, a := range s {
		if a == c {
			return true
		}
	}
	return false
}

// Print Outputs Function to only print specified outputs
func printOutputs(message sqs.Message, outputs []string) {

	valid_outputs := validOutputs(message)

	for _, output := range outputs {
		// if supplied output is valid, print it
		if contains(valid_outputs, output) {
			//				fmt.Println(*message.Body)

			//fmt.Println(f)
			fmt.Println(reflect.ValueOf(message).FieldByName("Body"))
		}
	}
}

// Function that takes the received message and returns a slice
// of field names from returned Messages Struct via power of reflection.
func validOutputs(message sqs.Message) []string {
	valid_outputs := make([]string, 0)
	s := reflect.ValueOf(message)
	typeOfMessage := s.Type()
	for i := 0; i < s.NumField(); i++ {
		valid_outputs = append(valid_outputs, typeOfMessage.Field(i).Name)
	}
	return valid_outputs
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
	outputs := config.getOutputs()

	// Create new sqs queue Object with Config Supplied
	queue := sqs.New(&aws.Config{Region: region})

	// Retrieve a message from the queue
	message, err := getMessage(queue, queue_url, message_attr)

	printOutputs(message, outputs)

	if err != nil {
		catchError(err)
	}

}
