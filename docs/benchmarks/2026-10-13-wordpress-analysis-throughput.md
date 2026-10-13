# WordPress analysis throughput iteration

Date: 2026-10-13. Baseline: `e9809ebe`; candidate: `e24bdcf6`. Host: Apple
MacBookPro17,1 (M1, 8 logical CPUs, 8 GiB), macOS 26.5.2, Go 1.27.1,
`GOWORK=off`. The WordPress workload discovered and parsed 5,362 files,
1,451,710 lines, and 47,359,328 bytes with eight workers. Full runs emitted
54,216 diagnostics with zero failed files. Raw generated samples are retained
locally at `/tmp/tusk-analysis-interleaved.json` and
`/tmp/tusk-analysis-paired.json`.

The change shortens shard selection in `asciiLowerIdent`'s shared bounded cache.
It mixes the key length and first, middle, and last bytes instead of hashing
every byte with FNV-1a. The cache map still hashes the complete key; only the
choice among eight locks changes.

## Throughput results

Ten process-cold runs per version were interleaved through the benchmark
harness. CV was 1.23% for the baseline and 1.05% for the candidate. The cold
metric includes process startup, discovery, parsing, indexing, semantic
snapshot work, and analysis rules.

| Variant | Mean | CV | Throughput | Max RSS | Diagnostics |
| --- | ---: | ---: | ---: | ---: | ---: |
| Baseline | 2,522.9 ms | 1.23% | 575,413 LOC/s | 433 MB | 54,216 |
| Candidate | 2,417.6 ms | 1.05% | 600,476 LOC/s | 424 MB | 54,216 |

That is a 4.36% cold full-pipeline throughput gain. A second matched measure
timed the worker after file discovery and parsing, covering `BuildProjectIndex`
and `runAnalysis`. Across ten alternating pairs, mean time fell from 2,116.3 ms
to 2,028.0 ms (CV 1.21% and 3.05%), raising that index-and-analysis kernel from
685,966 to 715,833 LOC/s, also a 4.35% gain. This kernel is not a rules-only
timer.

The before/after WordPress diagnostic snapshots matched exactly after
normalizing the corpus root: 3,187 analysis files, zero parse errors, and
54,216 issue entries. The interleaved cold-run peak RSS was also slightly lower
for the candidate (424 MB versus 433 MB).

## Validation and limits

- `GOWORK=off go test ./analyse`, `GOWORK=off go test ./...`, and
  `GOWORK=off go vet ./...` passed.
- Both variants accounted for all 5,362 files with zero failures and the same
  diagnostic count; the full snapshot comparison was exact.
- This is one host and one corpus. Carry the candidate forward for the next
  iteration, but keep the release-wide performance gate open until the other
  pinned workloads have stable matched baselines.

Reproduce the interleaved cold runs by building the baseline harness binary
from `e9809ebe`, then running the candidate harness:

```sh
git worktree add --detach /tmp/tusk-analysis-baseline e9809ebe
(cd /tmp/tusk-analysis-baseline && GOWORK=off go build -o /tmp/tusk-analysis-baseline-bin ./cmd/benchmark)
git worktree remove /tmp/tusk-analysis-baseline
GOWORK=off go run ./cmd/benchmark --root test_projects/wordpress-develop --cold-runs 10 --extra-cold-runs 0 --skip-warm --workers 8 --baseline-binary /tmp/tusk-analysis-baseline-bin --json --output /tmp/tusk-analysis.json
```

For the timed index-and-analysis kernel, invoke each build with
`--harness-worker-phase=full --harness-worker-repeat=1`, alternating baseline
and candidate runs. The worker emits its measured `durationMs` and file,
diagnostic, and memory accounting as JSON.
