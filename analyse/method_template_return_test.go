package analyse

import (
	"testing"

	"github.com/ayanozturk/tusk/ast"
)

func TestMethodTemplateBindingFromStructuredGenericParameters(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
		actual   GenericInstance
		want     string
	}{
		{"list", "list<T>", GenericInstance{ClassName: "list", TypeArguments: []string{"string"}}, "string"},
		{"array", "array<T>", GenericInstance{ClassName: "array", TypeArguments: []string{"int", "string"}}, "string"},
		{"nested", "Box<list<T>>", GenericInstance{ClassName: "Box", TypeArguments: []string{"list<string>"}}, "string"},
		{"incompatible", "list<T>", GenericInstance{ClassName: "Box", TypeArguments: []string{"string"}}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			method := ResolvedMethod{
				TemplateParams: []string{"T"},
				Params:         []ResolvedParam{{Name: "values", Type: test.expected}},
				ReturnType:     "T",
			}
			scope := &functionScope{}
			scope.setGenericContext("values", test.actual)
			args := []ast.Node{&ast.NamedArgumentNode{Name: "values", Value: &ast.VariableNode{Name: "values"}}}
			if got := bindCallSiteMethodTemplates(method, args, scope, nil, "")["T"]; got != test.want {
				t.Fatalf("binding = %q, want %q", got, test.want)
			}
		})
	}
}

func TestMethodTemplateVariadicBindingsJoin(t *testing.T) {
	method := ResolvedMethod{
		TemplateParams: []string{"Value"},
		Params:         []ResolvedParam{{Name: "values", Type: "Value", IsVariadic: true}},
		ReturnType:     "Value",
	}
	scope := &functionScope{}
	scope.setVariable("number", ParseType("int"))
	scope.setVariable("text", ParseType("string"))
	args := []ast.Node{&ast.VariableNode{Name: "number"}, &ast.VariableNode{Name: "text"}}
	got := ParseType(bindCallSiteMethodTemplates(method, args, scope, nil, "")["Value"])
	if !got.Accepts(ParseType("int")) || !got.Accepts(ParseType("string")) || got.Accepts(ParseType("bool")) {
		t.Fatalf("variadic method bindings must join, got %q", got.String())
	}
}

func TestMethodTemplateBindingRetainsViolatedBound(t *testing.T) {
	method := ResolvedMethod{
		TemplateParams: []string{"Value"},
		TemplateBounds: []string{"int"},
		Params:         []ResolvedParam{{Name: "values", Type: "array<Value>"}},
		ReturnType:     "Value",
	}
	scope := &functionScope{}
	scope.setGenericContext("values", GenericInstance{ClassName: "array", TypeArguments: []string{"string"}})
	args := []ast.Node{&ast.VariableNode{Name: "values"}}
	if got := bindCallSiteMethodTemplates(method, args, scope, nil, "")["Value"]; got != "int" {
		t.Fatalf("violated template bound = %q, want int", got)
	}
}

func TestMethodTemplateBindingPreservesClassStringForms(t *testing.T) {
	for _, test := range []struct {
		param string
		bound string
	}{
		{"class-string<T>", "object"},
		{"class-string<T>|string", "object"},
		{"T", "class-string"},
	} {
		method := ResolvedMethod{
			TemplateParams: []string{"T"},
			TemplateBounds: []string{test.bound},
			Params:         []ResolvedParam{{Name: "class", Type: test.param}},
			ReturnType:     "T",
		}
		args := []ast.Node{&ast.ClassConstFetchNode{Class: "Item", Const: "class"}}
		if got := bindCallSiteMethodTemplates(method, args, nil, nil, "")["T"]; got != "Item" {
			t.Errorf("%s bound by %s: got %q, want Item", test.param, test.bound, got)
		}
	}
}

func TestGenericMethodArgumentsAcceptForwardedNullableTemplates(t *testing.T) {
	const source = `<?php
class Selector {
    /**
     * @template Value of class-string|null
     * @param Value $name
     */
    public function select(?string $name = null): void {}
    /**
     * @template T of class-string|null
     * @param T $name
     */
    public function forward(?string $name = null): void {
        $this->select($name);
    }
}
`
	if issues := runAnalysisLevelOnFiles(t, map[string]string{"forward.php": source}, 5); hasArgTypeIssue(issues) {
		t.Fatalf("nullable caller templates must expand to their bounds, got %#v", issues)
	}
}

func TestArgumentTemplateBoundsPreserveUnionAndIntersectionStructure(t *testing.T) {
	for _, test := range []struct {
		actual string
		bound  string
		want   string
	}{
		{"T|null", "string", "string|null"},
		{"T|int", "", "mixed|int"},
		{"T&Marker", "Left|Right", "(Left|Right)&Marker"},
		{"Unrelated|null", "string", "Unrelated|null"},
	} {
		scope := &functionScope{functionScopeContext: &functionScopeContext{
			templateBounds: map[string]string{"t": test.bound},
		}}
		scope.setVariable("value", ParseType(test.actual))
		got := inferArgumentTypeWithFacts("", &ast.VariableNode{Name: "value"}, scope, nil)
		if got.dnfString() != ParseType(test.want).dnfString() {
			t.Errorf("%s bounded by %q = %s, want %s", test.actual, test.bound, got.dnfString(), test.want)
		}
	}
}

func TestMethodTemplateReturnBindsStructuredParameters(t *testing.T) {
	const source = `<?php
class Selector {
    /** @template Value @param list<Value> $values @return Value */
    public function first(array $values) { return $values[0]; }
    /** @template Value @param list<Value> $values @return Value */
    public static function select(array $values) { return $values[0]; }
}
class ChildSelector extends Selector {}
interface Selection {
    /** @template Value @param list<Value> $values @return Value */
    public function first(array $values);
}
/** @param list<string> $values */
function clean(ChildSelector $selector, Selection $interface, array $values): string {
    $selector->first(values: $values);
    $interface->first($values);
    return Selector::select($values);
}
/** @param list<string> $values */
function wrongInherited(ChildSelector $selector, array $values): int { return $selector->first(values: $values); }
/** @param list<string> $values */
function wrongInterface(Selection $selector, array $values): int { return $selector->first($values); }
/** @param list<string> $values */
function wrongStatic(array $values): int { return Selector::select($values); }
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"method-templates.php": source}, 5)
	count := 0
	for _, issue := range issues {
		if issue.Code == "A.RETURN.TYPE" {
			count++
			if issue.Line != 20 && issue.Line != 22 && issue.Line != 24 {
				t.Errorf("unexpected return mismatch: %#v", issue)
			}
		}
		if issue.Code == "A.ARG.TYPE" {
			t.Errorf("unexpected argument mismatch: %#v", issue)
		}
	}
	if count != 3 {
		t.Fatalf("want three structured method return mismatches, got %#v", issues)
	}
}

func TestStructuredMethodTemplateArgumentsRespectDeclaredBounds(t *testing.T) {
	const source = `<?php
class Selector {
    /** @template Value of int @param array<Value> $values @return Value */
    public function first(array $values) { return $values[0]; }
}
class Collector {
    /** @template Value of int @param array<Value> $values */
    public function __construct(array $values) {}
}
/** @param array<int> $values */
function clean(Selector $selector, array $values): int {
    new Collector($values);
    return $selector->first(values: $values);
}
/** @param array<string> $values */
function wrong(Selector $selector, array $values): void {
    $selector->first(values: $values);
    new Collector($values);
}
`
	issues := runAnalysisLevelOnFiles(t, map[string]string{"method-template-bounds.php": source}, 5)
	count := 0
	for _, issue := range issues {
		if issue.Code == "A.ARG.TYPE" {
			count++
			if issue.Line != 17 && issue.Line != 18 {
				t.Errorf("unexpected bound mismatch: %#v", issue)
			}
		}
	}
	if count != 2 {
		t.Fatalf("want bounded method and constructor mismatches, got %#v", issues)
	}
}
