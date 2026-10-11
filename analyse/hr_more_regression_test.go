package analyse

import "testing"

func TestIntegerRangeMayBeCoveredByUnionOfAdjacentContracts(t *testing.T) {
	if !ParseType("0|positive-int").Accepts(ParseType("int<0, max>")) {
		t.Fatal("the union 0|positive-int should accept the complete non-negative integer range")
	}
	if !ParseType("0|null|positive-int").Accepts(ParseType("int<0, max>")) {
		t.Fatal("an unrelated nullable alternative should not prevent non-negative integer coverage")
	}
	if ParseType("0|positive-int").Accepts(ParseType("int")) {
		t.Fatal("0|positive-int must not accept negative integers")
	}
}

func TestNeverDoesNotMakeACompatibleReturnUnionInvalid(t *testing.T) {
	if !ParseType("string").Accepts(ParseType("'ready'|never")) {
		t.Fatal("a never branch should not invalidate the remaining compatible return type")
	}
}

func TestCapturedVariableInsideCastInNestedCallbackIsUsed(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": `<?php
final class CallbackProbe {
    public function with(mixed ...$arguments): void {}
}
function checksWeek(CallbackProbe $probe, DateTime $weekStart): void {
    $probe->with(static function (DateTime $weekEnd) use ($weekStart): int {
        return (int) $weekEnd->diff($weekStart)->days;
    });
}
`}, 1)
	if hasIssueContaining(issues, level1CoreCode, "unused use $weekStart") {
		t.Fatalf("a captured value inside a type-cast expression must count as used, got %#v", issues)
	}
}

func TestStringAccumulatorTernaryReturnFitsString(t *testing.T) {
	issues := analysePHP(t, `<?php
function initials(?string $firstName, ?string $lastName): string {
    $initials = '';
    if ($firstName !== '') {
        $initials .= strtoupper(substr($firstName, 0, 1));
    }
    if ($lastName !== '') {
        $initials .= strtoupper(substr($lastName, 0, 1));
    }
    return $initials !== '' ? $initials : 'NA';
}
`)
	if hasReturnTypeIssue(issues) {
		t.Fatalf("a string accumulator with a literal fallback must satisfy string, got %#v", issues)
	}
}
