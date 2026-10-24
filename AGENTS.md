# Repository Guidelines

## Project Structure & Module Organization

This calendar-booking project combines a Go backend with a React/TypeScript frontend.

- `cmd/server/main.go`: HTTP server entry point; currently serves `GET /health` on port 8080.
- `internal/`, `modules/`, and `config/`: placeholder directories for backend implementation and configuration.
- `web/src/`: React components, styles, and bundled assets in `assets/`; `web/public/` contains static assets.
- `mise.toml` and `mise/tasks.toml`: pinned tools and development tasks.
- `docs/`: project and release documentation; `.github/workflows/`: CI and release automation.

## Build, Test, and Development Commands

Run commands from the repository root:

- `mise trust` then `mise install`: trust the configuration and install pinned tools.
- `mise run web:install`: install frontend dependencies using the frozen pnpm lockfile.
- `mise run run`: start the Go server.
- `mise run web:dev`: start the Vite development server, usually on port 5173.
- `mise run build`: compile the backend to `bin/server`.
- `mise run test` / `mise run lint`: run Go tests / golangci-lint.
- `mise run web:build` / `mise run web:lint`: check TypeScript and build to `web/dist/` / run Oxlint.
- `mise run web:preview`: preview the built frontend.

## Coding Style & Naming Conventions

Format Go with `gofmt` and use standard Go naming conventions: lowercase packages and exported identifiers in PascalCase. Follow existing frontend style: two-space indentation, single-quoted JavaScript strings, and no semicolons. Use PascalCase for React component files, such as `App.tsx`, and camelCase for functions and variables. Use pnpm for frontend dependencies and commit corresponding lockfile updates.

## Testing Guidelines

No test files or frontend test framework are currently present, and no coverage threshold is configured. Add Go tests alongside implementation in `*_test.go` files using `testing` and `TestXxx` functions. Run `mise run test` before submitting backend changes. For frontend changes, run lint and build, and manually verify affected interactions; these checks do not replace behavioral tests.

## Commit & Pull Request Guidelines

Recent commits use Conventional Commits, including `ci:` and `chore:`. Use descriptive `feat:`, `fix:`, or `docs:` prefixes; release-please uses commit messages to determine releases. Keep squash-merge titles conventional. PR descriptions should explain behavior changes, link relevant issues, report validation, and include screenshots for UI changes.

Do not edit or delete `.github/workflows/hexlet-check.yml` or rename the repository. Keep credentials and local `.env` files out of Git; provide example configuration without secrets when needed.
