package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"geliumui/internal/gates"
)

var shell = regexp.MustCompile(`\.[\w-]*(?:topbar|appbar|navbar|site-nav)[\w-]*\s*\{`)
var hex = regexp.MustCompile(`#[0-9a-fA-F]{3,8}(?:$|[^A-Za-z0-9_-])`)
var form = regexp.MustCompile(`(?is)<form\b[^>]*>.*?</form>`)
var validationHook = regexp.MustCompile(`(?i)validation-summary|X-Gelium-Validation|gelium-validation`)
var darkMedia = regexp.MustCompile(`(?i)prefers-color-scheme\s*:\s*dark`)
var image = regexp.MustCompile(`(?is)<img\b[^>]*>`)
var altAttribute = regexp.MustCompile(`(?is)\balt\s*=\s*(?:"([^"]*)"|'([^']*)')`)
var widthAttribute = regexp.MustCompile(`(?i)\bwidth\s*=`)
var heightAttribute = regexp.MustCompile(`(?i)\bheight\s*=`)

type paths []string

func (p *paths) String() string     { return "" }
func (p *paths) Set(v string) error { *p = append(*p, v); return nil }
func main()                         { os.Exit(run(os.Args[1:], os.Stdout)) }
func run(args []string, out io.Writer) int {
	f := flag.NewFlagSet("gelium-ux-detect", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	var owned, shared paths
	format := f.String("format", "text", "")
	exceptionPath := f.String("exceptions", "", "")
	f.Var(&owned, "owned", "")
	f.Var(&shared, "shared", "")
	if err := f.Parse(args); err != nil {
		writeDetectorResult(out, *format, gates.DetectorResult{Status: "invalid-configuration", Findings: []gates.DetectorFinding{}})
		return 2
	}
	if (*format != "text" && *format != "json") || len(owned) == 0 || len(f.Args()) == 0 || invalidScope(owned, shared, f.Args()) {
		writeDetectorResult(out, *format, gates.DetectorResult{Status: "invalid-configuration", Findings: []gates.DetectorFinding{}})
		return 2
	}
	findings := []gates.DetectorFinding{}
	for _, path := range f.Args() {
		b, err := os.ReadFile(path)
		if err != nil {
			writeDetectorResult(out, *format, gates.DetectorResult{Status: "invalid-configuration", Findings: []gates.DetectorFinding{}})
			return 2
		}
		content := string(b)
		if customShell(content) {
			findings = append(findings, gates.DetectorFinding{ID: "custom-shell", Rule: "shell-contract", Path: path})
		}
		if hasOneOffHex(content) {
			findings = append(findings, gates.DetectorFinding{ID: "color-literal", Rule: "token-contract", Path: path})
		}
		if filepath.Ext(path) == ".css" && darkMedia.MatchString(content) {
			findings = append(findings, gates.DetectorFinding{ID: "media-dark", Rule: "theme-contract", Path: path})
		}
		if hasFormWithoutValidationHook(content) {
			findings = append(findings, gates.DetectorFinding{ID: "form-validation", Rule: "form-contract", Path: path})
		}
		if hasUnknownInformativeMedia(content) {
			findings = append(findings, gates.DetectorFinding{ID: "media-metadata-unknown", Rule: "media-contract", Path: path})
		}
	}
	exceptions := []gates.DetectorException{}
	if *exceptionPath != "" {
		body, err := os.ReadFile(*exceptionPath)
		if err != nil || json.Unmarshal(body, &exceptions) != nil {
			writeDetectorResult(out, *format, gates.DetectorResult{Status: "invalid-configuration", Findings: []gates.DetectorFinding{}})
			return 2
		}
	}
	result := gates.ClassifyDetectorFindings(findings, gates.Scope{OwnedPaths: owned, SharedPaths: shared}, exceptions, now())
	if result.Status == "invalid-configuration" {
		result.Findings = []gates.DetectorFinding{}
	}
	writeDetectorResult(out, *format, result)
	if result.Status == "invalid-configuration" {
		return 2
	}
	if result.Status == "clean-pass" || result.Status == "pass-with-exceptions" {
		return 0
	}
	return 1
}
func now() (t time.Time) { return time.Now().UTC() }

func writeDetectorResult(out io.Writer, format string, result gates.DetectorResult) {
	if format == "json" {
		_ = json.NewEncoder(out).Encode(result)
		return
	}
	fmt.Fprintf(out, "gelium-ux-detect: %s\n", result.Status)
	for _, x := range result.Findings {
		fmt.Fprintf(out, "%s %s %s", x.Rule, x.Attribution, x.Path)
		if x.ExceptionID != "" {
			fmt.Fprintf(out, " exception=%s", x.ExceptionID)
		}
		if x.ExpiredExceptionID != "" {
			fmt.Fprintf(out, " expired-exception=%s", x.ExpiredExceptionID)
		}
		fmt.Fprintln(out)
	}
}

func invalidScope(owned, shared, scanned []string) bool {
	declared := make(map[string]string, len(owned)+len(shared))
	for _, path := range owned {
		if prior, exists := declared[path]; exists && prior != "owned" {
			return true
		}
		declared[path] = "owned"
	}
	for _, path := range shared {
		if prior, exists := declared[path]; exists && prior != "shared" {
			return true
		}
		declared[path] = "shared"
	}
	for _, path := range scanned {
		if _, exists := declared[path]; !exists {
			return true
		}
	}
	return false
}

func customShell(content string) bool {
	for _, selector := range shell.FindAllString(content, -1) {
		if !strings.HasPrefix(strings.TrimPrefix(strings.TrimSpace(selector), "."), "ui-") {
			return true
		}
	}
	return false
}

func hasOneOffHex(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "--ui-") && strings.Contains(line, ":") {
			continue
		}
		if hex.MatchString(line) {
			return true
		}
	}
	return false
}

func hasFormWithoutValidationHook(content string) bool {
	for _, candidate := range form.FindAllString(content, -1) {
		if !validationHook.MatchString(candidate) {
			return true
		}
	}
	return false
}

func hasUnknownInformativeMedia(content string) bool {
	for _, tag := range image.FindAllString(content, -1) {
		alt := altAttribute.FindStringSubmatch(tag)
		if len(alt) != 3 || strings.TrimSpace(alt[1]+alt[2]) == "" {
			continue
		}
		if !widthAttribute.MatchString(tag) || !heightAttribute.MatchString(tag) {
			return true
		}
	}
	return false
}
