# Reviewed corpus diagnostic quality

## Gate status

**Failed; release-blocking work remains.** The checked-in differential fixtures
pass, and the fresh PSL comparison analyzed every selected file. The corpus
differences have not all been reviewed and dispositioned. Symfony and Composer
candidate snapshots still have unreviewed issue-set changes, and Symfony has
one parse-error file in both snapshots. The cumulative comparison against the
saved `HEAD` snapshots has not been rerun after the latest fixes.

The analyzer correctness gate is conjunctive. Every condition below must pass;
aggregate scores cannot offset a failed condition:

1. **Differential fixtures:** zero engine mismatches and zero PHPStan reference
   mismatches in every pinned supported case. Unsupported cases must match
   their pinned PHPStan result and have an explicit unsupported disposition.
2. **Corpus accounting:** every selected PHP file is discovered and analyzed by
   both tools, with zero read, parse, timeout, or analysis failures.
3. **Corpus diagnostics:** every engine-only and PHPStan-only finding is
   reviewed against the selected tool versions and corpus revision. Every
   difference has a checked-in disposition; none remain unexplained at release
   review.
4. **Metric integrity:** publish precision, recall, F1, reviewed F1, unmatched
   counts, and provenance together for each workload and level. F1 is a
   workload-specific summary, not a standalone pass threshold; do not average
   across levels or corpora.

An intentional behavior boundary is acceptable only when it is explicitly
dispositioned and reflected in the supported limitations and release notes.

## Remaining work

Use this order for the remaining corpus-quality work. The gate stays open until
all differences have a source-backed disposition and the release-candidate
corpora pass the accounting and differential checks.

1. **Disposition every current PSL `A.ARG.TYPE` location.** The 2026-10-31
   20:34:43Z report has 150 engine-only argument diagnostics. Review each
   source call and its declared contract, then mark it as an analyser defect,
   an intentional supported boundary, or an unsupported construct. Repeated
   patterns include scalar ranges, non-empty strings, nullable refinements,
   interface/resource unions, and generic substitutions; these are only search
   buckets. Do not classify a finding from its message or bucket alone. The
   location inventory below pins the work to this report. The 356/109 counts
   quoted in the task do not match this checkout's latest reproducible report
   (150/7); this inventory is tied to the pinned 1,775-file PSL corpus and its
   2026-10-31T20:34:43Z run. If the larger counts refer to a different report or
   corpus selection, that input still needs reconciling before claiming it is
   covered here. Against the freshly pulled `main` baseline, this batch removes
   144 argument and 64 return findings and adds none in either family. The
   earlier 87/7 inventory predates the pulled template-inference commits and
   is not the current inventory.
2. **Disposition every current PSL `A.RETURN.TYPE` location.** The same report
   has 7 engine-only return diagnostics. Inspect the actual return paths and
   contracts, including the apparently missing-return paths. Treat annotations,
   resource wrappers, templates, integer ranges, and string non-emptiness as
   hypotheses to test at each source location, not blanket explanations.
3. **Fix confirmed defects.** For each reproduced defect, add a synthetic
   failing regression and a clean control at PHPStan level 8. Update the rule
   contract and supported-limitations documentation whenever behavior or a
   supported boundary changes. Re-run the pinned differential fixtures and
   compare exact PSL issue sets after each change.
4. **Record boundaries and unsupported cases by evidence.** A supported
   boundary needs a specific source example, a statement of current behavior,
   and documentation in the relevant rule/capability pages. An unsupported
   construct needs a minimal reproducer showing which syntax or semantic
   feature is outside the analyser's supported model. Structured parameter
   binding, dependent bounds, complex callable bodies, and shared integer-range
   normalization are known partial areas; none pre-classifies a listed
   finding.
5. **Refresh all release evidence.** Re-run the required differential packs,
   PSL, Symfony, and Composer comparisons on the release candidate. Verify
   complete file accounting, no execution failures, exact issue snapshots,
   disposition counts, and provenance. For shared CST/analyser changes, compare
   Symfony and Composer issue sets against their pre-change snapshots. Keep the
   gate failed until every difference is dispositioned.

### Current PSL location inventory

The entries below are the exact engine-only source locations in the latest
level-8 report (`2026-10-31T20:34:43Z`). Multiple diagnostics at the same path
and line are marked `×N`. All listed locations remain open until source-backed
dispositions are recorded; the inventory alone is not a disposition.

#### Argument findings (150 diagnostics)

- `packages/crypto/src/Psl/Crypto/Hkdf/extract.php`: 41
- `packages/crypto/src/Psl/Crypto/KeyExchange/generate_key_pair.php`: 25
- `packages/dict/src/Psl/Dict/diff.php`: 27, 28
- `packages/dict/src/Psl/Dict/diff_by_key.php`: 27, 28
- `packages/dict/src/Psl/Dict/drop.php`: 27
- `packages/dict/src/Psl/Dict/intersect.php`: 27, 28
- `packages/dict/src/Psl/Dict/intersect_by_key.php`: 27, 28
- `packages/dict/src/Psl/Dict/merge.php`: 37, 47
- `packages/dict/src/Psl/Dict/sort_by.php`: 54
- `packages/dict/src/Psl/Dict/take.php`: 27
- `packages/dict/src/Psl/Dict/unique.php`: 22
- `packages/dict/src/Psl/Dict/unique_scalar.php`: 30
- `packages/dns/src/Psl/DNS/HTTP/Connector.php`: 105 ×3, 140, 141
- `packages/graph/src/Psl/Graph/add_edge.php`: 43, 47
- `packages/graph/src/Psl/Graph/bfs.php`: 50
- `packages/graph/src/Psl/Graph/dfs.php`: 59
- `packages/graph/src/Psl/Graph/has_path.php`: 49
- `packages/graph/src/Psl/Graph/shortest_path.php`: 40 ×2
- `packages/graph/src/Psl/Graph/shortest_path_by.php`: 66, 103, 124
- `packages/graph/src/Psl/Graph/topological_sort.php`: 39, 50, 71
- `packages/h2/src/Psl/H2/Internal/ConnectionTrait.php`: 114, 269
- `packages/h2/src/Psl/H2/Internal/StateMachine.php`: 453, 578, 648, 750, 857 ×2, 860 ×2, 865, 869, 1042, 1067, 1242, 1265, 1304, 1330, 1338, 1382, 1387, 1468, 1510
- `packages/http-client/src/Psl/HTTP/Client/Connection/Connector.php`: 183, 187 ×2, 218, 222, 236, 283, 294 ×2, 298
- `packages/http-client/src/Psl/HTTP/Client/Connection/PooledConnector.php`: 403, 407 ×2, 433, 437, 451, 505, 516 ×2, 520
- `packages/http-client/src/Psl/HTTP/Client/Internal/H1/Transport.php`: 131
- `packages/http-client/src/Psl/HTTP/Client/Internal/H2/H2Session.php`: 80
- `packages/http-client/src/Psl/HTTP/Client/Internal/HttpTunnel.php`: 58, 62
- `packages/http-client/src/Psl/HTTP/Client/Internal/resolve_url.php`: 71, 80
- `packages/http-client/src/Psl/HTTP/Client/RedirectClient.php`: 223
- `packages/iri/src/Psl/IRI/IRI.php`: 57
- `packages/iri/src/Psl/IRI/from_uri.php`: 28, 35
- `packages/mime/src/Psl/MIME/MultiPart/Parser.php`: 300, 332, 404
- `packages/smtp/src/Psl/SMTP/Client/Transport.php`: 268, 388 ×2
- `packages/smtp/src/Psl/SMTP/Internal/DotStuffingReadHandle.php`: 49, 87
- `packages/socks/src/Psl/Socks/Connector.php`: 45, 46
- `packages/str/src/Psl/Str/splice.php`: 34, 38
- `packages/str/src/Psl/Str/wrap.php`: 65, 75
- `packages/tcp/src/Psl/TCP/Connector.php`: 38
- `packages/tcp/src/Psl/TCP/Socket.php`: 121
- `packages/tcp/src/Psl/TCP/listen.php`: 43
- `packages/terminal/src/Psl/Terminal/Application.php`: 349 ×2
- `packages/terminal/src/Psl/Terminal/Rect.php`: 80 ×4
- `packages/terminal/src/Psl/Terminal/Widget/Block.php`: 105, 106, 107, 108
- `packages/terminal/src/Psl/Terminal/Widget/Internal/LineWrapper.php`: 155
- `packages/tls/src/Psl/TLS/connect.php`: 36
- `packages/tree/src/Psl/Tree/all.php`: 37
- `packages/tree/src/Psl/Tree/any.php`: 37
- `packages/tree/src/Psl/Tree/filter.php`: 47
- `packages/tree/src/Psl/Tree/find.php`: 42
- `packages/tree/src/Psl/Tree/level_order.php`: 37
- `packages/tree/src/Psl/Tree/path_to.php`: 47
- `packages/tree/src/Psl/Tree/to_index.php`: 52
- `packages/type/src/Psl/Type/Internal/MapType.php`: 122, 177
- `packages/type/src/Psl/Type/Internal/MutableMapType.php`: 122, 178
- `packages/type/src/Psl/Type/Internal/MutableSetType.php`: 72, 117
- `packages/type/src/Psl/Type/Internal/MutableVectorType.php`: 99, 138
- `packages/type/src/Psl/Type/Internal/SetType.php`: 72, 114
- `packages/type/src/Psl/Type/Internal/VectorType.php`: 99, 140
- `packages/type/src/Psl/Type/backed_enum.php`: 22
- `packages/type/src/Psl/Type/backed_enum_value.php`: 27
- `packages/type/src/Psl/Type/intersection.php`: 24 ×2
- `packages/type/src/Psl/Type/non_empty_vec.php`: 20
- `packages/type/src/Psl/Type/union.php`: 22 ×2
- `packages/type/src/Psl/Type/unit_enum.php`: 22
- `packages/type/src/Psl/Type/vec.php`: 20
- `packages/unix/src/Psl/Unix/Socket.php`: 98
- `packages/unix/src/Psl/Unix/listen.php`: 28
- `packages/url/src/Psl/URL/from_uri.php`: 43, 46
- `packages/vec/src/Psl/Vec/unique_scalar.php`: 30

#### Return findings (7 diagnostics)

- `packages/date-time/src/Psl/DateTime/DateTimeConvenienceMethodsTrait.php`: 370, 464
- `packages/io/src/Psl/IO/Internal/ResourceHandle.php`: 307, 454
- `packages/iter/src/Psl/Iter/Iterator.php`: 230
- `packages/mime/src/Psl/MIME/ContentDisposition.php`: 199
- `packages/network/src/Psl/Network/Internal/socket_connect.php`: 44
checked against Symfony and Composer issue snapshots below.

For the seven current return findings, the integer-range local-assertion
group is an intentional supported boundary; the sequence-generic and callable
wrapper cases still need reduced fixtures before final disposition:

| PSL location | Classification | Evidence and remaining work |
| --- | --- | --- |
| `DateTimeConvenienceMethodsTrait.php:370`, `:464`; `ContentDisposition.php:199` | Intentional supported boundary | Each casts a string to `int` and relies on a local inline `@var int<...>` refinement. This analyser does not consume inline integer-range `@var` assertions; this exact local-assertion form is outside the current supported refinement model. `TestReturnTypeKeepsInlineReturnAssertionBoundary` demonstrates the unsupported inline assertion and a clean branch-guarded control. |
| `IO/Internal/ResourceHandle.php:307,454` | Open: sequence generic-flow defect vs unsupported construct | These return values through a sequence object populated elsewhere and awaited through `waitFor`. Reproduce property initialization and the operation callback separately; do not count this as covered by deferred generic-property work. |
| `Network/Internal/socket_connect.php:44` | Open: callable-wrapper defect vs unsupported construct | The outer return is a suppression wrapper around a closure with an inline resource return contract, cancellation callbacks, and suspension. Minimize the wrapper and closure contract independently before assigning a disposition. |
| `Iterator.php:230` | Open: integer-range assertion vs state-flow defect | The method combines generator exhaustion, cached nullable state, `count($this->entries)`, and an inline nonnegative-range assertion. Separate the native `count()` result from cached property flow in a fixture before deciding whether to support the inline assertion or the state transfer. |

The 150 remaining argument findings still require a source-level disposition
each. The location inventory is a queue, not evidence of a shared cause. For
each location, record the argument expression, callee contract, actual inferred
type, and whether a minimal positive/negative fixture reproduces a defect or an
explicit unsupported boundary. In particular, do not use structured parameter
binding, dependent bounds, complex callable bodies, or shared integer-range
normalization as a disposition without a reproducer for that finding.

## Existing evidence

### Differential fixtures

The current structured-template slice passes **56/56 level-3** and
**107/107 level-5** cases against PHPStan 2.2.16, with zero engine or
reference mismatches. The checked-in packs and their pinned PHPStan versions are declared in
`testdata/diagnostic-differential*/manifest.json`; `cmd/diagnostic-diff`
compares both diagnostic lists, including duplicates. Earlier recorded runs
covered 508 cases across levels 0–8. Re-run every pinned pack on the release
candidate before accepting this gate.

### PSL source corpus

The measured workload is PSL 6.2.1 at commit
`331f3ab363508825e62c42725f96c091bcceb970`: the `src` directory of each of 73
packages, totaling 1,775 PHP files. Tests, examples, and generated/vendor code
are excluded. The comparison indexes the same 25 `vendor/revolt` dependency
files for Tusk and PHPStan.

The prior recorded dependency-aware PHPStan comparison used PHP 8.4.25 and
PHPStan 2.2.16. The fresh rerun below used PHP 8.5.10 and PHPStan 2.2.16, with
the pinned PSL checkout, PHPStan configuration, first-party manifest, index
manifest, and crosswalk. Both tools analyzed all selected
files. At level 8, the recorded results were:

| F1 | Precision | Recall | Reviewed F1 | Exact | Tusk | PHPStan | Tusk-only | PHPStan-only | Unreviewed PHPStan-only |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 4.83% | 4.37% | 5.39% | 7.41% | 39 | 892 | 724 | 853 | 685 | 563 |

The latest report was generated at `2026-10-31T20:34:43Z` against the pinned
corpus and reference. It analyzed all 1,775 selected files in both tools with
zero parse/read failures. At level 8 it reports 892 Tusk diagnostics, 724
PHPStan diagnostics, 39 exact matches, 853 engine-only diagnostics, 685
PHPStan-only diagnostics, and 563 unreviewed PHPStan-only diagnostics. The
engine-only `A.ARG.TYPE` / `A.RETURN.TYPE` counts are 150 / 7. Their source
locations are listed above. Compared with the freshly pulled `main` baseline,
the candidate has 194 fewer Tusk diagnostics overall, including 144 fewer
`A.ARG.TYPE` and 64 fewer `A.RETURN.TYPE`, with no new locations in either
target family. The overall issue-set diff still needs source-level disposition;
this family-count reduction does not close the release gate. This batch fixes reproduced analyzer defects and documents a supported boundary. The current argument locations have not yet
been individually classified; several return-family reproductions also remain
open. The report has zero unmapped engine diagnostics; mapping is not a
source-level disposition.
Precision is 4.37%, recall is 5.39%, and F1 (compatibility) is 4.83% (reviewed
F1 is 7.41%). These remain far below the gate, and the family findings remain
open.

The fresh setup targets PHP 8.5 in `testdata/phpstan-psl.neon`, loads PSL's
Composer autoloader and the pinned PSL PHPStan extension, and uses the lockfiles
and manifests described in the run artifacts. The analysis does not execute
PSL code. The corpus host lacked `bcmath` and `intl`; Composer installation
used `--ignore-platform-reqs`. Preserve this environment limitation in any
updated evidence unless the run environment changes.

For the recorded PHP 8.4.25 / PHPStan 2.2.16 reference run, the support lock
SHA-256 was
`2a29f882adbd7b973e2742bba1b2a0e1ca6d2efd8f8e25dba2649b5d5ffb7819`, the PSL
`composer.lock` hash was
`85e11c6f9b12cfb0794eb1ac9ab1e33c1f07d1111b9b589c7034660ef11ed111`, and the
configuration hash was
`eca48699f663b6632c6b1ccdf27442d4886dfa52d814a853d118b5f237731d61`. The
first-party and combined index manifest hashes were
`6a955366217506f2cabe420293b4ab19265e5d38e590ba8c7ddc81aa56cc1b34` and
`38c2f5b1cf88c27c33daddabfdb90cc0616cd6a0ba8742eb35bd4f2f64c014ec`.

### Symfony and Composer

The current working candidate was compared with detached baseline `c1ae5a3d`
(the freshly pulled `main`) at level 8 using the same absolute corpus roots.
Symfony analyzed 10,026 files with the same one parse-error fixture in both
snapshots. The candidate removes 151 diagnostics and adds 76 across 123 files;
removals include 77 `A.ARG.TYPE` and 31 `A.RETURN.TYPE`, while additions include
38 `A.ARG.TYPE` and 13 `A.RETURN.TYPE`. Composer analyzed 532 files with no
parse errors; the candidate removes 36 diagnostics and adds 5 across 20 files.
Composer changes include 21 argument and 3 return removals, plus 3 argument,
1 property, and 1 return addition. The additions and all changes not yet tied
to an accepted source-level disposition remain open; these snapshots do not
establish corpus parity. Generated snapshots are in `/tmp/tusk-symfony-main-base.json`,
`/tmp/tusk-symfony-batch.json`, `/tmp/tusk-composer-main-base.json`, and
`/tmp/tusk-composer-batch.json` for this run. Repeat comparisons after each
shared-analyser change and retain the JSON as review evidence.
The earlier callable-semantics comparison remains historical evidence: it had
208 removals and 31 additions on Symfony and 25 removals and 3 additions on
Composer. Those figures must not be substituted for the current candidate
delta. Symfony's parse-error file and the corpus accounting issue remain open.

## Reproduction

Run from the repository root with PHP 8.4 or 8.5. Fetch the pinned test projects
with `go run ./cmd/fetch-test-projects`, install the support and PSL production
dependencies, and select the `src` directory of every PSL package in bytewise
order. Save the complete JSON report and verify its file-manifest,
index-manifest, configuration, and crosswalk hashes. Keep the PHPStan version,
PHP version, corpus commit, command, and tool output with each result.

The level-8 rerun used a temporary PHPStan wrapper to raise PHP's default
128 MB memory limit to 1 GB; without it, PHPStan exhausted memory before
producing JSON. On the case-insensitive macOS host used for the latest recorded
reference run, ordinary PSR-4 fallback autoloading stopped on a duplicate
function declaration. Regenerating the ignored vendor autoloader with
`composer dump-autoload --no-dev --classmap-authoritative --ignore-platform-reqs`
allowed analysis to complete. Its classmap and files-map hashes were
`ad6ce608bab96ccb91042f826506895d112b638feff55f51158c3b715311cfa8` and
`93b810b58a1870f1a08760beba279b0553b747d248ec930a5ad50e378b50e115`. Recheck
these hashes if reproducing that environment; do not assume they apply after
changing the PHP version, filesystem, dependencies, or autoloader generation.

```sh
mapfile -t source_dirs < <(find test_projects/psl/packages -mindepth 2 -maxdepth 2 -type d -name src | LC_ALL=C sort)
paths=$(IFS=,; printf '%s' "${source_dirs[*]#test_projects/psl/}")
(cd testdata/phpstan-psl && composer install --no-interaction --prefer-dist)
(cd test_projects/psl && composer install --no-dev --ignore-platform-reqs --no-interaction --prefer-dist)
go run ./cmd/phpstan-compat \
  --root test_projects/psl \
  --paths "$paths" \
  --index-paths vendor/revolt \
  --levels 0,1,2,3,4,5,6,7,8 \
  --phpstan-bin ../../testdata/phpstan-psl/vendor/bin/phpstan \
  --phpstan-config ../../testdata/phpstan-psl.neon \
  --workers 3 --json --output /tmp/tusk-psl-src-all-levels.json
```

Update this page whenever a corpus, manifest, configuration, threshold,
disposition, or result changes.
