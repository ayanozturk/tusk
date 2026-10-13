# Index and analysis throughput iteration

Date: 2026-10-13. Baseline: `cb6cc5c5`; candidate: `e9809ebe`. Host: Apple
MacBookPro17,1 (M1, 8 logical CPUs, 8 GiB), macOS 26.5.2, Go 1.27.1,
`GOWORK=off`. Workload: the checked-out WordPress corpus, all paths, 5,362
files, 1,451,710 lines, 47,359,328 bytes. Full runs produced 54,216
diagnostics with zero failed files. Raw samples are in
[the JSON report](2026-10-13-index-analysis-throughput.json).

The candidate adds a no-delimiter fast path to `splitTopLevelTypes`, avoiding
the nesting-state scan for single-part type strings. It does not alter parsing
or type-splitting results.

## Index throughput

Twenty alternating baseline/candidate runs measured the parent-observed cold
`index` worker process at eight workers. This includes process startup, file
discovery and reads, declaration parsing, project indexing, and result output.
It uses the same parent-boundary timing as the benchmark harness; the worker's
internal timer does not include those setup phases.

| Variant | Mean | CV | Throughput |
| --- | ---: | ---: | ---: |
| Baseline | 488.5 ms | 2.46% | 2,971,594 LOC/s |
| Candidate | 476.9 ms | 3.12% | 3,044,005 LOC/s |

The candidate improved the measured index rate by 2.44%. Both variants meet
the 5% CV gate and account for the same 5,362 files with no failures.

## Cold full-analysis throughput

The harness ran ten process-cold full-pipeline samples per variant at each
worker count, with one validation run and one warmup. This metric covers
startup, discovery, declaration parsing, indexing, full host parsing, semantic
snapshot construction, and analysis rules. The exact per-file issue snapshot
is unchanged: 3,187 analysis targets, zero parse errors, and 54,216 issue
entries before and after.

| Workers | Variant | Mean | CV | Throughput | Diagnostics |
| ---: | --- | ---: | ---: | ---: | ---: |
| 4 | Baseline | 3,048.5 ms | 0.71% | 476,205 LOC/s | 54,216 |
| 4 | Candidate | 3,068.3 ms | 0.75% | 473,132 LOC/s | 54,216 |
| 8 | Baseline | 2,518.4 ms | 0.79% | 576,441 LOC/s | 54,216 |
| 8 | Candidate | 2,516.5 ms | 1.26% | 576,877 LOC/s | 54,216 |

At a matched worker count, the candidate is neutral for cold full-analysis
throughput. On this eight-CPU host, using eight workers instead of four raised
the candidate's full-pipeline rate by 21.9%; the benchmark already defaults to
`runtime.NumCPU()`. Treat that worker-scaling result as host-specific.

## Checks and limits

- `GOWORK=off go test ./...` and `GOWORK=off go vet ./...` passed.
- The baseline and candidate WordPress issue snapshots matched exactly.
- The benchmark reports for both worker counts had exact file accounting and
  CV below 5%.
- This is one host and one corpus. It does not replace the release-wide
  resource baseline across all pinned workloads, and it does not claim a
  same-worker full-analysis speedup.

Reproduce the cold harness measurements with:

```sh
GOWORK=off go run ./cmd/benchmark --root test_projects/wordpress-develop --cold-runs 10 --extra-cold-runs 0 --skip-warm --workers 4 --json --output /tmp/tusk-wp-throughput.json
GOWORK=off go run ./cmd/benchmark --root test_projects/wordpress-develop --cold-runs 10 --extra-cold-runs 0 --skip-warm --workers 8 --json --output /tmp/tusk-wp-throughput-8.json
```
