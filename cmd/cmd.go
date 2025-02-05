package cmd

import (
	"context"

	"github.com/spf13/cobra"
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/gigapipehq/loggen/cmd/run"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/otel"
	"github.com/gigapipehq/loggen/internal/prom"
	_default "github.com/gigapipehq/loggen/internal/senders/default"
)

var tp *trace.TracerProvider
var rootCMD = &cobra.Command{
	Use:   "loggen",
	Short: "A fake log, metric and trace generator for qryn Cloud",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		sender := _default.New().WithHeaders(map[string]string{
			"X-API-Key":    cfg.APIKey,
			"X-API-Secret": cfg.APISecret,
		})
		prom.Initialize(context.Background(), sender.Client(), cfg)
		exporter := otel.NewExporter(cfg.URL, sender.Client())
		tp = otel.NewProvider(exporter, cfg)
		otelsdk.SetTracerProvider(tp)
	},
}

func init() {
	config.Load()
	rootCMD.AddCommand(run.CMD())
}

func Execute() error {
	defer func() {
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
	}()
	return rootCMD.Execute()
}
