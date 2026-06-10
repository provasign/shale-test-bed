# shale-test-bed

`shale-test-bed` is a deliberately small Go service used to exercise
[Shale](https://github.com/provasign/shale), Provasign's agent evidence workflow,
on real pull requests.

The application is intentionally simple: it gives agents and reviewers a
real codebase with source files, tests, runtime behavior, and plausible defects,
without the noise of a production service. Changes in this repository are meant
to create clear evidence about what an agent intended to do, what it changed,
what commands it ran, and whether the resulting checks passed or failed.

## Service behavior

The service exposes two HTTP endpoints:

```text
GET  /health
POST /login   {"user": "...", "password": "..."}
```

`GET /health` returns `200 OK` with `ok`.

`POST /login` checks credentials against a small in-memory user table in
`internal/auth`. Failed login attempts are tracked in memory so account lockout
behavior can be changed, tested, broken, and reviewed in pull requests.
Passwords must be 12 to 128 characters long and include lowercase, uppercase,
numeric, and special characters before they can authenticate.

The service is not intended for production use. Passwords are stored in memory
as plain test fixtures, state resets when the process restarts, and there is no
persistent datastore.

## Repository layout

```text
.
├── main.go                    # HTTP server and route handlers
├── internal/auth/auth.go      # Toy authentication and lockout logic
├── internal/auth/auth_test.go # Unit tests for authentication behavior
├── go.mod                     # Go module definition
├── AGENTS.md                  # Shale instructions for coding agents
└── CLAUDE.md                  # Same Shale instructions for Claude-based agents
```

## Running locally

Start the service:

```sh
go run .
```

The default address is `:8080`. Set `PORT` to use a different port:

```sh
PORT=9090 go run .
```

Run the test suite:

```sh
go test ./...
```

## What counts as test evidence here

This repository is designed so pull requests can show several kinds of evidence:

- Source changes: diffs to `main.go`, `internal/auth/auth.go`, tests, or docs.
- Unit tests: `internal/auth/auth_test.go` captures expected authentication and
  lockout behavior.
- Command evidence: agent-run commands such as `go test ./...` are recorded by
  Shale and attached to the agent session.
- Check outcomes: passing and failing test runs are both useful evidence because
  they show what was verified and what still needs attention.
- Coverage gaps: small, isolated code paths make it easy to see when a change is
  not covered by tests.
- Agent intent and completion notes: agents must record `shale intent` before
  editing and `shale done` after completing work, as described in `AGENTS.md`.

In practice, the most useful evidence for a PR is a combination of the code diff,
the updated or existing tests, the exact verification commands that were run, and
the final Shale evidence card generated from the agent session.

## Agent workflow

Agents working in this repository must follow the Shale instructions in
`AGENTS.md` before making edits. The short version is:

```sh
shale intent "one-line description of the planned change" --body "why and approach"
# make the change and run relevant checks
shale done --note "what changed and any deviations"
```

This keeps pull requests reviewable even when multiple agents or tools are used
to produce a change.
