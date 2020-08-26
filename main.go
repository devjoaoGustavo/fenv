package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/devjoaoGustavo/fenv/service"
)

func main() {
	filterKey := flag.String("k", "Name", "Key by which the filter will be applied valid values: Name | Tier | Path | Type | KeyID | Tag key | Data type")
	filterVariant := flag.String("v", "Contains", "How filter will be applied. Valid values: Equals | BeginsWith | Contains")
	rawTerms := flag.String("t", "", "Search terms separated by comma")
	maxResult := flag.Int64("limit", 50, "Maximum number of results")
	flag.Parse()
	svc := service.New(rawTerms, filterKey, filterVariant, maxResult)
	if len(svc.FilterValues) < 1 {
		fmt.Println("fenv: No term given")
		os.Exit(1)
	}
	svc.DescribeParameters()
	svc.GetParameters()
}
