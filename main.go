package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/devjoaoGustavo/grepenv/filter"
)

var (
	config        = &aws.Config{Region: aws.String("us-east-1")}
	filterKey     *string
	filterVariant *string
	rawTerms      *string
	maxResult     *int64
	svc           *ssm.SSM
)

func init() {
	filterKey = flag.String("k", "Name", "Key by which the filter will be applied valid values: Name | Tier | Path | Type | KeyID | Tag key | Data type")
	filterVariant = flag.String("v", "Contains", "How filter will be applied. Valid values: Equals | BeginsWith | Contains")
	rawTerms = flag.String("t", "", "Search terms separated by comma")
	maxResult = flag.Int64("limit", 50, "Maximum number of results")
	flag.Parse()
	awsSession, err := session.NewSession(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ssm-explorer: %q", err)
	}
	svc = ssm.New(awsSession)
}

func main() {
	var paramNames []string
	param := make(chan string)
	resps := make(chan int)
	filterValues := filter.SplitValues(rawTerms)
	if len(filterValues) < 1 {
		fmt.Fprintln(os.Stderr, "\x1b[31mgrepenv: No term given\x1b[0m")
		os.Exit(1)
	}
	for _, value := range filterValues {
		go performDescribeParams(value, param, resps)
	}
	responses := 0
	for range filterValues {
		select {
		case n := <-resps:
			responses += n
		case <-time.After(5 * time.Second):
			continue
		}
	}
	for i := 0; i < responses; i++ {
		select {
		case p := <-param:
			paramNames = append(paramNames, p)
		case <-time.After(5 * time.Second):
			continue
		}
	}
	if len(paramNames) < 1 {
		fmt.Println("grepenv: No parameter found on SSM")
		os.Exit(1)
	}
	var refParamNames []*string
	for i := range paramNames {
		refParamNames = append(refParamNames, &paramNames[i])
	}
	searchParameters(refParamNames)
}

func performDescribeParams(value string, param chan<- string, resps chan<- int) {
	values := []*string{&value}
	describeParametersInput := buildInput(
		filterKey,
		filterVariant,
		values,
	)
	result, err := svc.DescribeParameters(&describeParametersInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mgrepenv: %q\x1b[0m", err)
		return
	}
	resps <- len(result.Parameters)
	for _, parameter := range result.Parameters {
		param <- *parameter.Name
	}
}

func searchParameters(paramNames []*string) {
	withDecryption := true
	for _, terms := range group(paramNames) {
		out, err := svc.GetParameters(&ssm.GetParametersInput{
			Names:          terms,
			WithDecryption: &withDecryption,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "ssm-explorer: %q", err)
		}
		for _, result := range out.Parameters {
			fmt.Printf("\x1b[32m%s\x1b[0m = %s\n", *result.Name, *result.Value)
		}
	}
}

func buildInput(key, variant *string, values []*string) ssm.DescribeParametersInput {
	filter := ssm.ParameterStringFilter{Key: key, Option: variant, Values: values}
	filters := []*ssm.ParameterStringFilter{&filter}
	return ssm.DescribeParametersInput{
		MaxResults:       maxResult,
		ParameterFilters: filters,
	}
}

func group(all []*string) (grouped [][]*string) {
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
