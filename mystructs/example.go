package mystructs

import (
	"bytes"
	"fmt"
)

// Example demonstrates the usage of both VarString and KVGroupInput
func Example() {
	fmt.Println("=== VarString Example ===")

	// Create a VarString with the example from the user's description
	vs, err := NewVarString("{id:1}-{name:Henry Cohle}")
	if err != nil {
		fmt.Printf("Error creating VarString: %v\n", err)
		return
	}

	// Initially, it uses default values
	fmt.Printf("Default execution: %s\n", vs.Exec()) // Output: "1-Henry Cohle"

	// Inject new values
	vs.Inject(map[string]string{"id": "3"})
	fmt.Printf("After injection: %s\n", vs.Exec()) // Output: "3-Henry Cohle"

	// The Stringer interface calls Exec()
	fmt.Printf("Stringer output: %s\n", vs) // Output: "3-Henry Cohle"

	// Demonstrate GORM compatibility
	// This stores only the original string in the database
	value, _ := vs.Value()
	fmt.Printf("GORM Value (original string): %s\n", value)

	// Demonstrate GraphQL marshaling
	var buf bytes.Buffer
	vs.MarshalGQL(&buf)
	fmt.Printf("GraphQL Marshal: %s\n", buf.String())

	fmt.Println("\n=== KVGroupInput Example ===")

	// Create KVGroup from string
	kv, err := NewKVGroupFromString("name:John age:30 city:NewYork")
	if err != nil {
		fmt.Printf("Error creating KVGroupInput: %v\n", err)
		return
	}

	fmt.Printf("KVGroup string: %s\n", kv.String())

	// Demonstrate GraphQL marshaling
	var kvBuf bytes.Buffer
	kv.MarshalGQL(&kvBuf)
	fmt.Printf("GraphQL Marshal: %s\n", kvBuf.String())

	// Demonstrate conversion stays as KVGroup
	group := kv.ToKVGroup()
	fmt.Printf("Converted to KVGroup: %+v\n", group)
}
