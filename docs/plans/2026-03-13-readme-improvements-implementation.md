# README Improvements Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Go version requirement, Environment Variables section, and Troubleshooting section to README.md.

**Architecture:** Three additive edits to a single file. No code changes, no tests needed.

**Tech Stack:** Markdown

---

### Task 1: Add Go version requirement to "Build from source"

**Files:**
- Modify: `README.md:40`

**Step 1: Edit the intro sentence**

Replace:
```
Building from source requires the [voipbin-go](https://github.com/voipbin/voipbin-go) SDK as a sibling directory:
```

With:
```
Requires Go 1.23+ and the [voipbin-go](https://github.com/voipbin/voipbin-go) SDK as a sibling directory:
```

**Step 2: Verify the edit**

Run: `head -45 README.md | tail -10`
Expected: Line 40 shows "Requires Go 1.23+ and the..."

---

### Task 2: Add "Environment Variables" section after Profiles

**Files:**
- Modify: `README.md` (insert after line 112, before the `## Commands` section)

**Step 1: Insert the new section**

Add after the Profiles section (after the closing ` ``` ` of the profiles code block):

```markdown

## Environment Variables

| Variable | Description |
|----------|-------------|
| `VN_ACCESS_KEY` | API access key (overrides config file, overridden by `--access-key` flag) |
```

**Step 2: Verify the edit**

Run: `grep -n "Environment Variables" README.md`
Expected: Shows the new section heading between Profiles and Commands

---

### Task 3: Add "Troubleshooting" section at end of file

**Files:**
- Modify: `README.md` (append after Shell Completion section)

**Step 1: Append the new section**

Add at end of file, after "To make it permanent..." paragraph:

```markdown

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
```

**Step 2: Verify the edit**

Run: `tail -20 README.md`
Expected: Shows the Troubleshooting section at end of file

---

### Task 4: Commit

**Step 1: Commit the change**

```bash
git add README.md
git commit -m "Improve README with Go version, env vars, and troubleshooting"
```
