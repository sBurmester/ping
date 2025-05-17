package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sburmester/ping/pkg/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	_ "github.com/KimMachineGun/automemlimit" // for side effects
	_ "go.uber.org/automaxprocs"              // for side effects
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	slog.Info("starting ping server")

	shutdown, err := installTelemetry()

	if err != nil {
		return fmt.Errorf("cannot install telemetry: %w", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := shutdown(ctx); err != nil {
			slog.Default().Error("cannot shutdown telemetry", err.Error())
		}
	}()

	err = http.StartServer()
	if err != nil {
		slog.Default().Error("error starting server", err)
	}
	return err
}

func installTelemetry() (func(context.Context) error, error) {
	r, err := resource.New(
		context.Background(),

		resource.WithTelemetrySDK(),

		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceName("ping-server"),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("cannot create resource: %w", err)
	}

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	metricReader, err := prometheus.New()

	if err != nil {
		return nil, fmt.Errorf("cannot create metric reader: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metricReader),
	)

	otel.SetMeterProvider(meterProvider)

	shutdown := func(ctx context.Context) error {
		return errors.Join(
			meterProvider.ForceFlush(ctx),
			meterProvider.Shutdown(ctx),
		)
	}

	return shutdown, nil
}
