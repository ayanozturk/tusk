package analyse

import "testing"

func TestLevel1CoreAcceptsCommonBuiltinConstants(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": `<?php
function builtins(): array {
    return [JSON_ERROR_NONE, DATE_ATOM, FILE_APPEND, FILTER_VALIDATE_BOOLEAN,
        FILTER_VALIDATE_EMAIL, PATHINFO_EXTENSION, UPLOAD_ERR_OK];
}
`}, 1)
	if hasIssueContaining(issues, level1CoreCode, "Constant") {
		t.Fatalf("PHP core and extension constants must resolve, got %#v", issues)
	}
}

func TestLevel1CoreTracksCapturedVariablesInNestedExpressions(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": `<?php
function checksWeek(DateTime $weekStart, DateTime $weekEnd): bool {
    $check = static function (DateTime $end) use ($weekStart): bool {
        return $end->diff($weekStart)->days === 6;
    };
    return $check($weekEnd);
}
`}, 1)
	if hasIssueContaining(issues, level1CoreCode, "unused use $weekStart") {
		t.Fatalf("a captured variable used as a chained method argument must count as used, got %#v", issues)
	}
}

func TestLevel1CoreDoesNotAssumeMutableEmptyArrayStaysEmpty(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": `<?php
function hasItems(array $values): bool {
    $items = [];
    foreach ($values as $value) {
        $items[] = $value;
    }
    return empty($items);
}
`}, 1)
	if hasIssueContaining(issues, level1CoreCode, "empty() always exists and is always falsy") {
		t.Fatalf("an array that is appended to cannot be assumed to stay empty, got %#v", issues)
	}
}
