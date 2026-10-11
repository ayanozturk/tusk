package analyse

import (
	"testing"

	"github.com/ayanozturk/tusk/ast"
)

func TestNarrowingAcrossEarlyReturnCompoundGuards(t *testing.T) {
	issues := analysePHP(t, `<?php
function consume(string $value): void {}
function requireText(?string $value): string {
    if (!is_string($value) || $value === '') {
        throw new \InvalidArgumentException();
    }
    consume($value);
    return $value;
}
`)
	for _, issue := range issues {
		if issue.Code == "A.ARG.TYPE" || issue.Code == "A.RETURN.TYPE" {
			t.Fatalf("compound early-return guards must refine subsequent argument and return types, got %#v", issues)
		}
	}
}

func TestNonNegativeIntegerNormalizationSatisfiesRangeContract(t *testing.T) {
	issues := analysePHP(t, `<?php
/** @return int<0, max> */
function normalize(int $value): int {
    return $value >= 0 ? $value : 0;
}
`)
	if hasReturnTypeIssue(issues) {
		t.Fatalf("a non-negative integer branch join must satisfy int<0, max>, got %#v", issues)
	}
}

func TestModuloAcceptsIntegerLiteralAndIntegerVariable(t *testing.T) {
	issues := analysePHP(t, `<?php
function divisible(): bool {
    for ($i = 1; $i <= 30; $i++) {
        $divisible = $i % 7 === 0;
    }
    return $divisible;
}
`)
	if countIssuesWithCode(issues, "A.BINARY.OP.INVALID") != 0 {
		t.Fatalf("modulo accepts integer literals and integer variables, got %#v", issues)
	}
}

func TestGenericCallbackMethodReturnBindsClosureResult(t *testing.T) {
	issues := analysePHP(t, `<?php
final class Widget {}
final class Item {}
final class CacheItem {}
/** @template T */
interface CallbackInterface {}
final class CacheLike {
    /**
     * @template T
     * @param (callable(CacheItem,bool):T)|(callable(Item,bool):T)|CallbackInterface<T> $callback
     * @return T
     */
    public function get(string $key, callable $callback): mixed {
        return $callback(new Item());
    }
}
function widget(CacheLike $cache): Widget {
    return $cache->get('widget', static function (Item $item): Widget {
        return new Widget();
    });
}
`)
	if hasReturnTypeIssue(issues) {
		t.Fatalf("generic callback return types must bind to the closure result, got %#v", issues)
	}
}

func TestOpenEndedArrayShapePHPDocIsAnArrayContract(t *testing.T) {
	issues := analysePHP(t, `<?php
/**
 * @psalm-type Route = array{path: string}
 * @psalm-type Routes = array{"when@dev"?: array<string, Route>, ...<string, Route>}
 */
final class RouteConfig {
    /** @param Routes $config @psalm-return Routes */
    public static function config(array $config): array { return $config; }
}
`)
	for _, issue := range issues {
		if issue.Code == level2PHPDocClassCode || issue.Code == level2PHPDocParamTypeCode || issue.Code == level2PHPDocReturnTypeCode || issue.Code == "A.RETURN.TYPE" {
			t.Fatalf("open-ended array shapes must remain valid array contracts, got %#v", issues)
		}
	}
}

func TestIntegerRangeBoundKeywordsStayGlobalInsideNamespaces(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{
		"length.php": `<?php
namespace Symfony\Component\Validator\Constraints;
final class Length {
    /** @param int<0, max>|null $min */
    public function __construct(?int $min = null, ?int $max = null) {}
}
`,
		"form.php": `<?php
namespace Demo;
use Symfony\Component\Validator\Constraints\Length;
new Length(min: 2, max: 100);
`,
	}, 5)
	for _, issue := range issues {
		if issue.Filename == "form.php" && issue.Code == "A.ARG.TYPE" {
			t.Fatalf("max in an integer range is a bound keyword, not a namespaced class, got %#v", issues)
		}
	}
}

func TestEmptyArrayInferenceDoesNotProveMutableVariableFalsiness(t *testing.T) {
	ctx := &AnalysisContext{}
	scope := &functionScope{}
	scope.setVariable("items", ParseType("empty-array"))
	var issues []AnalysisIssue
	appendLevel1NullableCheck("test.php", &ast.VariableNode{Name: "items"}, "empty", scope, ctx, &issues)
	if len(issues) != 0 {
		t.Fatalf("empty-array inference alone cannot prove a mutable variable remains empty, got %#v", issues)
	}
}

func TestCapturedVariableInsideNestedCallbackCallIsUsed(t *testing.T) {
	issues := runAnalysisLevelOnFiles(t, map[string]string{"test.php": `<?php
final class CallbackProbe {
    public function callback(callable $callback): callable { return $callback; }
    public function with(mixed ...$arguments): void {}
}
function checksWeek(CallbackProbe $probe, DateTime $weekStart): void {
    $probe->with($probe->callback(static function (DateTime $weekEnd) use ($weekStart): bool {
        return $weekEnd->diff($weekStart)->days === 6;
    }));
}
`}, 1)
	if hasIssueContaining(issues, level1CoreCode, "unused use $weekStart") {
		t.Fatalf("a captured value referenced inside a nested callback call must count as used, got %#v", issues)
	}
}
