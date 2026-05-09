package compat

import "testing"

func TestSchemaCompatible(t *testing.T) {
	t.Parallel()
	ok, err := SchemaCompatible("0.0.2", "0.0.1")
	if err != nil {
		t.Fatalf("schema compatible error: %v", err)
	}
	if !ok {
		t.Fatalf("expected schema compatibility")
	}

	ok, err = SchemaCompatible("1.0.0", "0.0.1")
	if err != nil {
		t.Fatalf("schema incompatible error: %v", err)
	}
	if ok {
		t.Fatalf("expected schema incompatibility")
	}
}

func TestPluginVersionInRange(t *testing.T) {
	t.Parallel()
	ok, err := PluginVersionInRange("0.2.0", ">=0.1.0, <1.0.0")
	if err != nil {
		t.Fatalf("plugin range error: %v", err)
	}
	if !ok {
		t.Fatalf("expected plugin to be in range")
	}
}
