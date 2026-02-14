# cpk - CherryPick CLI

A command-line client for the [CherryPick](https://www.cherrypick.co/) API. Single static binary — no runtime dependencies.

## Install

### Homebrew (macOS/Linux)

```bash
brew install lollipopai/tap/cpk
```

### Download binary

Grab the latest release from [GitHub Releases](https://github.com/lollipopai/cli/releases), extract, and put `cpk` on your PATH.

### Go install

```bash
go install github.com/lollipopai/cli/cmd/cpk@latest
```

### Build from source

```bash
git clone https://github.com/lollipopai/cli.git && cd cli
make build      # produces ./cpk
make install    # copies to GOPATH/bin or /usr/local/bin
```

## Setup

Point the CLI at your CherryPick instance (defaults to `https://localhost:3000`):

```bash
cpk config set-url https://app.cherrypick.co
```

## Authentication

```bash
cpk login       # Opens browser for OAuth 2.1 PKCE login
cpk whoami      # Show current user profile
cpk logout      # Clear stored credentials
```

`cpk login` opens your browser, catches the OAuth callback on a local server (`127.0.0.1:9876`), and stores tokens in `~/.cpk/credentials.json`. Tokens auto-refresh when they expire.

## Commands

### Recipes

```bash
cpk recipes search "chicken tikka"      # Search recipes by name
cpk recipes get chicken-tikka-masala     # Get recipe by slug
cpk recipes get 42                       # Get recipe by ID
```

### Products

```bash
cpk products search milk                # Search products by keyword
cpk products get 123                    # Get a product by ID
```

### Basket

```bash
cpk basket                              # Show current basket
cpk basket show                         # Same as above
cpk basket add-recipe 456               # Add a recipe to the basket
cpk basket remove-recipe 456            # Remove a recipe from the basket
cpk basket add-product 789              # Add a product to the basket
cpk basket add-product 789 -q 2         # Add product with quantity
cpk basket remove-product 789           # Remove a product from the basket
cpk basket clear                        # Clear the entire basket
```

### Orders

```bash
cpk orders                              # List order summaries
cpk orders list                         # Same as above
cpk orders get 42                       # Get a specific order by ID
```

### Playlists

```bash
cpk playlists                           # List all playlists
cpk playlists list                      # Same as above
cpk playlists get 7                     # Get a specific playlist by ID
```

### Configuration

```bash
cpk config show                         # Show current config (base URL, auth status, token expiry)
cpk config set-url https://example.com  # Set the base API URL
```

### Raw Twirp calls

Call any Twirp RPC endpoint directly — useful for endpoints not wrapped by a named command:

```bash
cpk call <service> <method> [json-payload]
```

Examples:

```bash
cpk call lollipop.proto.recipe.v1.RecipeV1 Search '{"query":"curry"}'
cpk call lollipop.proto.user.v1.UserV1 Current
cpk call lollipop.proto.basket.v1.BasketV1 Show
cpk call lollipop.proto.product.v2.ProductV2 Search '{"keyword":"eggs"}'
```

Services follow the pattern `lollipop.proto.<domain>.<version>.<ServiceName>`.

### Shell completions

```bash
cpk completion bash > /etc/bash_completion.d/cpk    # Bash
cpk completion zsh > "${fpath[1]}/_cpk"              # Zsh
cpk completion fish > ~/.config/fish/completions/cpk.fish  # Fish
```

## Credentials

Stored in `~/.cpk/credentials.json` with `0600` permissions (directory `0700`).

The JSON format is identical to the previous Python CLI, so existing credentials carry over — no need to re-authenticate after upgrading.

Fields stored:

| Field | Description |
|-------|-------------|
| `base_url` | API base URL |
| `oauth_access_token` | Current OAuth access token |
| `oauth_refresh_token` | OAuth refresh token (for auto-renewal) |
| `oauth_expires_at` | Token expiry timestamp |
| `oauth_client_id` | Dynamically registered OAuth client ID |
| `jwt` | Legacy JWT token (read if present, not created by new login) |

## Uninstall

```bash
brew uninstall cpk          # if installed via Homebrew
rm "$(which cpk)"           # if installed manually
rm -rf ~/.cpk               # remove credentials and config
```

## Troubleshooting

| Error | Fix |
|-------|-----|
| `not logged in. Run: cpk login` | Run `cpk login` |
| `HTTP 401` / `Try: cpk login` | Token expired — run `cpk login` again |
| `Connection failed` | Server not running — check `cpk config show` for the base URL |
| `invalid JSON response` | Server returned non-JSON — check the URL is correct |
| `OAuth state mismatch` | Possible CSRF — retry `cpk login` |
| `failed to start callback server on port 9876` | Port in use — kill the process on 9876 and retry |

## Development

```bash
make build        # Build ./cpk with version from git tags
make test         # Run all tests
make test-cover   # Run tests with coverage report
make lint         # Run go vet
make clean        # Remove build artifacts
```

## License

MIT
