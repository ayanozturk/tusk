# Current implementation plan

Status: in progress

Roadmap source: short-term item 2, “Complete the semantic and type system,” explicitly reprioritized by the user ahead of the active 1.0 stabilization work.

Scope: complete PHPDoc/type semantics, control-flow and variable/property flow, interprocedural summaries with precise invalidation, and testable semantic extension contracts. Preserve parser losslessness, deterministic diagnostics, public API compatibility, and PHP Strom integration. The 1.0 release remains blocked and deferred; its acceptance evidence and outstanding gates remain in [docs/releasing-1.0.md](docs/releasing-1.0.md).

## Work sequence

1. [ ] Complete PHPDoc and generic type semantics: templates, bounds, variance, generic inheritance, aliases/imports, conditional and indexed types, key/value projections, callable signatures, shapes, non-empty and literal types, and integer ranges. Initial slices share function/method template inference nested in supported generic parameters from known argument generic contexts, while retaining declared bounds.
2. [ ] Extend CFG and variable/property flow through loops, exceptions, generators, closures, references, property initialization, readonly constraints, and dynamic properties. Keep state joins deterministic and add behavioral regression tests for each supported path.
3. [ ] Add interprocedural summaries and invalidate dependants only when exported semantic facts change. Cover cache identity, dependency invalidation, and concurrent deterministic output.
4. [ ] Define testable extension contracts for dynamic returns, assertions, generated members, reflection, service containers, and ORMs. Keep extension behavior isolated, documented, and independently testable.
5. [ ] Validate the completed semantic slices with focused unit and differential fixtures, relevant reviewed corpora, snapshot/incremental tests, and pinned plus sibling-development PHP Strom tests. Update capability documentation and release-plan evidence to reflect actual support.

## First milestone acceptance

- Structured PHPDoc generic parameters infer templates only from matching, known argument structure; unresolved or incompatible shapes remain conservative and respect declared bounds.
- Clean and mismatching argument/return fixtures cover the supported forms and compare against the reference analyser at the exact rule level.
- Existing generic, callable, shape, alias, conditional-return, differential, and snapshot tests remain green.
- Any newly supported diagnostic behavior has a regression fixture and matching rule/capability documentation.
- No public API change is made without compatibility review and pinned/sibling-development PHP Strom validation.

### Verified structured-template slices

- [x] Share declared-template inference between functions and methods for direct,
  callable-return, and matching known generic arguments, including nested generic
  arguments and repeated variadic contributions. Preserve legacy class-string
  forms and conservative method callback-signature checking.
- [x] Substitute method/constructor argument templates without mutating indexed
  signatures. Invalid inferred bindings retain their declared return bound.
  Expand forwarded local template bounds through nullable/union/intersection
  argument types.
- [x] Compare complete level-3 and level-5 packs against the restored pinned
  PHPStan 2.2.16 executable: 56/56 and 107/107 cases respectively, with zero
  engine or reference expectation mismatches.
- [x] Run `go test ./...`, `go vet ./...`, focused semantic/template/callable race
  tests, rule-inventory/style documentation checks, and PHP Strom sibling
  server tests/vet.
- [x] Compare production-registry corpus snapshots against the preceding
  structured-function commit: Composer has zero issue-set differences across
  532 files. Symfony has 10,026 files (one unchanged parse-error file), with
  nine changed files, 12 removed diagnostic locations, five refined messages
  at existing locations, and no new diagnostic locations. These are intentional
  semantic changes, not a claim of exact Symfony issue-set parity.

These slices do not complete the first milestone or the broader roadmap item.
Next, continue PHPDoc/type semantics with explicit regression/reference evidence
for unsupported union structures, dependent bounds, and generic callback
contracts before moving to flow or interprocedural work. Arguments without
recoverable generic context and static-method argument checking remain partial.

## Deferred release work

- The full-analyser 1.0 release is not authorized by this reprioritization. Continue its separate correctness, reviewed-corpus quality, performance/resource, compatibility, packaging, and publication gates only after this roadmap item is completed or the user changes priority again.
- Do not claim completion based on unit tests alone: keep full corpus accounting, differential parity, race safety, and reproducible performance evidence as release gates.
