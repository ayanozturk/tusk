package analyse

import (
	"strings"
	"testing"
)

func TestReturnTypeNarrowsAfterNeverCall(t *testing.T) {
	const source = `<?php
namespace Psl;
function invariant_violation(string $message): never { throw new \RuntimeException($message); }
namespace Psl\IO\Internal;
/** @return resource|false */
function openResource(): mixed { return fopen('php://temp', 'r+b'); }
/** @return resource */
function useResource(): mixed {
    $resource = openResource();
    if ($resource === false) { \Psl\invariant_violation('failed'); }
    return $resource;
}
/** @return resource */
function invalidControl(): mixed { return openResource(); }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the unchecked false return to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeSubstitutesGenericPropertyThroughAwaitableChain(t *testing.T) {
	const source = `<?php
namespace Psl\Async;
/** @template-covariant T */
class Awaitable {
    /** @var T */
    private mixed $value;

    /** @return T */
    public function await(): mixed { return $this->value; }
}
/** @template T */
class Deferred {
    /** @var Awaitable<T> */
    private Awaitable $awaitable;

    /** @return Awaitable<T> */
    public function getAwaitable(): Awaitable { return $this->awaitable; }
}
namespace Psl\HTTP;
use Psl\Async\Deferred;
class H2Stream {
    /** @var Deferred<array{int<100, 999>, string, bool}> */
    private Deferred $responseDeferred;

    /** @return array{int<100, 999>, string, bool} */
    public function awaitResponse(): array {
        return $this->responseDeferred->getAwaitable()->await();
    }
}
class InvalidStream {
    /** @var Deferred<int> */
    private Deferred $responseDeferred;

    /** @return string */
    public function awaitResponse(): mixed {
        return $this->responseDeferred->getAwaitable()->await();
    }
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"awaitable-property.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "actual: [int]") {
		t.Fatalf("expected only the incompatible delegated return to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeKeepsInlineReturnAssertionBoundary(t *testing.T) {
	const source = `<?php
/** @return int<0, 5> */
function inlineAssertion(): int {
    /** @var int<0, 5> */
    return (int) '3';
}
/** @return int<0, 5> */
function checkedControl(int $value): int {
    if ($value < 0 || $value > 5) {
        throw new \\InvalidArgumentException();
    }
    return $value;
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"inline-return-assertion.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "inlineAssertion") {
		t.Fatalf("expected the inline assertion boundary and a clean guarded control, got %#v", returnIssues)
	}
}

func TestReturnTypePreservesMatchLiteralBranches(t *testing.T) {
	const source = `<?php
enum Month { case January; case February; case April; }
/** @return int<28, 31> */
function monthDays(Month $month): int {
    return match ($month) {
        Month::January => 31,
        Month::February => 28,
        Month::April => 30,
    };
}
/** @return 1|2 */
function ipFamily(bool $ipv4): int {
    return match ($ipv4) { true => 1, false => 2 };
}
/** @return int<28, 31> */
function invalidMonthDays(): int { return match (true) { true => 32, false => 28 }; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"match-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidMonthDays") {
		t.Fatalf("expected only the out-of-range match control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeInfersNonEmptyConcatenations(t *testing.T) {
	const source = `<?php
/** @return non-empty-string */
function prefix(string $value): string { return 'prefix:' . $value; }
/** @return non-empty-string */
function uri(string $scheme, string $host): string { return "tcp://{$host}:" . $scheme; }
/** @return non-empty-string */
function invalidControl(string $value): string { return $value; }
enum Scheme: string { case TCP = 'tcp'; case Unix = 'unix'; }
/** @return non-empty-string */
function enumUri(Scheme $scheme, string $host): string { return "{$scheme->value}://{$host}"; }
enum MaybeEmptyScheme: string { case Empty = ''; case TCP = 'tcp'; }
/** @return non-empty-string */
function invalidEnumValue(MaybeEmptyScheme $scheme): string { return $scheme->value; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"non-empty-concat.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 2 || !strings.Contains(returnIssues[0].Message, "invalidControl") || !strings.Contains(returnIssues[1].Message, "invalidEnumValue") {
		t.Fatalf("expected only the possibly-empty variable control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeNarrowsMatchTrueInstanceofArms(t *testing.T) {
	const source = `<?php
class AddressList {}
class Mailbox {}
function parseAddress(string $value): AddressList { return new AddressList(); }
function fromMailbox(Mailbox $value): AddressList { return new AddressList(); }
function addressList(string|Mailbox|AddressList $value): AddressList {
    return match (true) {
        $value instanceof AddressList => $value,
        $value instanceof Mailbox => fromMailbox($value),
        default => parseAddress($value),
    };
}

function invalidControl(string|Mailbox|AddressList $value): Mailbox { return match (true) {
    $value instanceof AddressList => $value,
    $value instanceof Mailbox => $value,
    default => $value,
}; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"match-true-narrowing.php": source}, 8)
	var typeIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			typeIssues = append(typeIssues, issue)
		}
	}
	if len(typeIssues) != 1 || !strings.Contains(typeIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the invalid match return control to fail, got %#v", typeIssues)
	}
}
func TestReturnTypeTracksNonEmptyCompoundStrings(t *testing.T) {
	const source = `<?php
/**
 * @param int<0, 3> $index
 * @return non-empty-string
 */
function indexedCharacter(int $index): string {
    $alphabet = '0123';
    return $alphabet[$index];
}
/** @return non-empty-string */
function appendedSuffix(string $value): string {
    $result = $value;
    $result .= ':suffix';
    return $result;
}
/** @return non-empty-string */
function invalidEmptyControl(): string { return ''; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"non-empty-compound.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidEmptyControl") {
		t.Fatalf("expected only the empty-string control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypePropagatesStraightLineDoWhileAssignments(t *testing.T) {
	const source = `<?php
/**
 * @param int<0, max> $number
 * @param int<2, 36> $base
 * @return non-empty-string
 */
function encoded(int $number, int $base): string {
    $alphabet = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $result = '';
    do {
        $quotient = intdiv($number, $base);
        /** @var int<0, 61> $index */
        $index = $number - ($quotient * $base);
        $result = $alphabet[$index] . $result;
        $number = $quotient;
    } while (0 !== $number);

    return $result;
}
/** @return non-empty-string */
function invalidControl(): string { return ''; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"do-while-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the empty-string control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeRecognizesInfiniteLoopReturns(t *testing.T) {
	const source = `<?php
/** @return int */
function readUntilReady(bool $ready): int {
    while (true) {
        if ($ready) { return 1; }
    }
}
/** @return int */
function invalidBreakControl(bool $break): int {
    while (true) {
        if ($break) { break; }
    }
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"infinite-loop-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidBreakControl") {
		t.Fatalf("expected only the breakable infinite-loop control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeRecognizesRetryLoopReturnsThroughCatch(t *testing.T) {
	const source = `<?php
function retryForever(): string {
    while (true) {
        try { return 'ready'; }
        catch (RuntimeException) { continue; }
    }
}
function retryUntilReady(bool $ready): string {
    while ($ready) {
        try { return 'ready'; }
        catch (RuntimeException) { continue; }
    }
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"retry-loop-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "retryUntilReady") {
		t.Fatalf("expected try-return retry loop clean and conditional loop control failing, got %#v", returnIssues)
	}
}

func TestReturnTypePreservesEnumCaseReturnPHPDoc(t *testing.T) {
	const source = `<?php

enum DateStyle { case Short; case Medium; }
/** @return DateStyle::Medium */
function defaultStyle(): DateStyle { return DateStyle::Medium; }
/** @return DateStyle::Medium */
function invalidCaseControl(): DateStyle { return DateStyle::Short; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"enum-case-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidCaseControl") {
		t.Fatalf("expected only the wrong enum case to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeNarrowsIsNumericStringBranch(t *testing.T) {
	const source = `<?php
/** @return numeric-string */
function numericString(mixed $value): string {
    if (is_string($value) && is_numeric($value)) { return $value; }
    throw new InvalidArgumentException();
}
/** @return numeric-string */
function invalidControl(): string { return 'plain text'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"numeric-string-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the non-numeric return control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeInfersNumericStringCasts(t *testing.T) {
	const source = `<?php
/** @return numeric-string */
function castInteger(int|float $number): string { return (string) $number; }
/** @return numeric-string */
function invalidControl(string $value): string { return $value; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"numeric-string-cast.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the unproven string control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeInfersNumericStringCastAfterIsNumeric(t *testing.T) {
	const source = `<?php
/** @return numeric-string */
function castNumeric(mixed $value): string {
    if (is_numeric($value)) { return (string) $value; }
    throw new InvalidArgumentException();
}
/** @return numeric-string */
function invalidControl(mixed $value): string { return (string) $value; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"numeric-string-predicate-cast.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the unguarded cast control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeTracksNonEmptyCallConcatenationThroughVariable(t *testing.T) {
	const source = `<?php
function maybeEmpty(): string { return ''; }
/** @return non-empty-string */
function generated(): string {
    $result = maybeEmpty() . '://';
    return $result;
}
/** @return non-empty-string */
function invalidControl(): string {
    $result = maybeEmpty();
    return $result;
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"non-empty-call-concat-variable.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the possibly-empty call return to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeNarrowsArrayAfterEmptyGuard(t *testing.T) {
	const source = `<?php
/** @return non-empty-array<int, string> */
function oneOrMore(array $items): array {
    $result = [];
    foreach ($items as $item) { $result[] = $item; }
    if ([] === $result) { throw new \\InvalidArgumentException(); }
    return $result;
}
/** @return non-empty-array<int, string> */
function invalidControl(): array { return []; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"non-empty-array-guard.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected the guarded result to pass and the empty control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeIntersectsConjunctiveIntegerBounds(t *testing.T) {
	const source = `<?php
/** @return int<-128, 127> */
function narrowed(int $value): int {
    if ($value < -128 || $value > 127) { throw new \\InvalidArgumentException(); }
    return $value;
}
/** @return int<-128, 127> */
function invalidControl(): int { return 128; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"integer-range-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected bounds-checked value to pass and out-of-range control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeAcceptsDeclaredRecordOrNullFromLoop(t *testing.T) {
	const source = `<?php
class NSEC3Record {}
/** @param list<NSEC3Record> $records @return NSEC3Record|null */
function findCovering(array $records, bool $match): ?NSEC3Record {
    foreach ($records as $record) {
        if ($match) { return $record; }
    }
    return null;
}
/** @return NSEC3Record|null */
function invalidControl(): ?NSEC3Record { return 'wrong'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"record-null-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected record-or-null loop to pass and invalid control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypePreservesNonEmptyPropertyPrefixedString(t *testing.T) {
	const source = `<?php
class Authority { public function toString(): string { return ''; } }
class URL {
    public string $scheme;
    public string $path = '';
    public Authority $authority;
    /** @return non-empty-string */
    public function toString(): string {
        $result = $this->scheme . '://' . $this->authority->toString();
        $result .= $this->path;
        return $result;
    }
}
/** @return non-empty-string */
function invalidControl(): string { return ''; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"url-non-empty-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected literal protocol prefix to prove non-empty return and empty control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeResolvesImportedGenericArrayElements(t *testing.T) {
	const source = `<?php
namespace Psl\DNS\Record { class NSEC3Record {} }
namespace Psl\DNSSEC\Internal\NSEC {
use Psl\DNS\Record\NSEC3Record;
/** @param list<NSEC3Record> $records */
function findCovering(array $records, bool $match): null|NSEC3Record {
    foreach ($records as $record) {
        if ($match) { return $record; }
    }
    return null;
}

/** @return null|NSEC3Record */
function invalidControl(): null|NSEC3Record { return 'wrong'; }
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"imported-record-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected imported record-or-null loop to pass and invalid control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypePreservesClassTemplatePropertyReturn(t *testing.T) {
	const source = `<?php
class AlphaError extends RuntimeException {}
class BetaError extends RuntimeException {}
/** @template-covariant Te of Throwable */
class Failure {
    /** @var Te */
    private Throwable $throwable;
    /** @param Te $throwable */
    public function __construct(Throwable $throwable) { $this->throwable = $throwable; }
    /** @return Te */
    public function getThrowable(): Throwable { return $this->throwable; }
    /** @template Ts @param Closure(): Ts $success @return Failure<Ts, Te> */
    public function map(Closure $success): Failure { return new Failure($this->throwable); }
}
/** @template I @template O */
interface TypeInterface {
    /** @param I $value @return O */
    public function assert(mixed $value): mixed;
}
/** @template I @template O */
class ConvertedType {
    /** @param TypeInterface<O> $into */
    public function __construct(private TypeInterface $into) {}
    /** @return O */
    public function assert(mixed $value): mixed { return $this->into->assert($value); }
}
/** @return AlphaError */
function invalidFailure(BetaError $error): Throwable { return $error; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"generic-class-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidFailure") {
		t.Fatalf("expected generic property return and incompatible control result, got %#v", returnIssues)
	}
}

func TestArgumentTypeNarrowsNonEmptyArrayGuardsAndSpread(t *testing.T) {
	const source = `<?php
class CompositeException {
    /** @param non-empty-array<Throwable> $errors */
    public function __construct(array $errors) {}
}

function guarded(array $errors): never {
    if ([] !== $errors) {
        if (1 === count($errors)) { throw $errors[0]; }
        throw new CompositeException($errors);
    }
    throw new RuntimeException();
}
function spread(Throwable $exception, array $errors): never {
    if ([] === $errors) { throw $exception; }
    throw new CompositeException([$exception, ...$errors]);
}
function invalidControl(): void { new CompositeException([]); }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"non-empty-array-args.php": source}, 8)
	var argIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.ARG.TYPE" {
			argIssues = append(argIssues, issue)
		}
	}
	if len(argIssues) != 1 || argIssues[0].Line != 18 {
		t.Fatalf("expected strict-empty guard and non-empty spread calls to pass with empty control failing, got %#v", argIssues)
	}
}

func TestDateIntervalDateStringHasVersionCorrectObjectReturn(t *testing.T) {
	const source = `<?php
use DateInterval;
function duration(int $seconds): DateInterval {
    return DateInterval::createFromDateString($seconds . ' seconds');
}

function period(int $years, int $months, int $days): DateInterval {
    $parts = [];
    if (0 !== $years) { $parts[] = $years . ' years'; }
    if (0 !== $months) { $parts[] = $months . ' months'; }
    if (0 !== $days) { $parts[] = $days . ' days'; }
    if ([] === $parts) { $parts[] = '0 days'; }
    return DateInterval::createFromDateString(implode(' ', $parts));
}
function invalidControl(): DateInterval { return new stdClass(); }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"date-interval-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected valid generated DateInterval strings to pass and wrong object control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeBindsNestedStaticMethodTemplatesForAwaitableIterate(t *testing.T) {
	const awaitableSource = `<?php
namespace Psl\Async;
use Generator;
/** @template-covariant T */
class Awaitable {
    /** @var T */
    private mixed $value;

    /** @param T $value */
    public function __construct(mixed $value) { $this->value = $value; }

    /** @return T */
    public function await(): mixed { return $this->value; }

	/**
	 * @template Tk
	 * @template Tv
	 * @param iterable<Tk, Awaitable<Tv>> $awaitables
	 * @return Generator<Tk, Awaitable<Tv>, null, void>
	 */
	public static function iterate(iterable $awaitables): Generator {
        foreach ($awaitables as $key => $awaitable) { yield $key => $awaitable; }
    }
}
`
	const firstSource = `<?php
namespace Psl\Async;
use Psl\Async\Awaitable;
/**
 * @template T
 * @param iterable<Awaitable<T>> $awaitables
 * @return T
 */
function first(iterable $awaitables): mixed {
    foreach (Awaitable::iterate($awaitables) as $first) {
        foreach ($awaitables as $awaitable) {
            if ($awaitable === $first) { continue; }
        }
        return $first->await();
    }
    throw new \\InvalidArgumentException();
}
/** @return string */
function validFirst(): mixed { return first([new Awaitable('value')]); }
/** @return int */
function invalidControl(): mixed { return 'wrong'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"awaitable.php": awaitableSource, "first.php": firstSource}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected only the wrong-return control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeTracksGenericRecursiveSearch(t *testing.T) {
	const source = `<?php
/** @template T */
interface Awaitable { /** @return T */ public function await(): mixed; }
/** @template T */
interface NodeInterface { /** @return T */ public function getValue(): mixed; }
class TreeNode implements NodeInterface {
    /** @param list<TreeNode> $children */
    public function __construct(private mixed $value, private array $children = []) {}
    public function getValue(): mixed { return $this->value; }
    /** @return list<TreeNode> */
    public function getChildren(): array { return $this->children; }
}

/** @template T @param iterable<Awaitable<T>> $awaitables @return T */
function first(iterable $awaitables): mixed {
    foreach ($awaitables as $awaitable) { return $awaitable->await(); }
    throw new RuntimeException();
}
/** @template T @param NodeInterface<T> $tree @param Closure(T): bool $predicate @return T|null */
function find(NodeInterface $tree, Closure $predicate): mixed {
    $value = $tree->getValue();
    if ($predicate($value)) { return $value; }
    if (!$tree instanceof TreeNode) { return null; }
    foreach ($tree->getChildren() as $child) {
        /** @var T|null $result */
        $result = find($child, $predicate);
        if (null !== $result) { return $result; }
    }
    return null;
}
/** @return int */
function invalidControl(): string { return 'wrong'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"generic-recursive-search.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected generic await/recursive search to pass with invalid control failing, got %#v", returnIssues)
	}
}

func TestReturnTypePreservesTemplateInRecursiveLocalVarAssertion(t *testing.T) {
	const source = `<?php
namespace Psl\Tree;
use Closure;
/** @template-covariant T */
interface NodeInterface {
    /** @return T */
    public function getValue(): mixed;
}
/** @template T @implements NodeInterface<T> */
final class TreeNode implements NodeInterface {
    /** @var T */
    private mixed $value;
    /** @var list<NodeInterface<T>> */
    private array $children;

    public function __construct(mixed $value, array $children = []) {
        $this->value = $value;
        $this->children = $children;
    }

    /** @return T */
    public function getValue(): mixed { return $this->value; }
    /** @return list<NodeInterface<T>> */
    public function getChildren(): array { return $this->children; }
}
/**
 * @template T
 * @param NodeInterface<T> $tree
 * @param Closure(T): bool $predicate
 * @return T|null
 */
function find(NodeInterface $tree, Closure $predicate): mixed {
    $value = $tree->getValue();
    if ($predicate($value)) { return $value; }
    if (!$tree instanceof TreeNode) { return null; }
    foreach ($tree->getChildren() as $child) {
        /** @var T|null $result */
        $result = find($child, $predicate);
        if (null !== $result) { return $result; }
    }
    return null;
}
/** @return int */
function invalidControl(): mixed { return 'wrong'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"tree-find.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected recursive generic result to pass and the bad control to fail, got %#v", returnIssues)
	}
}

func TestReturnTypeSubstitutesPromotedGenericInterfaceProperty(t *testing.T) {
	const source = `<?php
namespace Psl\Type;
/** @template T */
interface TypeInterface {
    /** @return T */
    public function assert(mixed $value): mixed;
}
namespace Psl\Type\Internal;
use Psl\Type\TypeInterface;
/**
 * @template I
 * @template O
 */
final class ConvertedType {
    /**
     * @param TypeInterface<I> $from
     * @param TypeInterface<O> $into
     */
    public function __construct(private TypeInterface $from, private TypeInterface $into) {}
    /** @return O */
    public function assert(mixed $value): mixed {
        return $this->into->assert($value);
    }
}
/** @return int */
function invalidControl(): string { return 'wrong'; }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"converted-generic-return.php": source}, 8)
	var returnIssues []AnalysisIssue
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			returnIssues = append(returnIssues, issue)
		}
	}
	if len(returnIssues) != 1 || !strings.Contains(returnIssues[0].Message, "invalidControl") {
		t.Fatalf("expected delegated generic assertion to pass with incompatible control failing, got %#v", returnIssues)
	}
}
