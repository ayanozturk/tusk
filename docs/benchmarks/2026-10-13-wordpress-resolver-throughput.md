# WordPress resolver throughput iteration

Date: 2026-10-13. Baseline: `e24bdcf6`; candidate: `da1c63b1`. Host: Apple
MacBookPro17,1 (M1, 8 logical CPUs, 8 GiB), macOS 26.5.2, Go 1.27.1,
`GOWORK=off`. The WordPress workload discovered and parsed 5,362 files,
1,451,710 lines, and 47,359,328 bytes with eight workers. Runs emitted 54,216
diagnostics with zero failed files. Raw generated measurements remain locally
at `/tmp/tusk-analysis-iter3-lineage-interleaved.json` and
`/tmp/tusk-analysis-iter3-lineage-kernel.json`.

The change reuses the normalized lineage key in
`resolveMethodFromNonGenericLineage` for both class and method map lookups,
avoiding a second normalization of `resolved.Name` for each ancestor.

## Throughput results

Ten process-cold runs per version were interleaved through the benchmark
harness. CV was 1.11% for the baseline and 1.66% for the candidate. Cold
throughput includes process startup, discovery, parsing, indexing, semantic
snapshot work, and analysis rules.

| Variant | Mean | CV | Throughput | Max RSS | Diagnostics |
| --- | ---: | ---: | ---: | ---: | ---: |
| Baseline | 2,423.9 ms | 1.11% | 598,915 LOC/s | 447 MB | 54,216 |
| Candidate | 2,387.3 ms | 1.66% | 608,097 LOC/s | 432 MB | 54,216 |

The cold full-pipeline rate improved 1.54%. A second matched measure timed the
worker after discovery and parsing, covering `BuildProjectIndex` and
`runAnalysis`. Across ten alternating pairs, mean time fell from 2,019.7 ms to
1,971.6 ms (CV 0.64% and 1.42%), raising that index-and-analysis kernel from
718,775 to 736,311 LOC/s, a 2.44% gain. This kernel is not a rules-only timer.

The before/after WordPress diagnostic snapshots matched exactly: 3,187
analysis files, zero parse errors, and 54,216 issue entries. The interleaved
cold-run peak RSS fell from 447 MB to 432 MB.

## Validation and limits

- `GOWORK=off go test ./analyse`, `GOWORK=off go test ./...`, and
  `GOWORK=off go vet ./...` passed.
- Both variants accounted for all 5,362 files with zero failures and the same
  diagnostic count; the full issue snapshot matched exactly.
- This is one host and one corpus. Keep the release-wide performance gate open
  until the other pinned workloads have stable matched baselines.

Reproduce the interleaved cold runs by building the baseline harness from
`e24bdcf6`, then running the candidate harness:

```sh
git worktree add --detach /tmp/tusk-analysis-baseline e24bdcf6
(cd /tmp/tusk-analysis-baseline && GOWORK=off go build -o /tmp/tusk-analysis-baseline-bin ./cmd/benchmark)
git worktree remove /tmp/tusk-analysis-baseline
GOWORK=off go run ./cmd/benchmark --root test_projects/wordpress-develop --cold-runs 10 --extra-cold-runs 0 --skip-warm --workers 8 --baseline-binary /tmp/tusk-analysis-baseline-bin --json --output /tmp/tusk-analysis.json
```

For the timed index-and-analysis kernel, invoke each build with
`--harness-worker-phase=full --harness-worker-repeat=1`, alternating baseline
and candidate runs. Each worker prints `durationMs` and file, diagnostic, and
memory accounting as JSON.
