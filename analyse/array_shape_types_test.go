package analyse

import "testing"

func TestArrayShapeIsCompatibleWithNativeArray(t *testing.T) {
	shape := ParseType(`array{"route": string, ...<string, int>}`)
	if !ParseType("array").Accepts(shape) {
		t.Fatalf("native array should accept array-shape types, got shape %#v", shape)
	}
	if !ParseType("iterable").Accepts(shape) {
		t.Fatalf("iterable should accept array-shape types, got shape %#v", shape)
	}
	if ParseType("non-empty-array").Accepts(ParseType("array{}")) {
		t.Fatal("an empty array shape must not satisfy non-empty-array")
	}
}
