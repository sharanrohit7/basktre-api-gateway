package appinit

import (
	"context"

	brconfig "github.com/basktre/api-gateway/pkg/config"
	"github.com/basktre/api-gateway/pkg/logger"
	brnewrelic "github.com/basktre/api-gateway/pkg/newrelic"
	"github.com/basktre/api-gateway/pkg/notifications/slack"
)

func Initialize(_ context.Context) {
	initializeLogger()
	initializeNewRelic()
	initializeSlack()
}

func initializeLogger() {
	level := brconfig.GetStringDefault("log.level", "info")
	format := brconfig.GetStringDefault("log.format", "json")
	if err := logger.Initialize(level, format); err != nil {
		panic("api-gateway: failed to initialize logger: " + err.Error())
	}
	logger.Infof("Logger initialized (level=%s, format=%s)", level, format)
}

func initializeNewRelic() {
	cfg := brnewrelic.Config{
		Enabled:           brconfig.GetBool("newrelic.enabled"),
		AppName:           brconfig.GetStringDefault("newrelic.appName", "api-gateway"),
		LicenseKey:        brconfig.GetString("newrelic.licenseKey"),
		DistributedTracer: brconfig.GetBool("newrelic.distributedTracer"),
	}
	if err := brnewrelic.Initialize(cfg); err != nil {
		logger.Warnf("NewRelic initialization failed (non-fatal): %v", err)
		return
	}
}

func initializeSlack() {
	webhookURL := brconfig.GetString("notification.slack.webhookURL")
	channelID := brconfig.GetString("notification.slack.errorChannelID")
	client := slack.NewClient(webhookURL, channelID)
	slack.Set(client)
}
