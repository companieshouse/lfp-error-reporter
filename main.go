package main

import (
	"fmt"

	. "github.com/aws/aws-lambda-go/lambda"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/lambda"
)

func main() {

	log.Namespace = "lfp-error-reporter"

	cfg, err := config.Get()
	if err != nil {
		log.Error(fmt.Errorf("error configuring service: %s. - exiting", err), nil)
		return
	}

	log.Trace("Config", log.Data{"Config": cfg})
	log.Info("Penalty payment error reporter lambda started")

	errorReporterLambda := lambda.New(cfg)

	Start(errorReporterLambda.Execute)
}
