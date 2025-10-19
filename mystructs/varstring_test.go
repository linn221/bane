package mystructs

import (
	"bytes"
	"testing"
)

func TestVarString(t *testing.T) {
	// Test the example from the user's description
	originalString := "{id:1}-{name:Henry Cohle}"

	// Create VarString
	vs, err := NewVarString(originalString)
	if err != nil {
		t.Fatalf("Failed to create VarString: %v", err)
	}

	// Test initial execution (should use default values)
	expected := "1-Henry Cohle"
	result := vs.Exec()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test Stringer interface
	if vs.String() != expected {
		t.Errorf("String() method failed. Expected %s, got %s", expected, vs.String())
	}

	// Test variable injection
	vs.Inject(map[string]string{"id": "3"})
	expected = "3-Henry Cohle"
	result = vs.Exec()
	if result != expected {
		t.Errorf("After injection, expected %s, got %s", expected, result)
	}
}

func TestVarStringComplex(t *testing.T) {
	// Test with more complex string
	originalString := "User {id:1} ({name:John Doe}) has {count:0} items in {category:general}"

	vs, err := NewVarString(originalString)
	if err != nil {
		t.Fatalf("Failed to create VarString: %v", err)
	}

	// Test default values
	expected := "User 1 (John Doe) has 0 items in general"
	result := vs.Exec()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test partial injection
	vs.Inject(map[string]string{
		"id":    "42",
		"count": "5",
	})

	expected = "User 42 (John Doe) has 5 items in general"
	result = vs.Exec()
	if result != expected {
		t.Errorf("After partial injection, expected %s, got %s", expected, result)
	}

	// Test full injection
	vs.Inject(map[string]string{
		"name":     "Alice Smith",
		"category": "electronics",
	})

	expected = "User 42 (Alice Smith) has 5 items in electronics"
	result = vs.Exec()
	if result != expected {
		t.Errorf("After full injection, expected %s, got %s", expected, result)
	}
}

func TestVarStringGORM(t *testing.T) {
	originalString := "{id:1}-{name:Henry Cohle}"

	vs, err := NewVarString(originalString)
	if err != nil {
		t.Fatalf("Failed to create VarString: %v", err)
	}

	// Test Value() method (serialization) - should return original string
	value, err := vs.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}

	// Verify Value() returns the original string
	if value != originalString {
		t.Errorf("Value() should return original string. Expected %s, got %s", originalString, value)
	}

	// Test Scan() method (deserialization)
	newVs := &VarString{}
	err = newVs.Scan(originalString)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	// Verify the deserialized VarString works the same
	if newVs.Exec() != vs.Exec() {
		t.Error("Deserialized VarString doesn't match original")
	}

	// Test that injection still works after deserialization
	newVs.Inject(map[string]string{"id": "42"})
	expected := "42-Henry Cohle"
	result := newVs.Exec()
	if result != expected {
		t.Errorf("After deserialization and injection, expected %s, got %s", expected, result)
	}
}

func TestVarStringGraphQL(t *testing.T) {
	originalString := "{id:1}-{name:Henry Cohle}"

	vs, err := NewVarString(originalString)
	if err != nil {
		t.Fatalf("Failed to create VarString: %v", err)
	}

	// Test MarshalGQL
	var buf bytes.Buffer
	vs.MarshalGQL(&buf)
	marshaled := buf.Bytes()

	// Should return the executed string (with default values)
	expected := `"1-Henry Cohle"`
	if string(marshaled) != expected {
		t.Errorf("MarshalGQL() expected %s, got %s", expected, string(marshaled))
	}

	// Test with injection
	vs.Inject(map[string]string{"id": "99"})
	var buf2 bytes.Buffer
	vs.MarshalGQL(&buf2)
	marshaled = buf2.Bytes()

	expected = `"99-Henry Cohle"`
	if string(marshaled) != expected {
		t.Errorf("MarshalGQL() after injection expected %s, got %s", expected, string(marshaled))
	}

	// Test UnmarshalGQL
	newVs := &VarString{}
	err = newVs.UnmarshalGQL(originalString)
	if err != nil {
		t.Fatalf("UnmarshalGQL() failed: %v", err)
	}

	// Should work the same as the original (both should have default values)
	if newVs.Exec() != "1-Henry Cohle" {
		t.Errorf("UnmarshalGQL() expected default values, got %s", newVs.Exec())
	}

	// Test injection after unmarshaling
	newVs.Inject(map[string]string{"id": "77"})
	expectedResult := "77-Henry Cohle"
	if newVs.Exec() != expectedResult {
		t.Errorf("After UnmarshalGQL and injection, expected %s, got %s", expectedResult, newVs.Exec())
	}
}
