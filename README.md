# chp - Cherrypick CLI

A command-line client for the [Cherrypick](https://alpha.lollipopai.com/) API. Single static binary — no runtime dependencies.

## Install

### Homebrew (macOS/Linux)

```bash
brew install lollipopai/tap/chp
```

### Download binary

Grab the latest release from [GitHub Releases](https://github.com/lollipopai/cli/releases), extract, and put `chp` on your PATH.

### Go install

```bash
go install github.com/lollipopai/cli/cmd/chp@latest
```

### Build from source

```bash
git clone https://github.com/lollipopai/cli.git && cd cli
make build      # produces ./chp
make install    # copies to GOPATH/bin or /usr/local/bin
```

## Setup

```bash
chp login       # Opens browser for OAuth 2.1 PKCE login
chp whoami      # Show current user profile
chp logout      # Clear stored credentials
```

`chp login` opens your browser, catches the OAuth callback on a local server (`127.0.0.1:9876`), and stores tokens in `~/.chp/credentials.json`. Tokens auto-refresh when they expire.

## Commands

### Recipes

```bash
chp recipes search "chicken tikka"      # Search recipes by name
chp recipes get chicken-tikka-masala     # Get recipe by slug
chp recipes get 42                       # Get recipe by ID
```

### Products

Products are identified by their Sainsbury's product UID.

```bash
chp products search milk                # Search products by keyword
chp products get 7834128                # Get a product by Sainsbury's UID
```

### Basket

```bash
chp basket                              # Show current basket
chp basket show                         # Same as above

# Recipes (batch)
chp basket add-recipe 1 2 3             # Add multiple recipes at once
chp basket remove-recipe 1 2            # Remove multiple recipes

# Products (batch, with quantities)
chp basket add-product 7834128 7209381           # Add products (qty defaults to 1)
chp basket add-product 7834128:2 7209381:3       # Per-item quantities via uid:qty
chp basket add-product 7834128 7209381 -q 2      # -q sets default quantity for all
chp basket add-product 7834128:3 7209381 -q 2    # 7834128→3, 7209381→2 (explicit overrides -q)
chp basket remove-product 7834128 7209381        # Remove multiple products

# Quantity management
chp basket set-quantity 7834128 4        # Set quantity of a product already in basket

chp basket clear                         # Clear the entire basket
```

### Orders

```bash
chp orders                              # List order summaries
chp orders list                         # Same as above
chp orders get 42                       # Get order by ID (also prints product UIDs for re-ordering)
```

`orders get` prints the full order JSON followed by a summary of all Sainsbury's product UIDs found in the order, with a ready-to-use `chp basket add-product` command for re-ordering.

### Slots

Manage Sainsbury's delivery slots. Also available as `chp delivery`.

```bash
chp slots                               # List available delivery slots
chp slots list                          # Same as above
chp slots get 5                         # Get slot details
chp slots book 5                        # Book a delivery slot
chp delivery                            # Alias for chp slots
```

### Plan

Manage your weekly meal plan.

```bash
chp plan                                # Show current meal plan
chp plan show                           # Same as above
chp plan list                           # List available plans
chp plan get 1                          # Get a specific plan
chp plan add-recipe 1 100 101 102       # Add recipes to plan 1 (batch)
chp plan remove-recipe 1 100            # Remove recipe 100 from plan 1
```

### Playlists

```bash
chp playlists                           # List all playlists
chp playlists list                      # Same as above
chp playlists get 7                     # Get a specific playlist by ID
```

### Configuration

```bash
chp config show                         # Show current config (base URL, auth status, token expiry)
chp config set-url https://example.com  # Set the base API URL
```

### Raw Twirp calls

Call any Twirp RPC endpoint directly — useful for endpoints not wrapped by a named command:

```bash
chp call <service> <method> [json-payload]
```

The `lollipop.proto.` prefix is added automatically, so you just need the short service name:

```bash
chp call recipe.v1.RecipeV1 Search '{"query":"curry"}'
chp call user.v1.UserV1 Current
chp call basket.v1.BasketV1 Show
chp call product.v2.ProductV2 Search '{"keyword":"eggs"}'
chp call slot.v1.SlotV1 List
chp call plan.v1.PlanV1 Show
```

### Shell completions

```bash
chp completion bash > /etc/bash_completion.d/chp    # Bash
chp completion zsh > "${fpath[1]}/_chp"              # Zsh
chp completion fish > ~/.config/fish/completions/chp.fish  # Fish
```

## Credentials

Stored in `~/.chp/credentials.json` with `0600` permissions (directory `0700`).

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
brew uninstall chp          # if installed via Homebrew
rm "$(which chp)"           # if installed manually
rm -rf ~/.chp               # remove credentials and config
```

## Troubleshooting

| Error | Fix |
|-------|-----|
| `not logged in. Run: chp login` | Run `chp login` |
| `HTTP 401` / `Try: chp login` | Token expired — run `chp login` again |
| `Connection failed` | Server not running — check `chp config show` for the base URL |
| `invalid JSON response` | Server returned non-JSON — check the URL is correct |
| `OAuth state mismatch` | Possible CSRF — retry `chp login` |
| `failed to start callback server on port 9876` | Port in use — kill the process on 9876 and retry |

## Development

```bash
make build        # Build ./chp with version from git tags
make test         # Run all tests
make test-cover   # Run tests with coverage report
make lint         # Run go vet
make clean        # Remove build artifacts
```

## License

MIT
