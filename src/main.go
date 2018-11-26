package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	env "github.com/czerasz/go-lambda-sns-to-es/src/env"
	esIndexName "github.com/czerasz/go-lambda-sns-to-es/src/es_index_name"
	elastic "github.com/olivere/elastic"
)

var debug bool = os.Getenv("DEBUG") == "true"

// This function is executed by lambda
func handler(event events.SNSEvent) error {
	if debug {
		// Stringify and output event in debug mode
		str, err := json.Marshal(event)
		if err != nil {
			return err
		}
		fmt.Printf("%s", string(str))
	}

	ctx := context.Background()

	// Prepare placeholder for the SNS message
	snsMsg := make(map[string]interface{})

	// Make sure the SNS event contains actual records
	if len(event.Records) > 0 {
		snsMsgStr := event.Records[0].SNS.Message

		if err := json.Unmarshal([]byte(snsMsgStr), &snsMsg); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Error: no events were send\n%v", event)
	}

	indexName, err := esIndexName.Generate()
	if err != nil {
		return err
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(env.GetEnv("ES_URL", "http://127.0.0.1:9200")),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		return err
	}

	// Add a document to the index
	_, err = client.Index().
		Index(indexName).
		Type(env.GetEnv("ES_DOC_TYPE_NAME", "notification")).
		BodyJson(snsMsg).
		Refresh("wait_for").
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Used for local testing
	if env.GetEnv("LOCAL_TEST", "false") == "true" {
		// Get JSON object from the pipe
		var e events.SNSEvent
		if err := json.Unmarshal(readPipe(), &e); err != nil {
			panic(err)
		}

		// Execute serverless logic
		err := handler(e)

		if err != nil {
			panic(err)
		}
	} else {
		lambda.Start(handler)
	}
}

// Read data send via the pipe
func readPipe() []byte {
	reader := bufio.NewReader(os.Stdin)
	var output []byte

	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return output
}
