package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sethvargo/go-envconfig"
	"io/ioutil"
	"log"
	"net/http"
)

type Environment struct {
	Application   string `env:"APPLICATION,required"`
	Environment   string `env:"ENVIRONMENT,required"`
	Configuration string `env:"CONFIGURATION,required"`
}

func (c Environment) ConfigURI() string {
	return fmt.Sprintf("http://localhost:2772/applications/%s/environments/%s/configurations/%s", c.Application, c.Environment, c.Configuration)
}

func getConfiguration(ctx context.Context, uri string) (string, error) {
	log.Printf("Getting config from %s", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//Convert the body to type string
	return string(body), nil
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func withEnv(c Environment) func(ctx context.Context) (Response, error) {
	return func(ctx context.Context) (Response, error) {

		conf, err := getConfiguration(ctx, c.ConfigURI())
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Printf("Configuration : %s\n", conf)

		var buf bytes.Buffer

		body, err := json.Marshal(map[string]interface{}{
			"config": conf,
		})
		if err != nil {
			return Response{StatusCode: 404}, err
		}
		json.HTMLEscape(&buf, body)

		r := Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Body:            buf.String(),
			Headers: map[string]string{
				"Content-Type":           "application/json",
				"X-MyCompany-Func-Reply": "hello-handler",
			},
		}

		return r, nil
	}
}

func main() {

	ctx := context.Background()

	var c Environment
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	log.Printf("Configuration : %+v", c)

	lambda.Start(withEnv(c))
}
