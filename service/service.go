package service

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/devjoaoGustavo/fenv/filter"
)

var config = &aws.Config{Region: aws.String("us-east-1")}

type Service struct {
	FilterValues []string
	SVC          *ssm.SSM
	Key, Variant *string
	MaxResult    *int64
	Result
}

type Result struct {
	Responses  int
	ParamNames []string
}

func New(rawTerms, key, variant *string, maxResult *int64) Service {
	session, err := session.NewSession(config)
	if err != nil {
		panic(fmt.Sprintf("fenv: %q", err))
	}
	return Service{
		FilterValues: filter.SplitValues(rawTerms),
		SVC:          ssm.New(session),
		Key:          key,
		Variant:      variant,
		MaxResult:    maxResult,
	}
}

func (self *Service) DescribeParameters() {
	results := make(chan string)
	n := make(chan int)
	for _, value := range self.FilterValues {
		values := []*string{&value}
		describeParametersInput := self.buildInput(
			self.Key,
			self.Variant,
			values,
		)
		go self.performDescribeParameters(&describeParametersInput, n, results)
	}
	qtd := countValues(n, len(self.FilterValues))
	for i := 0; i < qtd; i++ {
		select {
		case p := <-results:
			self.ParamNames = append(self.ParamNames, p)
		case <-time.After(5 * time.Second):
			continue
		}
	}
}

func (self Service) buildInput(key, variant *string, values []*string) ssm.DescribeParametersInput {
	filter := ssm.ParameterStringFilter{Key: key, Option: variant, Values: values}
	filters := []*ssm.ParameterStringFilter{&filter}
	return ssm.DescribeParametersInput{
		MaxResults:       self.MaxResult,
		ParameterFilters: filters,
	}
}

func (self Service) performDescribeParameters(input *ssm.DescribeParametersInput, n chan<- int, results chan<- string) {
	result, err := self.SVC.DescribeParameters(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mfenv: %q\x1b[0m", err)
		return
	}
	n <- len(result.Parameters)
	for _, parameter := range result.Parameters {
		results <- *parameter.Name
	}
}

func (self Result) HasParams() bool {
	return len(self.ParamNames) > 0
}

func (self Service) GetParameters() {
	results := make(chan string)
	n := make(chan int)
	grouped := self.groupedParams()
	for _, terms := range grouped {
		go self.performGetParameters(terms, n, results)
	}
	qtd := countValues(n, len(grouped))
	for i := 0; i < qtd; i++ {
		select {
		case param := <-results:
			fmt.Print(param)
		case <-time.After(5 * time.Second):
			continue
		}
	}
}

func countValues(n chan int, limit int) (qtd int) {
	for i := 0; i < limit; i++ {
		select {
		case x := <-n:
			qtd += x
		case <-time.After(10 * time.Second):
			continue
		}
	}
	return
}

func (self Service) performGetParameters(terms []*string, n chan<- int, results chan<- string) {
	withDecryption := true
	output, err := self.SVC.GetParameters(&ssm.GetParametersInput{
		Names:          terms,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "fenv: %q", err)
	}
	n <- len(output.Parameters)
	for _, result := range output.Parameters {
		results <- fmt.Sprintf("\x1b[32m%s\x1b[0m = %s\n", *result.Name, *result.Value)
	}
}

func (self Result) groupedParams() (grouped [][]*string) {
	all := self.RefParamNames()
	maxTermsPerRequest := 10
	slicesCount := len(all) / maxTermsPerRequest
	start := 0
	end := maxTermsPerRequest
	for i := 0; i < slicesCount; i++ {
		grouped = append(grouped, all[start:end])
		start = end
		end += maxTermsPerRequest
	}
	if slicesCount*maxTermsPerRequest < len(all) {
		grouped = append(grouped, all[slicesCount*maxTermsPerRequest:])
	}
	return
}

func (self Result) RefParamNames() (result []*string) {
	for i := range self.ParamNames {
		result = append(result, &self.ParamNames[i])
	}
	return
}
