package newrelic

import (
	"context"
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

type Config struct {
	Enabled           bool
	AppName           string
	LicenseKey        string
	DistributedTracer bool
}

func Initialize(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}
	a, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(cfg.DistributedTracer),
		newrelic.ConfigEnabled(true),
	)
	if err != nil {
		return err
	}
	app = a
	return nil
}

func GetApplication() *newrelic.Application { return app }
func IsEnabled() bool                       { return app != nil }

func NoticeErrorString(ctx context.Context, msg string) {
	if app == nil {
		return
	}
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		txn.NoticeError(errors.New(msg))
	}
}
