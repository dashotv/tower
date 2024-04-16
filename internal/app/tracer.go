package app

// TODO: opentelemetry

// func init() {
// 	initializers = append(initializers, setupTracer)
// 	starters = append(starters, startTracer)
// }
//
// func setupTracer(a *Application) error {
// 	return nil
// }
//
// func startTracer(ctx context.Context, a *Application) error {
// 	_ = initTracer() // TODO: deal with shutdown
// 	a.Router.Use(otelecho.Middleware("tower"))
// 	return nil
// }
//
// // setupOTelSDK bootstraps the OpenTelemetry pipeline.
// // If it does not return an error, make sure to call shutdown for proper cleanup.
// func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
// 	var shutdownFuncs []func(context.Context) error
//
// 	// shutdown calls cleanup functions registered via shutdownFuncs.
// 	// The errors from the calls are joined.
// 	// Each registered cleanup will be invoked once.
// 	shutdown = func(ctx context.Context) error {
// 		var err error
// 		for _, fn := range shutdownFuncs {
// 			err = errors.Join(err, fn(ctx))
// 		}
// 		shutdownFuncs = nil
// 		return err
// 	}
//
// 	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
// 	handleErr := func(inErr error) {
// 		err = errors.Join(inErr, shutdown(ctx))
// 	}
//
// 	// Set up propagator.
// 	prop := newPropagator()
// 	otel.SetTextMapPropagator(prop)
//
// 	// Set up trace provider.
// 	tracerProvider, err := newTraceProvider()
// 	if err != nil {
// 		handleErr(err)
// 		return
// 	}
// 	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
// 	otel.SetTracerProvider(tracerProvider)
//
// 	// Set up meter provider.
// 	meterProvider, err := newMeterProvider()
// 	if err != nil {
// 		handleErr(err)
// 		return
// 	}
// 	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
// 	otel.SetMeterProvider(meterProvider)
//
// 	return
// }
//
// func newPropagator() propagation.TextMapPropagator {
// 	return propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	)
// }
//
// func newTraceProvider() (*trace.TracerProvider, error) {
// 	traceExporter, err := stdouttrace.New(
// 		stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	traceProvider := trace.NewTracerProvider(
// 		trace.WithBatcher(traceExporter,
// 			// Default is 5s. Set to 1s for demonstrative purposes.
// 			trace.WithBatchTimeout(time.Second)),
// 	)
// 	return traceProvider, nil
// }
// func newSignozTraceProvider() (*trace.TracerProvider, error) {
// 	traceExporter, err := stdouttrace.New(
// 		stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	traceProvider := trace.NewTracerProvider(
// 		trace.WithBatcher(traceExporter,
// 			// Default is 5s. Set to 1s for demonstrative purposes.
// 			trace.WithBatchTimeout(time.Second)),
// 	)
// 	return traceProvider, nil
// }
// func newMeterProvider() (*metric.MeterProvider, error) {
// 	metricExporter, err := stdoutmetric.New()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	meterProvider := metric.NewMeterProvider(
// 		metric.WithReader(metric.NewPeriodicReader(metricExporter,
// 			// Default is 1m. Set to 3s for demonstrative purposes.
// 			metric.WithInterval(3*time.Second))),
// 	)
// 	return meterProvider, nil
// }
//
// func initTracer() func(context.Context) error {
// 	// insecure := os.Getenv("OTEL_INSECURE_MODE")
// 	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
// 	serviceName := os.Getenv("OTEL_SERVICE_NAME")
//
// 	// var secureOption otlptracegrpc.Option
// 	// if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
// 	// 	secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
// 	// } else {
// 	// 	secureOption = otlptracegrpc.WithInsecure()
// 	// }
//
// 	// exporter, err := otlptracegrpc.New(
// 	// 	context.Background(),
// 	// 	otlptracegrpc.WithInsecure(),
// 	// 	otlptracegrpc.WithEndpoint(collectorURL),
// 	// )
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create exporter: %v", err)
// 	// }
//
// 	// HTTP
// 	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(collectorURL))
// 	if err != nil {
// 		log.Fatalf("Failed to create exporter: %v", err)
// 	}
//
// 	resources, err := resource.New(
// 		context.Background(),
// 		resource.WithAttributes(
// 			attribute.String("service.name", serviceName),
// 			attribute.String("library.language", "go"),
// 		),
// 	)
// 	if err != nil {
// 		log.Fatalf("Could not set resources: %v", err)
// 	}
//
// 	otel.SetTracerProvider(
// 		trace.NewTracerProvider(
// 			trace.WithSampler(trace.AlwaysSample()),
// 			trace.WithBatcher(exporter),
// 			trace.WithResource(resources),
// 		),
// 	)
// 	return exporter.Shutdown
// }
