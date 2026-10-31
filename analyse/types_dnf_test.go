package analyse

import "testing"

func TestParseTypePreservesDNFStructureWithoutChangingFlatCompatibility(t *testing.T) {
	typ := ParseType("(A&B)|(C&D)")
	if got := typ.dnfString(); got != "(A&B)|(C&D)" {
		t.Fatalf("lossless DNF string = %q, want %q", got, "(A&B)|(C&D)")
	}
	if got := typ.String(); got != "A|B|C|D" {
		t.Fatalf("flat compatibility string = %q, want %q", got, "A|B|C|D")
	}
	if _, ok := typ.SingleClassName(); ok {
		t.Fatal("DNF type must not become a single class")
	}
	if reparsed := ParseType(typ.dnfString()).dnfString(); reparsed != typ.dnfString() {
		t.Fatalf("DNF round trip = %q, want %q", reparsed, typ.dnfString())
	}
}

func TestParseTypeOnlyStripsBalancedOuterParentheses(t *testing.T) {
	for raw, expected := range map[string]string{
		"((A&B)|(C&D))": "(A&B)|(C&D)",
		"(A&B)|(C&D)":   "(A&B)|(C&D)",
		"A|(B&C)":       "(B&C)|A",
	} {
		if got := ParseType(raw).dnfString(); got != expected {
			t.Errorf("ParseType(%q).dnfString() = %q, want %q", raw, got, expected)
		}
	}
}

func TestParseTypeDNFDeduplicatesAndNullableRemovalDoesNotMutateCache(t *testing.T) {
	original := ParseType("(A&B)|(A&B)|null")
	if got := original.dnfString(); got != "(A&B)|null" {
		t.Fatalf("deduplicated DNF = %q", got)
	}
	refined := original.withoutBuiltin("null")
	if got := refined.dnfString(); got != "A&B" {
		t.Fatalf("refined DNF = %q, want %q", got, "A&B")
	}
	if got := original.dnfString(); got != "(A&B)|null" {
		t.Fatalf("cached original mutated after refinement: %q", got)
	}
	if got := ParseType("(A&B)|(A&B)|null").dnfString(); got != "(A&B)|null" {
		t.Fatalf("cached DNF changed after refinement: %q", got)
	}
}

func TestParseTypeDNFKeepsExistingFlatAcceptanceBehavior(t *testing.T) {
	declared := ParseType("A|B")
	actual := ParseType("A&B")
	if !declared.Accepts(actual) {
		t.Fatal("DNF metadata must not change existing flat acceptance semantics")
	}
}

func TestNumericStringSatisfiesNonEmptyString(t *testing.T) {
	if !ParseType("non-empty-string").Accepts(ParseType("numeric-string")) {
		t.Fatal("numeric-string is necessarily non-empty and must satisfy non-empty-string")
	}
	if ParseType("non-empty-string").Accepts(ParseType("string")) {
		t.Fatal("general string must not satisfy non-empty-string")
	}
}

func TestWithoutExactClassAlternative(t *testing.T) {
	current := ParseType("Constraint|string|null")
	refined := current.withoutExactClassAlternative("Constraint")
	if refined.Accepts(ParseType("Constraint")) || !ParseType("string|null").Accepts(refined) {
		t.Fatalf("withoutExactClassAlternative did not remove only the exact class branch: current=%s refined=%s", current, refined)
	}
}
