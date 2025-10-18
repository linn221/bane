package mystructs

import (
	"testing"
)

func TestVarKVGroupValue(t *testing.T) {
	tests := []struct {
		name     string
		vkg      VarKVGroup
		expected string
	}{
		{
			name:     "empty VarKVGroup",
			vkg:      VarKVGroup{VarKVs: []VarKV{}},
			expected: "",
		},
		{
			name: "single VarKV",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
				},
			},
			expected: "name|John",
		},
		{
			name: "multiple VarKVs",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
					{
						Key:   VarString{OriginalString: "age"},
						Value: VarString{OriginalString: "30"},
					},
					{
						Key:   VarString{OriginalString: "city"},
						Value: VarString{OriginalString: "New York"},
					},
				},
			},
			expected: "name|John age|30 city|New York",
		},
		{
			name: "VarKVs with complex strings",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "url"},
						Value: VarString{OriginalString: "https://example.com"},
					},
					{
						Key:   VarString{OriginalString: "description"},
						Value: VarString{OriginalString: "A_test_with_spaces_and_symbols!@#$%"},
					},
				},
			},
			expected: "url|https://example.com description|A_test_with_spaces_and_symbols!@#$%",
		},
		{
			name: "VarKVs with pipe characters in values",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "pattern"},
						Value: VarString{OriginalString: "key|value"},
					},
					{
						Key:   VarString{OriginalString: "separator"},
						Value: VarString{OriginalString: "||"},
					},
				},
			},
			expected: "pattern|key|value separator|||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.vkg.Value()
			if err != nil {
				t.Fatalf("Value() failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestVarKVGroupScan(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    VarKVGroup
		expectError bool
		errorMsg    string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: VarKVGroup{VarKVs: []VarKV{}},
		},
		{
			name:     "empty string",
			input:    "",
			expected: VarKVGroup{VarKVs: []VarKV{}},
		},
		{
			name:  "single key-value pair",
			input: "name|John",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
				},
			},
		},
		{
			name:  "multiple key-value pairs",
			input: "name|John age|30 city|NewYork",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
					{
						Key:   VarString{OriginalString: "age"},
						Value: VarString{OriginalString: "30"},
					},
					{
						Key:   VarString{OriginalString: "city"},
						Value: VarString{OriginalString: "NewYork"},
					},
				},
			},
		},
		{
			name:  "complex strings with spaces and symbols",
			input: "url|https://example.com description|A_test_with_spaces_and_symbols!@#$%",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "url"},
						Value: VarString{OriginalString: "https://example.com"},
					},
					{
						Key:   VarString{OriginalString: "description"},
						Value: VarString{OriginalString: "A_test_with_spaces_and_symbols!@#$%"},
					},
				},
			},
		},
		{
			name:  "pipe characters in values",
			input: "pattern|key|value separator|||",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "pattern"},
						Value: VarString{OriginalString: "key|value"},
					},
					{
						Key:   VarString{OriginalString: "separator"},
						Value: VarString{OriginalString: "||"},
					},
				},
			},
		},
		{
			name:  "byte slice input",
			input: []byte("name|John age|30"),
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
					{
						Key:   VarString{OriginalString: "age"},
						Value: VarString{OriginalString: "30"},
					},
				},
			},
		},
		{
			name:        "invalid input type",
			input:       123,
			expectError: true,
			errorMsg:    "cannot scan int into VarKVGroup",
		},
		{
			name:        "missing pipe in key-value pair",
			input:       "nameJohn",
			expectError: true,
			errorMsg:    "invalid format: missing pipe in 'nameJohn'",
		},
		{
			name:        "mixed valid and invalid pairs",
			input:       "name|John age30 city|New York",
			expectError: true,
			errorMsg:    "invalid format: missing pipe in 'age30'",
		},
		{
			name:  "empty key",
			input: "|value",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: ""},
						Value: VarString{OriginalString: "value"},
					},
				},
			},
		},
		{
			name:  "empty value",
			input: "key|",
			expected: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "key"},
						Value: VarString{OriginalString: ""},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vkg VarKVGroup
			err := vkg.Scan(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Scan() failed: %v", err)
			}

			if len(vkg.VarKVs) != len(tt.expected.VarKVs) {
				t.Errorf("Expected %d VarKVs, got %d", len(tt.expected.VarKVs), len(vkg.VarKVs))
				return
			}

			for i, expectedKV := range tt.expected.VarKVs {
				if i >= len(vkg.VarKVs) {
					t.Errorf("Missing VarKV at index %d", i)
					continue
				}

				actualKV := vkg.VarKVs[i]
				if actualKV.Key.OriginalString != expectedKV.Key.OriginalString {
					t.Errorf("VarKV[%d].Key.OriginalString: expected '%s', got '%s'", i, expectedKV.Key.OriginalString, actualKV.Key.OriginalString)
				}
				if actualKV.Value.OriginalString != expectedKV.Value.OriginalString {
					t.Errorf("VarKV[%d].Value.OriginalString: expected '%s', got '%s'", i, expectedKV.Value.OriginalString, actualKV.Value.OriginalString)
				}
			}
		})
	}
}

func TestVarKVGroupRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		vkg  VarKVGroup
	}{
		{
			name: "empty VarKVGroup",
			vkg:  VarKVGroup{VarKVs: []VarKV{}},
		},
		{
			name: "single VarKV",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
				},
			},
		},
		{
			name: "multiple VarKVs",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "name"},
						Value: VarString{OriginalString: "John"},
					},
					{
						Key:   VarString{OriginalString: "age"},
						Value: VarString{OriginalString: "30"},
					},
					{
						Key:   VarString{OriginalString: "city"},
						Value: VarString{OriginalString: "NewYork"},
					},
				},
			},
		},
		{
			name: "complex strings",
			vkg: VarKVGroup{
				VarKVs: []VarKV{
					{
						Key:   VarString{OriginalString: "url"},
						Value: VarString{OriginalString: "https://example.com/path?param=value"},
					},
					{
						Key:   VarString{OriginalString: "description"},
						Value: VarString{OriginalString: "A_test_with_spaces_and_symbols!@#$%"},
					},
					{
						Key:   VarString{OriginalString: "pattern"},
						Value: VarString{OriginalString: "key|value"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Value() -> Scan() round trip
			value, err := tt.vkg.Value()
			if err != nil {
				t.Fatalf("Value() failed: %v", err)
			}

			var newVkg VarKVGroup
			err = newVkg.Scan(value)
			if err != nil {
				t.Fatalf("Scan() failed: %v", err)
			}

			// Compare the results
			if len(newVkg.VarKVs) != len(tt.vkg.VarKVs) {
				t.Errorf("Round trip failed: expected %d VarKVs, got %d", len(tt.vkg.VarKVs), len(newVkg.VarKVs))
				return
			}

			for i, expectedKV := range tt.vkg.VarKVs {
				actualKV := newVkg.VarKVs[i]
				if actualKV.Key.OriginalString != expectedKV.Key.OriginalString {
					t.Errorf("Round trip failed: VarKV[%d].Key.OriginalString: expected '%s', got '%s'", i, expectedKV.Key.OriginalString, actualKV.Key.OriginalString)
				}
				if actualKV.Value.OriginalString != expectedKV.Value.OriginalString {
					t.Errorf("Round trip failed: VarKV[%d].Value.OriginalString: expected '%s', got '%s'", i, expectedKV.Value.OriginalString, actualKV.Value.OriginalString)
				}
			}
		})
	}
}

func TestVarKVGroupWithVarStringPlaceholders(t *testing.T) {
	// Test with VarStrings that have placeholders
	keyVarString, err := NewVarString("{id:1}")
	if err != nil {
		t.Fatalf("Failed to create key VarString: %v", err)
	}

	valueVarString, err := NewVarString("{name:John}")
	if err != nil {
		t.Fatalf("Failed to create value VarString: %v", err)
	}

	vkg := VarKVGroup{
		VarKVs: []VarKV{
			{
				Key:   *keyVarString,
				Value: *valueVarString,
			},
		},
	}

	// Test Value() method
	value, err := vkg.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}

	expected := "{id:1}|{name:John}"
	if value != expected {
		t.Errorf("Expected %s, got %s", expected, value)
	}

	// Test Scan() method
	var newVkg VarKVGroup
	err = newVkg.Scan(value)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	// Verify the scanned VarKVs have the correct OriginalString values
	if len(newVkg.VarKVs) != 1 {
		t.Fatalf("Expected 1 VarKV, got %d", len(newVkg.VarKVs))
	}

	scannedKV := newVkg.VarKVs[0]
	if scannedKV.Key.OriginalString != "{id:1}" {
		t.Errorf("Expected key OriginalString '{id:1}', got '%s'", scannedKV.Key.OriginalString)
	}
	if scannedKV.Value.OriginalString != "{name:John}" {
		t.Errorf("Expected value OriginalString '{name:John}', got '%s'", scannedKV.Value.OriginalString)
	}

	// Test that the scanned VarStrings work correctly
	keyResult := scannedKV.Key.Exec()
	expectedKey := "1"
	if keyResult != expectedKey {
		t.Errorf("Expected key Exec() result '%s', got '%s'", expectedKey, keyResult)
	}

	valueResult := scannedKV.Value.Exec()
	expectedValue := "John"
	if valueResult != expectedValue {
		t.Errorf("Expected value Exec() result '%s', got '%s'", expectedValue, valueResult)
	}
}

func TestVarKVGroupEdgeCases(t *testing.T) {
	// Test with empty key but valid value
	vkg := VarKVGroup{
		VarKVs: []VarKV{
			{
				Key:   VarString{OriginalString: ""},
				Value: VarString{OriginalString: "value"},
			},
		},
	}

	value, err := vkg.Value()
	if err != nil {
		t.Fatalf("Value() with empty key failed: %v", err)
	}

	expected := "|value"
	if value != expected {
		t.Errorf("Expected '%s', got '%s'", expected, value)
	}

	// Test scanning the empty key case
	var newVkg VarKVGroup
	err = newVkg.Scan(value)
	if err != nil {
		t.Fatalf("Scan() with empty key failed: %v", err)
	}

	if len(newVkg.VarKVs) != 1 {
		t.Fatalf("Expected 1 VarKV, got %d", len(newVkg.VarKVs))
	}

	scannedKV := newVkg.VarKVs[0]
	if scannedKV.Key.OriginalString != "" {
		t.Errorf("Expected empty key OriginalString, got '%s'", scannedKV.Key.OriginalString)
	}
	if scannedKV.Value.OriginalString != "value" {
		t.Errorf("Expected value OriginalString 'value', got '%s'", scannedKV.Value.OriginalString)
	}
}
