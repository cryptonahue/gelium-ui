// Command gelium-preflight validates a versioned Gelium gate ledger.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"geliumui/internal/gates"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, time.Now().UTC()))
}

func run(args []string, output io.Writer, now time.Time) int {
	if len(args) == 0 || (args[0] != "prebuild" && args[0] != "release") {
		fmt.Fprintln(output, "usage: gelium-preflight <prebuild|release> --ledger <path> [--changed <path>] [--format text|json]")
		return 2
	}
	mode := args[0]
	flags := flag.NewFlagSet(mode, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	ledgerPath := flags.String("ledger", "", "path to JSON ledger")
	format := flags.String("format", "text", "text or json")
	var changed paths
	flags.Var(&changed, "changed", "changed UI-relevant path (repeatable)")
	if err := flags.Parse(args[1:]); err != nil || *ledgerPath == "" || (*format != "text" && *format != "json") {
		fmt.Fprintln(output, "usage: gelium-preflight <prebuild|release> --ledger <path> [--changed <path>] [--format text|json]")
		return 2
	}
	data, err := os.ReadFile(*ledgerPath)
	if err != nil {
		return emit(output, mode, *format, gates.PreflightResult{Status: "invalid-configuration", Issues: []gates.Issue{{Code: "ledger-read-failed", Field: "ledger", Message: err.Error()}}})
	}
	ledger, issues := gates.ValidateLedger(data, now)
	if mode == "release" {
		return emit(output, mode, *format, gates.EvaluateRelease(ledger, issues, changed))
	}
	return emit(output, mode, *format, gates.EvaluatePrebuild(ledger, issues, changed))
}

type paths []string

func (p *paths) String() string { return "" }
func (p *paths) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func emit(output io.Writer, mode, format string, result gates.PreflightResult) int {
	if format == "json" {
		_ = json.NewEncoder(output).Encode(result)
	} else {
		fmt.Fprintf(output, "gelium-preflight %s: %s\n", mode, result.Status)
		for _, issue := range result.Issues {
			fmt.Fprintf(output, "- %s: %s\n", issue.Code, issue.Message)
		}
	}
	if result.Status == "pass" {
		return 0
	}
	if result.Status == "invalid-configuration" {
		return 2
	}
	return 1
}
