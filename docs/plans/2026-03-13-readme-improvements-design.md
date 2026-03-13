# Design: README.md Improvements

**Date:** 2026-03-13
**Status:** Approved

## Goal

Improve the README with structural additions for better developer experience, while keeping the existing (verified accurate) content intact.

## Approach

Structural improvements (Approach B): fix small gaps and add missing sections that new users and contributors need. No rewrite of existing content.

## Accuracy Verification (Completed)

Before designing changes, the full README was verified against the codebase:
- All 47 registered commands match their documented verbs (100% accuracy)
- All 4 global flags match root.go definitions
- Release v0.1.0 exists with all 4 binary assets (download URLs valid)
- Shell completion implementation matches documentation
- No orphaned command files; all registered in root.go

## Changes

### 1. Build from source — add Go version requirement

**Location:** Line 40, replace intro sentence

Before:
```
Building from source requires the [voipbin-go](...) SDK as a sibling directory:
```

After:
```
Requires Go 1.23+ and the [voipbin-go](...) SDK as a sibling directory:
```

**Rationale:** go.mod specifies `go 1.23.2`. Users building from source need to know this upfront.

### 2. New "Environment Variables" section

**Location:** After the Profiles section (after line 112), before Commands section

```markdown
## Environment Variables

| Variable | Description |
|----------|-------------|
| `VN_ACCESS_KEY` | API access key (overrides config file, overridden by `--access-key` flag) |
```

**Rationale:** `VN_ACCESS_KEY` is a supported env var (used in `internal/auth/auth.go`) but only mentioned in the access key priority list. Giving it its own section makes it discoverable for CI/CD and scripting use cases.

### 3. New "Troubleshooting" section

**Location:** End of file (after Shell Completion)

```markdown
## Troubleshooting

**"no access key found"**
Set an access key via one of: `vn login`, `--access-key` flag, or `VN_ACCESS_KEY` environment variable.

**"API error" on login**
The access key is validated against the API during login. Verify the key is correct and the API is reachable. If using a custom API URL, pass `--api-url` explicitly.

**Build fails with missing `voipbin-go` module**
The SDK is referenced via a local `replace` directive in `go.mod`. Clone [voipbin-go](https://github.com/voipbin/voipbin-go) as a sibling directory:

├── cli/
└── voipbin-go/
```

**Rationale:** Covers the three most common failure modes: missing auth, invalid key, and build setup. Error messages match actual strings from `internal/auth/auth.go`.

## Intentionally Omitted

- **License section:** No LICENSE file in the repo; nothing to reference
- **Contributing section:** Repo has 4 commits and no contribution guidelines; premature
- **Command reference changes:** Verified 100% accurate; no changes needed
- **Structural rewrite:** Current structure is solid; changes are additive only
