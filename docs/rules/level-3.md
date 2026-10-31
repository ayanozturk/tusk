# PHPStan-compatible rules: level 3

<!-- rule-inventory: level=3 introduced=5 cumulative=26 -->

[Back to the README static-analysis section](../../README.md#static-analysis) · [Open the analyser capability matrix](../analyser-capability-matrix.md)

## Rule inventory

- **Introduced at this level:** 5 registered levelled rules: `A.PROP.TYPE`, `A.RETURN.NEVER`, `A.RETURN.TYPE`, `A.RETURN.VOID`, and `Level3.ThrowType`.
- **Cumulative registered levelled rules:** 26. The assignment, binary-operation, void-purity, method, and PHPDoc rules from level 2 are cumulative and are not introduced here.
- **Checked-in differential pack:** 56 cases in `testdata/diagnostic-differential-level3`.

## Registered rules, one by one

### `A.PROP.TYPE`

**What it checks:** Checks assignments to typed properties against the property’s declared type.

**Why it helps:** Typed properties are runtime-enforced contracts; incompatible assignments fail when executed.

**Example that reports:**

```php
<?php
class Cart { public int $count; }
$cart = new Cart();
$cart->count = "many";
```

**A safer form:**

```php
<?php
$cart->count = 3;
```

### `A.RETURN.NEVER`

**What it checks:** Checks selected never-returning functions and methods for paths that return normally instead of terminating or throwing.

**Why it helps:** A never contract tells callers execution cannot continue; violating it makes control-flow reasoning unsafe.

**Example that reports:**

```php
<?php
function stop(): never { return; }
```

**A safer form:**

```php
<?php
function stop(): never { throw new RuntimeException(); }
```

### `A.RETURN.TYPE`

**What it checks:** Compares returned expressions with declared return types for supported inferred expressions.

**Why it helps:** Return mismatches violate the function contract and can cause downstream type errors.

**Example that reports:**

```php
<?php
function countItems(): int { return "three"; }
```

**A safer form:**

```php
<?php
function countItems(): int { return 3; }
```

Named generic function and method returns use call-site bindings from direct template
arguments and declared or supported inferred callable returns, including imported function aliases.
The level-3 clean/mismatch fixtures cover `preserve($value)` and a `Value`
template nested in a known `list<Value>` or array parameter structure. Repeated
direct arguments join their inferred types; unresolved or structurally unknown
bindings use a bound or `mixed`. Supported unannotated arrows and simple closures infer returns in isolated
parameter/capture scopes. Method template binding preserves arbitrary declared
template names and inherited callback metadata, and shares structured generic
inference with functions. Clean/mismatch fixtures cover inherited, interface,
named, and static method returns from known iterable generic arguments; unit
tests also cover nested generic arguments and variadic joins. A binding outside
its declared bound retains the bound for return inference, while the call still
reports the invalid argument. Unsupported generic structures and complex callback
control flow remain partial.

Bare class-template property annotations such as `@var T` now retain the
owning class template when compared with a method's `@return T`; the native
property type remains the fallback when no documented template refines it.
`TestReturnTypePreservesClassTemplatePropertyReturn` covers a
`Throwable` property refined to `T of Throwable` and a clean incompatible-return
control. Compound property annotations such as `T|null` remain conservative
because their correctness can depend on assignment and state correlation. This
does not bind a concrete `T` for every unparameterized object or prove dependent
generic flow across deferred state.

Generic PHPDoc on constructor-promoted properties is retained when a method
calls through that property. `TestReturnTypeSubstitutesPromotedGenericInterfaceProperty`
covers `TypeInterface<O>` delegated returns and an incompatible-return control;
unparameterized receivers and deferred generic state remain outside this case.

Compatible generic PHPDoc on an ordinary property also refines its native class
type when a return is delegated through that property. The reduced
`Deferred<T> -> Awaitable<T> -> await()` case in
`TestReturnTypeSubstitutesGenericPropertyThroughAwaitableChain` checks the
resolved tuple return and an incompatible-return control. For nested method
templates, a method accepting `iterable<Awaitable<T>>` can bind `T` from the
value position of a declared `iterable<TKey, Awaitable<TValue>>` parameter;
`TestReturnTypeBindsNestedStaticMethodTemplatesForAwaitableIterate` covers the
two-argument generator contract and a wrong-return control. This remains a
narrow iterable value-position case; structured parameter binding and general
dependent template bounds remain partial.

Local `@var` assertions preserve a class template in a simple union such as
`T|null` inside a recursive generic method. The regression
`TestReturnTypePreservesTemplateInRecursiveLocalVarAssertion` covers the
recursive return and an incompatible-return control. Inline integer-range
assertions such as `/** @var int<0, max> $value */` are not currently applied to
local variables; use an explicit runtime guard when a return depends on that
range, as covered by `TestReturnTypeKeepsInlineReturnAssertionBoundary`.

An empty-array strict comparison also refines the surviving path: after
`if ([] === $result) { throw ...; }`, the fallthrough type is
`non-empty-array`. `TestReturnTypeNarrowsArrayAfterEmptyGuard` covers a built
array with a throwing empty control and a clean negative control that returns
`[]`. Loop-produced array element/key shapes remain outside this refinement.

Anonymous `@var` annotations immediately before a return, or named assertions
matching the returned variable, refine its observed type after CST lowering.
Integer-range `@var` assertions at a return are currently an explicit unsupported
case: range-bound normalization between inline assertion and PHPDoc declaration
is incomplete. The PSL `getYearShort`, `getDayOfYear`, and `ContentDisposition::size`
findings use this form. `int<min, max>` refinements from signed-literal branch
checks are supported, including a negative lower bound; the regression
`TestReturnTypeIntersectsConjunctiveIntegerBounds` has a bounded clean return and
an out-of-range control.
Declared scalar template bounds are honored; an
assertion naming a different variable does not override the return. Integer-range
assertions remain conservative pending shared declaration/expression normalization.

Return-completeness analysis treats a literal `while (true)`, `do { ... }
while (true)`, or `for (;;)` without a possible `break` as non-fallthrough.
For example, `function read(): int { while (true) { if (ready()) return 1; } }`
does not fall through; adding a reachable `break` makes a missing return
reportable. A try/catch inside a non-breakable literal infinite loop also counts
as non-fallthrough when the successful path returns and caught failures
continue; `TestReturnTypeRecognizesRetryLoopReturnsThroughCatch` pairs that
case with a nonliteral-loop missing-return control. Break detection stays
conservative around nested or unsupported statement forms, and this is not a
general loop termination proof.

For direct enum-case returns, an exact PHPDoc case contract such as
`@return DateStyle::Medium` accepts only that case. The enum's broad declared
type remains the native return type; a different case is still reported.

### `A.RETURN.VOID`

**What it checks:** Reports returned values from functions or methods declared void.

**Why it helps:** A void declaration promises no value; returning one makes the API contract inconsistent.

**Example that reports:**

```php
<?php
function logMessage(): void { return "done"; }
```

**A safer form:**

```php
<?php
function logMessage(): void { echo "done"; return; }
```

### `Level3.ThrowType`

**What it checks:** Checks selected throw expressions to ensure the thrown value is a Throwable.

**Why it helps:** Throwing a non-exception value is invalid PHP and fails at runtime.

**Example that reports:**

```php
<?php
throw "failed";
```

**A safer form:**

```php
<?php
throw new RuntimeException("failed");
```

## Coverage and boundaries

Level 3 checks declared return types, values returned from `void` functions and methods, `never` return contracts, typed-property assignments, and selected throw targets. Return inference covers arithmetic, comparison, logical, unary-not, unary-numeric, and spaceship results in addition to the previously covered literal, call, property, cast, conditional, coalesce, and match expressions. Match branches preserve integer and string literal return values, and `match (true)` arms narrow supported `instanceof` conditions. Resolved `never`-returning calls terminate branches for return completeness and variable narrowing. Constant non-empty prefixes in either concatenation form, including binary `.` expressions carried through local assignments, non-empty backed-enum values, and casts from numeric values preserve `non-empty-string` when inferred; regressions cover the direct URL-like return and an empty control. Boolean literals retain `true`/`false` precision, and the built-in `fopen` signature is `resource|false`. Generic `$this` return inference prefers project-index `ResolveMethod` over unbound same-class AST signatures. Method-level `@template T` binds from callable arguments (`callable(): T` and unions such as `callable(...): T|CallbackInterface<T>`) using declared or supported inferred closure/arrow return types rather than inferred `Closure`; enclosing class templates also retain their identity in local return checks. Typed receivers without generic args use class `@extends` generic parents for method return inference (for example `$repo->find()` on a child repository, including a multi-hop `ServiceEntityRepository` parent). Unbound class templates inside generic arguments such as `EntityRepository<T>` stay as `T` when the current class still owns that template. Normally completing `finally` blocks preserve the selected try/catch outcome for return completeness. Return and property checks use immutable semantic facts and known class/property metadata where available. Return checks also consume compatible narrowed PHPDoc types on explicit and constructor-promoted properties (for example, a promoted `positive-int` property returned from a `positive-int` method), and retain nested class-template arguments on promoted property receivers when resolving delegated returns.

Ternary inference preserves narrowed literal branches, compatible adjacent integer intervals can cover a wider declared range, and `never` contributes no value to a return union. Coverage remains partial for PHPDoc types, dynamic expressions, complex unions, every assignment operator, and the complete PHPStan type system. Levels 0 through 2 are cumulative.
