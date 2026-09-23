# Graph Report - wd  (2026-09-22)

## Corpus Check
- Corpus is ~5,593 words - fits in a single context window. You may not need a graph.

## Summary
- 91 nodes · 229 edges · 14 communities (8 shown, 6 thin omitted)
- Extraction: 99% EXTRACTED · 1% INFERRED · 0% AMBIGUOUS · INFERRED: 2 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Root Command & Wiring
- CLI Verb Integration Tests
- Platform Directory Opener
- Warp Point Store
- Shell Wrapper Tests
- Open, Remove & Version
- Init & Show Commands
- List & Clean Commands
- Add Commands
- Store Unit Tests
- Test Process Helpers
- LS Command
- Go Module

## God Nodes (most connected - your core abstractions)
1. `Store` - 19 edges
2. `runWD()` - 12 edges
3. `runShellWrapperSuite()` - 12 edges
4. `tempDir()` - 11 edges
5. `writeWarpConfig()` - 8 edges
6. `runWDFailure()` - 5 edges
7. `TestShowMatchesSymlinkedDirectory()` - 5 edges
8. `TestCleanVerbRemovesMissingWarpPoints()` - 5 edges
9. `TestLSVerbListsWarpPointContents()` - 5 edges
10. `TestOpenVerbUsesPlatformOpener()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Execute()`  [EXTRACTED]
  main.go → cmd/root.go
- `TestWDHelperProcess()` --calls--> `Execute()`  [EXTRACTED]
  main_test.go → cmd/root.go
- `getWarpPoint()` --calls--> `newStore()`  [INFERRED]
  cmd/root.go → cmd/store.go
- `addPoint()` --calls--> `newStore()`  [INFERRED]
  cmd/root.go → cmd/store.go

## Import Cycles
- None detected.

## Communities (14 total, 6 thin omitted)

### Community 0 - "Root Command & Wiring"
Cohesion: 0.22
Nodes (9): addPoint(), Execute(), getWarpPoint(), passthroughArgs(), newStore(), go_pkg_slices, github.com/spf13/cobra.Command, main() (+1 more)

### Community 1 - "CLI Verb Integration Tests"
Cohesion: 0.51
Nodes (11): testing.T, runWD(), tempDir(), TestAddcdStoresAbsolutePath(), TestCleanVerbRemovesMissingWarpPoints(), TestInformationalVerbs(), TestLSVerbListsWarpPointContents(), TestOpenVerbUsesPlatformOpener() (+3 more)

### Community 2 - "Platform Directory Opener"
Cohesion: 0.28
Nodes (6): openDir(), openDir(), openDir(), go_pkg_os_exec, go_pkg_syscall, os/exec.Cmd

### Community 3 - "Warp Point Store"
Cohesion: 0.47
Nodes (3): Store, WarpPoint, go_pkg_bufio

### Community 4 - "Shell Wrapper Tests"
Cohesion: 0.39
Nodes (8): go_pkg_runtime, lastLine(), runShellWrapperSuite(), TestBashWrapperSuite(), TestPowerShellWrapperSuite(), TestZshWrapperSuite(), writeShim(), shellAdapter

### Community 6 - "Init & Show Commands"
Cohesion: 0.29
Nodes (3): sameDir(), go_pkg_strings, os.FileInfo

### Community 9 - "Store Unit Tests"
Cohesion: 0.40
Nodes (4): TestStoreLoadRejectsMalformedConfig(), TestStorePutReplacesExistingPoint(), TestStoreRemoveAndCleanWriteConsistentConfig(), go_pkg_testing

### Community 10 - "Test Process Helpers"
Cohesion: 0.50
Nodes (4): helperEnv(), runWDFailure(), TestCommandErrors(), invocation

## Knowledge Gaps
- **1 isolated node(s):** `github.com/OrganizedMayhem/wd`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 19 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Warp Point Store` to `Root Command & Wiring`, `Shell Wrapper Tests`, `Open, Remove & Version`, `Init & Show Commands`, `List & Clean Commands`, `Add Commands`, `LS Command`, `Path Command & Entry`?**
  _High betweenness centrality (0.191) - this node is a cross-community bridge._
- **What connects `github.com/OrganizedMayhem/wd` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._