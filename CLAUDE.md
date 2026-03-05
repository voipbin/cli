# VoIPBIN CLI (vn)

## Quick Reference

```
make build       # Build binary to bin/vn
make test        # Run all tests (go test ./...)
make lint        # Run linter (golangci-lint run)
make install     # Install to $GOPATH/bin (go install ./cmd/vn)
make clean       # Remove bin/ directory
```

## Project Structure

```
cmd/vn/main.go              Entry point
internal/
  config/config.go           Profile and config management (~/.vn/config.yaml)
  auth/auth.go               API client creation, access key resolution
  output/                    Output formatters (table, JSON, YAML)
    output.go                Formatter interface and helpers
    table.go                 Table formatter (tablewriter)
    json.go                  JSON formatter (pretty-printed)
    yaml.go                  YAML formatter (via JSON intermediate)
  commands/                  All CLI commands (one file per resource)
    root.go                  Root command, global flags, command registration
    login.go                 Login/logout with key validation
    calls.go, agents.go...   ~45 resource command files
Makefile                     Build automation
go.mod                       Module definition (github.com/voipbin/vn-cli)
```

## Architecture

- **CLI framework:** Cobra (`spf13/cobra`). Root command in `internal/commands/root.go`.
- **API client:** `voipbin-go` SDK, auto-generated from OpenAPI. Local module via `replace` directive in `go.mod` pointing to `../voipbin-go`.
- **Config:** YAML at `~/.vn/config.yaml`, file permissions 0600. Multi-profile support.
- **Auth:** Access key resolution priority: `--access-key` flag > `VN_ACCESS_KEY` env var > config profile. Key is passed as `accesskey` query parameter via custom `http.RoundTripper`.
- **Output:** Pluggable formatters behind `output.Formatter` interface. Table (default), JSON, YAML. YAML goes through JSON intermediate to handle pointer fields.
- **Command pattern:** kubectl-style `vn <resource> <verb> [args]`. Each resource has its own file defining the parent command, subcommands, and column definitions for table output.

## How to Add a New Resource

1. Create `internal/commands/<resource>.go`
2. Define the parent command function `newResourceCmd()` returning `*cobra.Command`
3. Add subcommand functions: `newResourceListCmd()`, `newResourceGetCmd()`, etc.
4. Define `resourceListColumns` and `resourceDetailColumns` as `[]output.Column`
5. In each subcommand's `RunE`, use `auth.NewClientFromContext(cmd)` to get the API client
6. Use `output.PrintList()` or `output.PrintItem()` for formatted output
7. Register in `root.go` via `cmd.AddCommand(newResourceCmd())`
8. Add tests in `internal/commands/commands_test.go` (verify command registration)

Follow the pattern in `calls.go` as the canonical example.

## Code Conventions

- One resource per file in `internal/commands/`
- Column definitions (`[]output.Column`) for both list and detail views
- Pagination via `--page-token` and `--page-size` flags on list commands
- Use `RunE` (not `Run`) for commands that perform API calls or I/O, so errors propagate correctly
- Error handling: check `resp.StatusCode() != 200` and return `fmt.Errorf("API error: %s", resp.Status())`
- Errors printed to stderr, non-zero exit code (via Cobra's `RunE` mechanism)
- No new external dependencies without strong justification
- Access keys are never logged or printed in output

## Dependencies

| Module | Purpose |
|--------|---------|
| `github.com/spf13/cobra` | CLI framework and command tree |
| `github.com/olekukonko/tablewriter` | Table-formatted output |
| `github.com/voipbin/voipbin-go` | VoIPBIN API SDK (local replace) |
| `gopkg.in/yaml.v3` | YAML output formatting |

## Testing

- Unit tests live alongside source: `*_test.go`
- Coverage: config load/save, auth key resolution, output formatting, command registration
- Run single package: `go test ./internal/config/`
- Run with verbose: `go test -v ./...`

## Gotchas

- **Local SDK dependency:** `voipbin-go` is referenced via `replace` directive in `go.mod` (`../voipbin-go`). The sibling directory must exist for builds to work.
- **Config file permissions:** Always 0600. Created automatically on `vn login`.
- **Access key validation:** `vn login` calls `GetMe` to validate the key before saving. Invalid keys are rejected.
- **YAML via JSON:** The YAML formatter marshals to JSON first, then to YAML. This ensures pointer fields and omitempty tags are handled correctly.
- **Version injection:** The `version` variable in `root.go` defaults to `"dev"`. Override at build time with `-ldflags`.
