# Repository Guidelines

## Project Structure & Module Organization

This calendar-booking project combines a Go backend with a React/TypeScript frontend.

- `cmd/server/main.go`: HTTP server entry point on port 8080; generated router and handwritten handlers in `internal/booking/`.
- `typespec/`: authoritative API contract. `openapi/`, `gen/api/`, and `web/src/generated/api/` are committed generated artifacts; never edit them manually.
- Change TypeSpec, run `mise run generate`, and commit contract and all generated changes together. `mise run generate:check` detects drift.
- `internal/booking/`: in-memory booking implementation and HTTP tests.
- `web/src/`: React components, styles, and bundled assets in `assets/`; `web/public/` contains static assets.
- `mise.toml` and `mise/tasks.toml`: pinned tools and development tasks.
- `docs/`: project and release documentation; `.github/workflows/`: CI and release automation.

## Build, Test, and Development Commands

Run commands from the repository root:

- `mise trust` then `mise install`: trust the configuration and install pinned tools.
- `mise run install`: install workspace and Go dependencies. `web:install` remains an alias for frozen workspace pnpm installation.
- `mise run generate`: TypeSpec → OpenAPI → Go server and TS SDK. Individual `generate:typespec`, `generate:backend`, and `generate:frontend` tasks are available; backend/frontend consume existing OpenAPI.
- `mise run generate:check`: regenerate and check artifact drift, including missing and untracked files.
- `go test -race ./...`: verify concurrent in-memory storage.
- `mise run run`: start the Go server.
- `mise run web:dev`: start the Vite development server, usually on port 5173.
- `mise run build`: compile the backend to `bin/server`.
- `mise run test` / `mise run lint`: run Go tests / golangci-lint.
- `mise run web:build` / `mise run web:lint`: check TypeScript and build to `web/dist/` / run Oxlint.
- `mise run web:preview`: preview the built frontend.

## Coding Style & Naming Conventions

Format Go with `gofmt` and use standard Go naming conventions: lowercase packages and exported identifiers in PascalCase. Follow existing frontend style: two-space indentation, single-quoted JavaScript strings, and no semicolons. Use PascalCase for React component files, such as `App.tsx`, and camelCase for functions and variables. Use pnpm for frontend dependencies and commit corresponding lockfile updates.

## Testing Guidelines

HTTP tests cover the booking API. No frontend test framework or coverage threshold is configured. Add Go tests alongside implementation in `*_test.go` files using `testing` and `TestXxx` functions. Run `mise run test` before submitting backend changes. For frontend changes, run lint and build, and manually verify affected interactions; these checks do not replace behavioral tests.

## Commit & Pull Request Guidelines

Recent commits use Conventional Commits, including `ci:` and `chore:`. Use descriptive `feat:`, `fix:`, or `docs:` prefixes; release-please uses commit messages to determine releases. Keep squash-merge titles conventional. PR descriptions should explain behavior changes, link relevant issues, report validation, and include screenshots for UI changes.

Do not edit or delete `.github/workflows/hexlet-check.yml` or rename the repository. Keep credentials and local `.env` files out of Git; provide example configuration without secrets when needed.

## Agent skills

### Issue tracker

Issues and specs live in GitHub Issues for
`one-thousand-twoth/ai-for-developers-project-386`.
See `docs/agents/issue-tracker.md`.

### Triage labels

Use the five default triage labels.
See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout: root `GLOSSARY.md` and `docs/adr/`.
See `docs/agents/domain.md`.
