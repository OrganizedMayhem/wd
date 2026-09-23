# Graph Report - wd  (2026-09-22)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 122 nodes · 298 edges · 15 communities (10 shown, 5 thin omitted)
- Extraction: 99% EXTRACTED · 1% INFERRED · 0% AMBIGUOUS · INFERRED: 3 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7d221828`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- CLI Integration Tests
- Root Command & Entry
- CI & Release Pipelines
- Platform Directory Opener
- Warp Point Store
- Ls & Version Commands
- Warp Point Completion
- Show Command
- Product Concepts (README)
- Clean & List Commands
- Add Commands & Store Tests
- Pages Graph Publishing
- Open Command
- Go Module
- Init & Path Commands

## God Nodes (most connected - your core abstractions)
1. `Store` - 19 edges
2. `runShellWrapperSuite()` - 13 edges
3. `runWD()` - 13 edges
4. `tempDir()` - 13 edges
5. `completeWarpPoint()` - 11 edges
6. `writeWarpConfig()` - 9 edges
7. `CI build job` - 6 edges
8. `Release build job` - 6 edges
9. `wd (warp directory)` - 6 edges
10. `completions()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `Cross-platform build matrix (linux/darwin/windows x amd64/arm64)` --conceptually_related_to--> `wd (warp directory)`  [INFERRED]
  .github/workflows/ci.yml → README.md
- `main()` --calls--> `Execute()`  [EXTRACTED]
  main.go → cmd/root.go
- `TestWDHelperProcess()` --calls--> `Execute()`  [EXTRACTED]
  main_test.go → cmd/root.go
- `CI build job` --semantically_similar_to--> `Release build job`  [INFERRED] [semantically similar]
  .github/workflows/ci.yml → .github/workflows/release.yml
- `CI checksums job` --semantically_similar_to--> `Release publish job (checksums + GitHub release)`  [INFERRED] [semantically similar]
  .github/workflows/ci.yml → .github/workflows/release.yml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **CI pipeline: test -> build -> checksums** — _github_workflows_ci_test, _github_workflows_ci_build, _github_workflows_ci_checksums [EXTRACTED 1.00]
- **Release pipeline: validate-tag -> build -> release** — _github_workflows_release_validate_tag, _github_workflows_release_build, _github_workflows_release_release_job [EXTRACTED 1.00]
- **wd shell integration (wrapper, completion, fzf picker)** — readme_shell_wrapper, readme_tab_completion, readme_fzf_picker [INFERRED 0.85]

## Communities (15 total, 5 thin omitted)

### Community 0 - "CLI Integration Tests"
Cohesion: 0.25
Nodes (26): os/exec.Cmd, testing.T, completions(), helperEnv(), lastLine(), runShellWrapperSuite(), runWD(), runWDFailure() (+18 more)

### Community 1 - "Root Command & Entry"
Cohesion: 0.24
Nodes (8): addPoint(), Execute(), getWarpPoint(), passthroughArgs(), newStore(), go_pkg_slices, main(), TestWDHelperProcess()

### Community 2 - "CI & Release Pipelines"
Cohesion: 0.29
Nodes (12): CI build job, CI checksums job, CI Workflow, Cross-platform build matrix (linux/darwin/windows x amd64/arm64), CI security job (Trivy scan), CI test job (vet + race tests, OS matrix), Trivy filesystem scan, Version injection via -ldflags cmd.version (+4 more)

### Community 3 - "Platform Directory Opener"
Cohesion: 0.28
Nodes (7): openDir(), openDir(), runOpener(), openDir(), go_pkg_syscall, go_pkg_unsafe, github.com/spf13/cobra.Command

### Community 4 - "Warp Point Store"
Cohesion: 0.47
Nodes (3): Store, WarpPoint, go_pkg_bufio

### Community 6 - "Warp Point Completion"
Cohesion: 0.40
Nodes (3): completeWarpPoint(), github.com/spf13/cobra.Completion, github.com/spf13/cobra.ShellCompDirective

### Community 7 - "Show Command"
Cohesion: 0.33
Nodes (3): sameDir(), go_pkg_strings, os.FileInfo

### Community 8 - "Product Concepts (README)"
Cohesion: 0.48
Nodes (7): wd commands (add, list, rm, show, open, ls, path, clean), fzf warp point picker (wd with no args), mfaerevaag/wd zsh plugin, Shell wrapper (wd init), Tab completion (wd completion), Warp point, wd (warp directory)

### Community 10 - "Add Commands & Store Tests"
Cohesion: 0.18
Nodes (7): TestStoreLoadRejectsMalformedConfig(), TestStorePutReplacesExistingPoint(), TestStoreRemoveAndCleanWriteConsistentConfig(), go_pkg_os, go_pkg_path_filepath, go_pkg_runtime, go_pkg_testing

### Community 11 - "Pages Graph Publishing"
Cohesion: 0.67
Nodes (3): Pages deploy job (stage graphify-out site), graphify-out outputs (graph.html, graph.json, GRAPH_REPORT.md), Pages Workflow

## Knowledge Gaps
- **5 isolated node(s):** `github.com/OrganizedMayhem/wd`, `Pages Workflow`, `graphify-out outputs (graph.html, graph.json, GRAPH_REPORT.md)`, `Trivy filesystem scan`, `mfaerevaag/wd zsh plugin`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 24 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Warp Point Store` to `CLI Integration Tests`, `Root Command & Entry`, `Ls & Version Commands`, `Show Command`, `Clean & List Commands`, `Add Commands & Store Tests`, `Init & Path Commands`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Why does `completeWarpPoint()` connect `Warp Point Completion` to `Root Command & Entry`, `Platform Directory Opener`, `Ls & Version Commands`, `Show Command`, `Open Command`, `Init & Path Commands`?**
  _High betweenness centrality (0.043) - this node is a cross-community bridge._
- **What connects `github.com/OrganizedMayhem/wd`, `Pages Workflow`, `graphify-out outputs (graph.html, graph.json, GRAPH_REPORT.md)` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._