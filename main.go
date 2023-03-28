package main

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

func main() {
	provider, _ := otelcol.NewConfigProvider(newDefaultConfigProviderSettings([]string{
		"file:logging.yaml",
		"file:otlp-grpc.yaml",
		"file:otlp-receiver.yaml",
		"file:pipeline.yaml",
		"file:otlp-exporter.yaml",
		// "file:pipeline-with-exporter.yaml",
	}))

	if forwarder {
		"inmem:otlp-forwarder"
	}

	factories, _ := components()

	cfg, err := provider.Get(context.Background(), factories)
	if err != nil {
		panic(err)
	}
	spew.Dump(cfg)
}

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Receivers, err = receiver.MakeFactoryMap(
		otlpreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Processors, err = processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Extensions, err = extension.MakeFactoryMap()
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Exporters, err = exporter.MakeFactoryMap(
		otlphttpexporter.NewFactory(),
		loggingexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}

func makeMapProvidersMap(providers ...confmap.Provider) map[string]confmap.Provider {
	ret := make(map[string]confmap.Provider, len(providers))
	for _, provider := range providers {
		ret[provider.Scheme()] = provider
	}
	return ret
}

func newDefaultConfigProviderSettings(uris []string) otelcol.ConfigProviderSettings {
	return otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       uris,
			Providers:  makeMapProvidersMap(fileprovider.New()),
			Converters: []confmap.Converter{New()},
		},
	}
}

type dumpconverter struct{}

func (dumpconverter) Convert(_ context.Context, conf *confmap.Conf) error {
	for _, k := range conf.AllKeys() {
		fmt.Println(k)
		spew.Dump(conf.Get(k))
	}
	return nil
}

func New() *dumpconverter {
	return &dumpconverter{}
}
