# CLAUDE.md

<!-- Canonical source: AGENTS.md. This file is auto-generated for Claude Code compatibility. -->

This file provides guidance to AI coding assistants when working with this repository.

## Project Overview

OCM CSV Parser — a CLI tool to parse and format CSV data for consumption by OCM services. Reads CSV files and transforms them into structured output suitable for OCM API ingestion.

## Build & Test Commands

```bash
go build ./...       # Build the project
go test ./...        # Run all tests
```

## Architecture

- **cmd/**: CLI entry point
- **pkg/**: Core parsing and formatting logic

## Key Conventions

- Module path: `github.com/openshift-online/ocm-csv-parser`
- Standard Go CLI application structure
