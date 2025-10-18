package mystructs

import (
	"testing"
)

func TestKVGroupInput(t *testing.T) {
	// Test creating KVGroupInput from string
	inputString := "name:John age:30 city:NewYork"

	kv, err := NewKVGroupInputFromString(inputString)
	if err != nil {
		t.Fatalf("Failed to create KVGroupInput from string: %v", err)
	}

	// Test that we have the expected pairs
	expectedPairs := []KVPair{
		{Key: "name", Value: "John"},
		{Key: "age", Value: "30"},
		{Key: "city", Value: "NewYork"},
	}

	if len(kv.KVPairs) != len(expectedPairs) {
		t.Fatalf("Expected %d pairs, got %d", len(expectedPairs), len(kv.KVPairs))
	}

	for i, expected := range expectedPairs {
		if kv.KVPairs[i].Key != expected.Key || kv.KVPairs[i].Value != expected.Value {
			t.Errorf("Pair %d: expected %+v, got %+v", i, expected, kv.KVPairs[i])
		}
	}

	// Test String() method
	result := kv.String()
	if result != inputString {
		t.Errorf("String() expected %s, got %s", inputString, result)
	}
}

func TestKVGroupInputEmpty(t *testing.T) {
	// Test empty input
	kv, err := NewKVGroupInputFromString("")
	if err != nil {
		t.Fatalf("Failed to create KVGroupInput from empty string: %v", err)
	}

	if len(kv.KVPairs) != 0 {
		t.Errorf("Expected empty KVPairs, got %d pairs", len(kv.KVPairs))
	}

	// Test String() with empty input
	if kv.String() != "" {
		t.Errorf("String() with empty input expected empty string, got %s", kv.String())
	}
}

func TestKVGroupInputGraphQL(t *testing.T) {
	// Test MarshalGQL
	kv := KVGroupInput{
		KVPairGroup: KVPairGroup{
			KVPairs: []KVPair{
				{Key: "name", Value: "Alice"},
				{Key: "age", Value: "25"},
			},
		},
	}

	marshaled, err := kv.MarshalGQL()
	if err != nil {
		t.Fatalf("MarshalGQL() failed: %v", err)
	}

	expected := `"name:Alice age:25"`
	if string(marshaled) != expected {
		t.Errorf("MarshalGQL() expected %s, got %s", expected, string(marshaled))
	}

	// Test UnmarshalGQL
	newKv := &KVGroupInput{}
	err = newKv.UnmarshalGQL("name:Bob age:35")
	if err != nil {
		t.Fatalf("UnmarshalGQL() failed: %v", err)
	}

	if len(newKv.KVPairs) != 2 {
		t.Fatalf("Expected 2 pairs after UnmarshalGQL, got %d", len(newKv.KVPairs))
	}

	if newKv.KVPairs[0].Key != "name" || newKv.KVPairs[0].Value != "Bob" {
		t.Errorf("First pair: expected {name:Bob}, got %+v", newKv.KVPairs[0])
	}

	if newKv.KVPairs[1].Key != "age" || newKv.KVPairs[1].Value != "35" {
		t.Errorf("Second pair: expected {age:35}, got %+v", newKv.KVPairs[1])
	}
}

func TestKVGroupInputComplexValues(t *testing.T) {
	// Test with values containing special characters (no spaces for now)
	inputString := "title:HelloWorld description:ThisIsATestWithNoSpaces url:https://example.com/path?param=value"

	kv, err := NewKVGroupInputFromString(inputString)
	if err != nil {
		t.Fatalf("Failed to create KVGroupInput with complex values: %v", err)
	}

	expectedPairs := []KVPair{
		{Key: "title", Value: "HelloWorld"},
		{Key: "description", Value: "ThisIsATestWithNoSpaces"},
		{Key: "url", Value: "https://example.com/path?param=value"},
	}

	if len(kv.KVPairs) != len(expectedPairs) {
		t.Fatalf("Expected %d pairs, got %d", len(expectedPairs), len(kv.KVPairs))
	}

	for i, expected := range expectedPairs {
		if kv.KVPairs[i].Key != expected.Key || kv.KVPairs[i].Value != expected.Value {
			t.Errorf("Pair %d: expected %+v, got %+v", i, expected, kv.KVPairs[i])
		}
	}
}

func TestKVGroupInputInvalidFormat(t *testing.T) {
	// Test invalid format (missing colon)
	_, err := NewKVGroupInputFromString("name:John age30 city:New York")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}

	// Test invalid format (no colon at all)
	_, err = NewKVGroupInputFromString("name:John age30 city")
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}

func TestKVGroupInputToKVPairGroup(t *testing.T) {
	kv := KVGroupInput{
		KVPairGroup: KVPairGroup{
			KVPairs: []KVPair{
				{Key: "name", Value: "Test"},
				{Key: "value", Value: "123"},
			},
		},
	}

	group := kv.ToKVPairGroup()

	if len(group.KVPairs) != 2 {
		t.Fatalf("Expected 2 pairs in KVPairGroup, got %d", len(group.KVPairs))
	}

	if group.KVPairs[0].Key != "name" || group.KVPairs[0].Value != "Test" {
		t.Errorf("First pair in KVPairGroup: expected {name:Test}, got %+v", group.KVPairs[0])
	}
}

func TestNewKVGroupInput(t *testing.T) {
	// Test creating KVGroupInput from KVPairGroup
	group := KVPairGroup{
		KVPairs: []KVPair{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	kv := NewKVGroupInput(group)

	if len(kv.KVPairs) != 2 {
		t.Fatalf("Expected 2 pairs, got %d", len(kv.KVPairs))
	}

	if kv.KVPairs[0].Key != "key1" || kv.KVPairs[0].Value != "value1" {
		t.Errorf("First pair: expected {key1:value1}, got %+v", kv.KVPairs[0])
	}
}
