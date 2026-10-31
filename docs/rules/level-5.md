# PHPStan-compatible rules: level 5

<!-- rule-inventory: level=5 introduced=1 cumulative=28 -->

[Back to the README static-analysis section](../../README.md#static-analysis) · [Open the analyser capability matrix](../analyser-capability-matrix.md)

## Rule inventory

- **Introduced at this level:** 1 registered levelled rule, `A.ARG.TYPE`.
- **Cumulative registered levelled rules:** 28.
- **Checked-in differential pack:** 107 cases in `testdata/diagnostic-differential-level5`.

## Registered rules, one by one

### `A.ARG.TYPE`

**What it checks:** Checks selected calls to known functions, methods, and constructors against parameter types, including supported PHPDoc and inferred types.

**Why it helps:** Passing the wrong value can cause runtime type errors or incorrect library behavior; checking the call site localizes the correction.

**Example that reports:**

```php
<?php
function save(int $id): void {}
save("abc");
```

**A safer form:**

```php
<?php
function save(int $id): void {}
save(42);
```

## Coverage and boundaries

`A.ARG.TYPE` reports incompatible arguments passed to resolved named functions, methods, and constructors. The 107-case pack covers positional, named, method, constructor, `self`-parameter, PHPDoc alias/list, generic function parameter bounds, structured generic function/method/constructor parameters (including clean and bound-mismatch cases), unbound class/method-template parameters and bounds, closure-filter, Symfony form and dependency-injection return, Laravel Eloquent builder, string-narrowing, scalar/nested array-shape offsets, indexed PHPDoc shape returns, generic array value offsets, `key-of`/`value-of` shape and generic array projections, callable parameter and return signatures, nested PHPDoc shape offsets, optional offsets, and union-shape projections, native/PHPDoc refinement (including refinements flowing through function parameters that also declare native types), parameter-dependent function and method return branches, strict exact-string equality branches, and valid literal branch joins; failing cases match PHPStan's `argument.type` identifier at level 5. Failing cases remain silent at level 4. Nullable-only argument mismatches are gated at level 8, while an extra `false` union arm is gated at level 7, matching the reference analyser's thresholds. Method `self`/`static`/`parent` parameter types bind to the callee's declaring class, called class, and parent rather than the caller. Direct `list<T>` parameters and class-level `@phpstan-type` / `@psalm-type` list aliases accept compatible arrays rather than being resolved as namespaced classes. Unbound method/class templates with bounds expand to those bounds instead of fake classes unless a call-site argument binds them. Interface method `@template T` stays a template instead of a namespaced fake class. Property PHPDoc generics such as `InputBag<string>` refine native property types, so `query->get('search', '')` infers `string`. Arrow functions and closures are `Closure`, and `callable` accepts `Closure`. A local `@var` on the next expression statement asserts that variable's type, including after a nullable `getUser()` assignment. Conditional PHPDoc returns use a branch selected from a statically known argument; unresolved argument types retain the branch union. Conditional annotations continue to fall back to the native return when their tested class types are unresolved. A terminating `if (!is_string($a) || !is_string($b)) { return; }` narrows both variables to `string` afterwards. `is_string($x) && $x !== '' ? (int)$x : null` is a ternary of `int|null`, not a boolean `&&`. Symfony container `get(class-string<T>)` and Doctrine repository `find(class-string<T>)` bind the service/entity type and preserve nullable return behavior. Scalar and nested shape offsets preserve declared builtin types, including unions of quoted PHPDoc string values. Unbound template parameters are treated as mixed (or their declared bound) when no concrete type argument is present, avoiding a fabricated class mismatch. A truthy `$x?->prop` (or `$x->prop`) in an `if` narrows `$x` to non-null inside that branch. A non-terminating `if (!$x) { $x = new T; }` or `if ($x === false) { $x = 0; }` joins the assignment with the implicit else so later statements see the filled-in type. A final backed enum excluded by `instanceof` narrows the false path; after a non-terminating `if ($value instanceof E) { $value = $value->value; }`, the assignment joins with the implicit else and later integer arguments remain clean. `$this->prop instanceof T` and `$obj->prop instanceof T` narrow that property in the then-branch so later arguments see `T`. A `while ($x !== null)` loop header narrows `$x` to non-null inside the loop body. Static methods, dynamic callables, unpacked arguments, extension-provided signatures, `@var` on statements other than expressions and returns, conditions beyond the supported type tests, and the full PHPStan type lattice remain partial. Symfony's richer conditional service-id return expression also remains partial. The rule uses the existing shared argument-call traversal, so moving it from level 10 does not add a pass to the default analysis path.

After an `if (is_string($value))` branch replaces a `string|T` value, the
fallthrough type excludes `string`; an unguarded `string|T` argument still
reports. A `@var` assertion attached to a return statement is applied while
checking nested calls in that return expression. Ternary branches retain
integer literal values for bounded argument checks, so `flag ? 4 : 0` satisfies
`int<0, 255>` while `flag ? 256 : 0` does not. The level-8 regression cases
`TestArgumentTypeNarrowsFalseIsStringPathAfterReplacement`,
`TestArgumentTypeAppliesReturnVarDocBeforeNestedCalls`, and
`TestArgumentTypePreservesLiteralRangeAcrossTernary` cover these paths with
negative controls.

Named generic functions and methods bind declared direct template parameters from positional,
named, and variadic arguments, joining repeated contributions. Templates nested
in supported PHPDoc generic iterable parameters (including nested generic
arguments) bind when the argument has matching, known generic context; unresolved
or incompatible shapes remain at their declared bound or `mixed`. A parameter of
`callable(): Result` also binds `Result` from a declared or supported inferred closure/arrow return or
known callable-variable return. The substituted return feeds argument
checks, so `acceptInt(preserve('wrong'))` reports while preserving an integer is
clean. Template names need not start with `T`. Unresolved templates fall back to
their declared bound or `mixed`. Structured parameters on resolved instance
methods and constructors use the same bound-aware argument substitution as
functions, without mutating indexed signatures. An out-of-bound binding also
retains the declared bound in the inferred return, so both the invalid argument
and an incompatible use of that return can report. Static-method argument
checking, unpacked arguments, dependent template bounds, and complex callback
bodies remain partial.

Forwarded local template arguments expand their declared bounds even when they
occur in a nullable or union type. Method callback signatures containing open
templates retain their existing conservative gate; structured argument inference
does not enable full generic callback-signature validation.

False-path facts from the left operand of `||` refine the right operand while
checking a terminating compound guard. Integer range bound keywords such as
`min` and `max` remain global inside namespaces.

For non-terminating `if` statements, branch assignments are joined with the
implicit fallthrough. A guarded parser conversion such as
`if ($value !== null && !$value instanceof Constraint) { $value = $parser->parse($value); }`
leaves the following code with `Constraint|null`; the regression also checks
that passing an unrelated `int` to a `Constraint|null` parameter still fails.
The false path of a strict `!== null` guard is `null`, while its true path is
non-null. This behavior is covered by
`TestNegatedInstanceofAssignmentNarrowsUnionAcrossJoin`.

Known standalone `callable(...)` / `Closure(...)` parameter contracts check
callback parameter contravariance, return covariance, and required, optional,
and variadic arity at level 5. Relative callback types bind to the callee, and
inherited generic signatures retain substituted callback metadata. Supported
unannotated arrows and simple closure bodies use parameter scopes, implicit
arrow captures, explicit closure captures, and isolated local assignments.
Declared return types remain authoritative. Loops, try/switch bodies,
generators, reference parameters, unresolved templates, and nesting beyond 16
remain conservative; nullable or union callback contracts and the complete
PHPStan callable/array type lattice are not signature-checked.

In class methods, a value whose inferred type is a declared class template uses
that template's declared bound for argument compatibility. For example, a
property documented as `Te` with `@template Te of Throwable` may be passed to a
parameter declared as `Throwable`; the level-8 regression
`TestReturnTypePreservesClassTemplatePropertyReturn` covers this and an
incompatible return control. This does not infer a concrete type argument for
every generic object construction.

A local assertion that names a class template also retains its declared bound
when the normalized name is namespace-qualified. For example, after an
`int|string` guard, `/** @var Tk $key */` for `@template Tk of array-key` may
flow to a method that accepts `int|string`.
`TestArgumentTypeUsesClassTemplateBoundAfterLocalAssertion` covers this clean
path and an incompatible object control. Other dependent template bounds remain
partial.

Argument checks also refine an array after a strict non-empty guard such as
`if ([] !== $errors)`, and recognize a literal array as non-empty when it has
an explicit element before any unpacked elements. This accepts
`new CompositeException([$first, ...$rest])` when the constructor requires
`non-empty-array`, while `new CompositeException([])` still reports a mismatch.
`TestArgumentTypeNarrowsNonEmptyArrayGuardsAndSpread` covers both accepted
forms and the empty-array control. The rule does not infer non-emptiness from
unpacked elements alone.

The positive branch of `is_resource($value)` narrows a known resource/object/null
union to `resource`. A negative guard such as `if (!is_resource($value)) throw`
therefore permits the surviving value to be passed to resource-only APIs.
`TestResourcePredicateNarrowsUnionForArguments` covers that path and an
unguarded union control.
