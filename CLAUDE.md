# cpk - CherryPick CLI

## Project Overview

Go CLI tool for the CherryPick API. Communicates via Twirp RPC (protobuf-over-HTTP/JSON) to the CherryPick Rails backend.

## Architecture

```
cmd/cpk/main.go          Entrypoint, sets version from ldflags
internal/
  output/                 Info/Success/Warn/Error/Fatal + colored JSON
  httpclient/             HTTP client, localhost TLS skip, APIError
  auth/                   Credentials, OAuth 2.1 PKCE, token refresh
  twirp/                  Twirp RPC caller with auto-refresh
  cli/                    Cobra commands (one file per command group)
```

## Code Conventions

- **Go standard**: `camelCase` unexported, `PascalCase` exported
- **Package-level vars** for cobra commands: `var fooCmd = &cobra.Command{...}`
- **Error handling**: Return errors from library code, call `output.Fatal()` in CLI handlers
- **No structured logging**: CLI uses `output.Info/Warn/Error/Fatal` only
- **API responses**: Typed as `any`, printed as colored JSON via `output.PrintJSON`

## Security Rules

- **Never disable TLS verification for non-localhost URLs** — only `localhost` and `127.0.0.1` skip cert checks (see `httpclient.Client.clientFor`)
- **Credentials file uses 0600 permissions** via `os.OpenFile` — never `os.WriteFile` then `chmod`
- **Config dir uses 0700 permissions** via `os.MkdirAll`
- **HTML-escape server content** in OAuth callback page — uses `html.EscapeString`
- **Use `crypto/rand`** for all random token generation — never `math/rand`

## Dependencies (3 runtime + 1 test)

| Dep | Purpose |
|-----|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | Colored output + TTY detection |
| `github.com/pkg/browser` | Cross-platform browser opening for OAuth |
| `github.com/stretchr/testify` | Test assertions (test only) |

## Adding New Commands

1. Create `internal/cli/<name>.go` with `var <name>Cmd = &cobra.Command{...}`
2. For parent commands, add subcommands in `init()`
3. Register top-level commands in `root.go`'s `init()`
4. Use `newTwirpCaller()` for authenticated API calls
5. Use `output.PrintJSON(result)` for response output
6. Use `output.Fatal(err.Error())` for errors

## Twirp Service Naming

Services follow: `lollipop.proto.<domain>.<version>.<ServiceName>`

Examples:
- `lollipop.proto.user.v1.UserV1`
- `lollipop.proto.recipe.v1.RecipeV1`
- `lollipop.proto.product.v2.ProductV2`
- `lollipop.proto.basket.v1.BasketV1`

Use `cpk call <service> <method> [payload]` for any endpoint not wrapped by a named command.

## Testing

- Uses `testing` + `github.com/stretchr/testify`
- Mock HTTP with `net/http/httptest`
- Credentials tests use `t.TempDir()` to isolate file system
- PKCE tests verify against RFC 7636 test vectors
- Run: `make test` or `go test ./...`

## Building

```bash
make build          # Builds ./cpk with version from git
make test           # Runs all tests
make install        # Copies to GOPATH/bin or /usr/local/bin
```

Version is injected via ldflags: `-X main.version=<tag>`
