//This package provides a struct to deal with the env variables searching
package finder

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Finder struct {
	*ssm.SSM
	Key       string
	Variant   string
	MaxResult int64
}

// NewFinder provides a new instance of Finder struct
func NewFinder() (Finder, error) {
	session, err := startSession()
	if err != nil {
		return Finder{}, err
	}
	return Finder{session, "Name", "Contains", int64(50)}, nil
}

func startSession() (*ssm.SSM, error) {
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		awsRegion = "us-east-2"
	}
	config := &aws.Config{Region: aws.String(awsRegion)}
	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return ssm.New(session), nil
}

// SearchParams search for env vars on SSM service
func (f Finder) SearchParams(input *ssm.DescribeParametersInput, n chan<- int, results chan<- string) error {
	var data []*ssm.ParameterMetadata
	result, err := f.DescribeParameters(input)
	if err != nil {
		return err
	}
	data = result.Parameters
	for result.NextToken != nil {
		input.SetNextToken(*result.NextToken)
		result, _ = f.DescribeParameters(input)
		data = append(data, result.Parameters...)
	}
	n <- len(data)
	for _, parameter := range data {
		results <- *parameter.Name
	}
	return nil
}

// FindParams uses the result of Finder.SearchParams in order to get the variable values
func (f Finder) FindParams(terms []*string, n chan<- int, results chan<- string) error {
	withDecryption := true
	output, err := f.GetParameters(&ssm.GetParametersInput{
		Names:          terms,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return err
	}
	n <- len(output.Parameters)
	for _, result := range output.Parameters {
		results <- fmt.Sprintf("\x1b[32m%s\x1b[0m = %s", *result.Name, *result.Value)
	}
	return nil
}
