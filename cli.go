package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

func printVersion() {
	fmt.Printf("volley version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		webhookUrl string
		severity   string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet("volley", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&webhookUrl, "webhook-url", "", "")
	flags.StringVar(&webhookUrl, "w", "", "(Short)")

	flags.StringVar(&severity, "severity", "", "")
	flags.StringVar(&severity, "s", "", "(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		printVersion()
		return ExitCodeOK
	}

	if webhookUrl == "" {
		logrus.Fatal("webhook-url is required")
	}
	if severity == "" {
		severity = "critical"
	}

	err := cli.send(os.Stdin, webhookUrl, severity)
	if err != nil {
		logrus.Error(os.Stderr, err)
		os.Exit(ExitCodeError)
	}

	return ExitCodeOK
}

type CommonLabels struct {
	Alertname   string `json:"alertname"`
	Node        string `json:"node"`
	ServiceName string `json:"service_name"`
	Severity    string `json:"severity"`
}
type CommonAnnotations struct {
	Summary string `json:"summary"`
	Text    string `json:"text"`
}

type PrometheusFormat struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname   string `json:"alertname"`
			Node        string `json:"node"`
			ServiceName string `json:"service_name"`
			Severity    string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			FiringText   string `json:"firing_text"`
			ResolvedText string `json:"resolved_text"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
		Severity  string `json:"severity"`
	} `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	ID                string            `json:"id"`
	Topic             struct {
	} `json:"topic"`
}

func (cli *CLI) send(stdin io.Reader, url, severity string) error {
	body := ""
	summary := ""
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		body += scanner.Text()
		if summary == "" {
			summary = body
		}
	}

	t := PrometheusFormat{
		CommonLabels: CommonLabels{
			Alertname: "snmptrap",
			Severity:  severity,
		},
		CommonAnnotations: CommonAnnotations{
			Summary: summary,
			Text:    body,
		},
	}

	buf, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(buf)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
