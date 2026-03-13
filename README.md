# vn

Command-line interface for the [VoIPBIN](https://voipbin.net) API. Manage calls, messages, agents, campaigns, and 40+ resources from your terminal.

## Overview

- **kubectl-style commands:** `vn <resource> <verb> [flags]`
- **Multiple output formats:** table (default), JSON, YAML
- **Multi-profile authentication:** switch between environments with `--profile`
- **Shell completion:** bash, zsh, fish, powershell

## Installation

### Download binary

Download the latest release for your platform:

**Linux (amd64)**
```bash
curl -sSL https://github.com/voipbin/cli/releases/latest/download/vn-linux-amd64 -o /usr/local/bin/vn && chmod +x /usr/local/bin/vn
```

**Linux (arm64)**
```bash
curl -sSL https://github.com/voipbin/cli/releases/latest/download/vn-linux-arm64 -o /usr/local/bin/vn && chmod +x /usr/local/bin/vn
```

**macOS (Apple Silicon)**
```bash
curl -sSL https://github.com/voipbin/cli/releases/latest/download/vn-darwin-arm64 -o /usr/local/bin/vn && chmod +x /usr/local/bin/vn
```

**macOS (Intel)**
```bash
curl -sSL https://github.com/voipbin/cli/releases/latest/download/vn-darwin-amd64 -o /usr/local/bin/vn && chmod +x /usr/local/bin/vn
```

### Build from source

Requires Go 1.23+ and the [voipbin-go](https://github.com/voipbin/voipbin-go) SDK as a sibling directory:

```bash
git clone https://github.com/voipbin/voipbin-go.git
git clone https://github.com/voipbin/cli.git
cd cli && make build
cp bin/vn /usr/local/bin/
```

## Quick Start

```bash
# Authenticate with your VoIPBIN access key
vn login

# List your calls
vn calls list

# Get details of a specific call
vn calls get <call-id>

# List agents in JSON format
vn agents list --output json

# Use a different profile
vn --profile staging calls list
```

## Configuration

Configuration is stored in `~/.vn/config.yaml`.

### Authentication

```bash
# Interactive login (prompts for access key)
vn login

# Non-interactive login
vn login --access-key <key> --profile production

# With a custom API URL
vn login --access-key <key> --profile staging --api-url https://staging-api.voipbin.net/v1.0

# Remove credentials
vn logout

# Remove credentials for a specific profile
vn logout --profile staging
```

The login command validates your access key against the API before saving. If the key is invalid, authentication will fail with an error and no credentials will be stored.

### Access key priority

The CLI resolves the access key in this order:

1. `--access-key` flag
2. `VN_ACCESS_KEY` environment variable
3. Config file profile (`~/.vn/config.yaml`)

### Profiles

Manage multiple environments with named profiles:

```bash
# Login to different environments
vn login --profile production
vn login --profile staging --api-url https://staging-api.voipbin.net/v1.0

# Use a specific profile
vn --profile staging agents list
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VN_ACCESS_KEY` | API access key (overrides config file, overridden by `--access-key` flag) |

## Commands

All commands follow the pattern `vn <resource> <verb> [args] [flags]`.

### Global flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `table`, `json`, `yaml` (default: `table`) |
| `--profile` | | Configuration profile to use |
| `--access-key` | | API access key (overrides config and env) |
| `--api-url` | | API base URL (overrides config) |

### Communication

| Resource | Verbs |
|----------|-------|
| `calls` | `list`, `get`, `create`, `delete`, `hangup`, `hold`, `unhold`, `mute`, `unmute`, `silence`, `unsilence`, `moh`, `unmoh`, `recording-start`, `recording-stop`, `talk` |
| `messages` | `list`, `get`, `create`, `delete` |
| `emails` | `list`, `get`, `create`, `delete` |
| `conferences` | `list`, `get`, `create`, `update`, `delete`, `recording-start`, `recording-stop`, `transcribe-start`, `transcribe-stop` |
| `conferencecalls` | `list`, `get`, `delete` |
| `groupcalls` | `list`, `get`, `create`, `delete`, `hangup` |
| `transfers` | `create` |

### AI Services

| Resource | Verbs |
|----------|-------|
| `ais` | `list`, `get`, `create`, `update`, `delete` |
| `aicalls` | `list`, `get`, `create`, `delete` |
| `aimessages` | `list`, `get`, `create`, `delete` |
| `aisummaries` | `list`, `get`, `create`, `delete` |

### Campaigns & Outbound

| Resource | Verbs |
|----------|-------|
| `campaigns` | `list`, `get`, `create`, `update`, `delete`, `set-status`, `set-service-level`, `set-actions`, `set-next-campaign`, `set-resource-info` |
| `campaigncalls` | `list`, `get`, `delete` |
| `outdials` | `list`, `get`, `create`, `update`, `delete`, `set-campaign`, `set-data`, `list-targets`, `create-target`, `delete-target` |
| `outplans` | `list`, `get`, `create`, `update`, `delete`, `set-dial-info` |

### Routing & Flow

| Resource | Verbs |
|----------|-------|
| `flows` | `list`, `get`, `create`, `update`, `delete` |
| `activeflows` | `list`, `get`, `create`, `delete`, `stop` |
| `routes` | `list`, `get`, `create`, `update`, `delete` |
| `queues` | `list`, `get`, `create`, `update`, `delete`, `set-routing-method`, `set-tag-ids` |
| `queuecalls` | `list`, `get`, `delete`, `kick` |
| `extensions` | `list`, `get`, `create`, `update`, `delete` |

### Chat & Conversation

| Resource | Verbs |
|----------|-------|
| `chats` | `list`, `get`, `create`, `update`, `delete`, `add-participant`, `remove-participant`, `set-room-owner` |
| `chatrooms` | `list`, `get`, `create`, `update`, `delete` |
| `chatmessages` | `list`, `get`, `create`, `delete` |
| `chatroommessages` | `list`, `get`, `create`, `delete` |
| `conversations` | `list`, `get`, `update`, `list-messages`, `create-message` |
| `conversation-accounts` | `list`, `get`, `create`, `update`, `delete` |

### Account Management

| Resource | Verbs |
|----------|-------|
| `agents` | `list`, `get`, `create`, `update`, `delete`, `update-addresses`, `update-password`, `update-permission`, `update-status`, `update-tag-ids` |
| `customers` | `list`, `get`, `create`, `update`, `delete`, `update-billing-account` |
| `accesskeys` | `list`, `get`, `create`, `update`, `delete` |
| `tags` | `list`, `get`, `create`, `update`, `delete` |

### Telecom

| Resource | Verbs |
|----------|-------|
| `numbers` | `list`, `get`, `create`, `update`, `delete`, `renew`, `update-flow-ids` |
| `available-numbers` | `list` |
| `providers` | `list`, `get`, `create`, `update`, `delete` |
| `trunks` | `list`, `get`, `create`, `update`, `delete` |

### Media & Storage

| Resource | Verbs |
|----------|-------|
| `recordings` | `list`, `get`, `delete` |
| `transcribes` | `list`, `get`, `create`, `delete`, `stop` |
| `files` | `list`, `get`, `delete` |
| `storage-accounts` | `list`, `get`, `create`, `delete` |
| `storage-files` | `list`, `get`, `delete` |

### Billing

| Resource | Verbs |
|----------|-------|
| `billing-accounts` | `list`, `get`, `create`, `update`, `delete`, `balance-add`, `balance-subtract`, `update-payment-info` |
| `billings` | `list`, `get` |

### Utility

| Command | Description |
|---------|-------------|
| `me` | Show current authenticated user |
| `version` | Print CLI version |
| `login` | Authenticate and store credentials |
| `logout` | Remove stored credentials |
| `completion` | Generate shell completion scripts |

### Pagination

List commands support pagination. When `--page-size` is not specified, the API uses its default page size.

```bash
vn calls list --page-size 50
vn calls list --page-size 50 --page-token <token-from-previous-response>
```

## Output Formats

### Table (default)

```
$ vn calls list
ID                                    SOURCE          DESTINATION     DIRECTION  STATUS   CREATED
550e8400-e29b-41d4-a716-446655440000  +12125551234    +14155559876    outbound   hangup   2026-03-05T10:30:00Z
```

### JSON

```bash
$ vn calls list --output json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "source": "+12125551234",
    "destination": "+14155559876",
    ...
  }
]
```

### YAML

```bash
$ vn calls list --output yaml
- id: 550e8400-e29b-41d4-a716-446655440000
  source: "+12125551234"
  destination: "+14155559876"
  ...
```

## Shell Completion

```bash
# Bash
source <(vn completion bash)

# Zsh
source <(vn completion zsh)

# Fish
vn completion fish | source

# PowerShell
vn completion powershell | Invoke-Expression
```

To make it permanent, add the appropriate line to your shell's config file (e.g., `~/.bashrc`, `~/.zshrc`).

## Troubleshooting

**"no access key found"**

Set an access key via one of: `vn login`, `--access-key` flag, or `VN_ACCESS_KEY` environment variable.

**"API error" on login**

The access key is validated against the API during login. Verify the key is correct and the API is reachable. If using a custom API URL, pass `--api-url` explicitly.

**Build fails with missing `voipbin-go` module**

The SDK is referenced via a local `replace` directive in `go.mod`. Clone [voipbin-go](https://github.com/voipbin/voipbin-go) as a sibling directory:

```
├── cli/
└── voipbin-go/
```
