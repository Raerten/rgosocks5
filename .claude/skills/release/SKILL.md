---
name: release
description: Prepare a new RGoSocks5 release — verify the tree is clean and green, summarize changes since the last tag, and create/push the next git tag that triggers the GoReleaser GitHub Actions workflow.
disable-model-invocation: true
---

# Release RGoSocks5

Releases are driven by **git tags**: pushing a tag matching `*` triggers
`.github/workflows/release.yml`, which runs GoReleaser to build multi-arch
binaries and publish Docker images (`raerten/rgosocks5`). This skill prepares and
cuts that tag. It has side effects (creating/pushing tags) — only run when the
user explicitly asks to release.

## Steps

1. **Pre-flight checks** — abort and report if any fail:
   - Working tree is clean: `git status --porcelain` (must be empty).
   - On the intended branch (usually `main`): `git branch --show-current`.
   - Tests and vet pass: `go test ./...` and `go vet ./...`.
   - `gofmt -l .` reports no files.

2. **Determine the next version**:
   - Find the latest tag: `git describe --tags --abbrev=0`.
   - Ask the user whether this is a patch / minor / major bump, or let them
     specify the exact version. Tags follow semver with a `v` prefix (e.g. `v1.2.7`).

3. **Summarize changes** since the last tag for the user to review:
   - `git log <last-tag>..HEAD --pretty=format:"%h %s"`.
   - Note that GoReleaser's changelog excludes `^docs:` and `^test:` commits, so
     highlight which commits will actually appear in release notes.

4. **Confirm with the user** the version and summary BEFORE tagging.

5. **Cut the release**:
   - `git tag <version>` (annotated tags are fine: `git tag -a <version> -m "<version>"`).
   - `git push origin <version>`.

6. **Report**: the new tag, what was pushed, and that the GoReleaser workflow is
   now running. Suggest the user watch GitHub Actions for build/publish status.

## Notes

- Version/commit/date are injected into the binary via `-ldflags` by GoReleaser
  (the `version`, `commit`, `date` vars in `main.go`). Don't hardcode versions.
- Never push a tag that already exists; verify with `git tag -l <version>` first.
- Do not `--force` push tags unless the user explicitly asks.
