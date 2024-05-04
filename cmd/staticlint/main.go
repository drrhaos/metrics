package main

import (
	"metrics/internal/analyzers/osexit"
	"strings"

	"github.com/mdempsky/maligned/passes/maligned"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	checks := map[string]bool{
		"ST1000": true,
		"ST1003": true,
	}

	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Analyzer.Name, "SA") {
			checks[v.Analyzer.Name] = true
		}
	}

	mychecks := []*analysis.Analyzer{
		structtag.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		appends.Analyzer,
		usesgenerics.Analyzer,
		maligned.Analyzer,
		bodyclose.Analyzer,
		osexit.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
