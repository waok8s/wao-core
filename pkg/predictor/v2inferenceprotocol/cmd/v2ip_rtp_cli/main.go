package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/waok8s/wao-core/pkg/predictor/v2inferenceprotocol"
	"github.com/waok8s/wao-core/pkg/util"
)

func main() {
	var urlTmpl string
	flag.StringVar(&urlTmpl, "urlTemplate", "http://localhost:5000/v2/models/myModel:{{.App}}/versions/v0.1.0/infer", "Prediction service URL template")
	var appName string
	flag.StringVar(&appName, "appName", "app", "App name")
	var cpuUsage float64
	flag.Float64Var(&cpuUsage, "cpuUsage", 0.0, "CPU usage")
	var basicAuth string
	flag.StringVar(&basicAuth, "basicAuth", "", "Basic auth in username@password format")
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout for the request")
	var logLevel int
	flag.IntVar(&logLevel, "v", 3, "klog-style log level")
	flag.Parse()

	var slogLevel slog.Level
	switch {
	case logLevel < 0:
		slogLevel = 100 // silent
	case logLevel == 0:
		slogLevel = slog.LevelError
	case logLevel == 1:
		slogLevel = slog.LevelWarn
	case logLevel == 2:
		slogLevel = slog.LevelInfo
	case logLevel == 3:
		slogLevel = slog.LevelDebug
	case logLevel > 3:
		slogLevel = -100 // verbose
	}

	lg := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevel,
	}))
	slog.SetDefault(lg.With("component", "ResponseTimePredictor (V2InferenceProtocol)"))

	requestEditorFns := []util.RequestEditorFn{}
	ss := strings.Split(basicAuth, ":")
	if len(ss) == 2 {
		requestEditorFns = append(requestEditorFns, util.WithBasicAuth(ss[0], ss[1]))
	}
	requestEditorFns = append(requestEditorFns, util.WithCurlLogger(lg.With("func", "WithCurlLogger(v2inferenceprotocol.ResponseTimePredictor.Predict)")))

	c := v2inferenceprotocol.NewResponseTimePredictor(urlTmpl, true, timeout, requestEditorFns...)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	v, err := c.Predict(ctx, appName, cpuUsage)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
}
