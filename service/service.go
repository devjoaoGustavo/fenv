package service

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/devjoaoGustavo/fenv/finder"
)

var (
	results = make(chan string)
	qtd     = make(chan int)
)

// SearchParameters searches for parameters that contain the given terms
func SearchParameters(terms []string) ([]string, error) {
	f, err := finder.NewFinder()
	if err != nil {
		return nil, err
	}
	params := terms[:0]
	for i := range terms {
		input := buildInput(f, []*string{&terms[i]})
		go f.SearchParams(&input, qtd, results)
	}
	params = handleChannels(qtd, results, len(terms))
	if len(params) == 0 {
		return nil, fmt.Errorf("describe parameter: not found")
	}
	return params, nil
}

// GetParameters searches for the values of the given param names
func GetParameters(params []string) ([]string, error) {
	f, err := finder.NewFinder()
	if err != nil {
		return nil, err
	}
	values := params[:0]
	grouped := group(params)
	for _, names := range grouped {
		go f.FindParams(names, qtd, results)
	}
	values = handleChannels(qtd, results, len(grouped))
	if len(values) == 0 {
		return nil, fmt.Errorf("get parameter: no parameter found")
	}
	return values, nil
}

// Auxiliar functions that I think could be moved to another package

func buildInput(f finder.Finder, values []*string) ssm.DescribeParametersInput {
	filter := ssm.ParameterStringFilter{Key: &f.Key, Option: &f.Variant, Values: values}
	filters := []*ssm.ParameterStringFilter{&filter}
	return ssm.DescribeParametersInput{
		MaxResults:       &f.MaxResult,
		ParameterFilters: filters,
	}
}

func handleChannels(qtd chan int, results chan string, limit int) []string {
	var output []string
	var n int
	for i := 0; i < limit; i++ {
		select {
		case x := <-qtd:
			n += x
		case <-time.After(time.Second << 5):
			continue
		}
	}
	for i := 0; i < n; i++ {
		select {
		case result := <-results:
			output = append(output, result)
		case <-time.After(time.Second << 5):
			continue
		}
	}
	return output
}

func group(params []string) (grouped [][]*string) {
	all := toRefs(params)
	maxTermsPerRequest := 10
	slicesCount := len(all) / maxTermsPerRequest
	start := 0
	end := maxTermsPerRequest
	for i := 0; i < slicesCount; i++ {
		grouped = append(grouped, all[start:end])
		start = end
		end += maxTermsPerRequest
	}
	if slicesCount*maxTermsPerRequest < len(params) {
		grouped = append(grouped, all[slicesCount*maxTermsPerRequest:])
	}
	return
}

func toRefs(input []string) (result []*string) {
	for i := range input {
		result = append(result, &input[i])
	}
	return
}
