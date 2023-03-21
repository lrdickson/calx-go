package kernel

import "testing"

func TestUpdate(t *testing.T) {
	// Run the output function
	input := make(map[string]string)
	input["a"] = "return 1"
	input["b"] = "return 2"
	goKernel := NewKernel()
	output := goKernel.Update(input)

	// Check the output
	expected := make(map[string]string)
	expected["a"] = "1"
	expected["b"] = "2"
	for key, expectedValue := range expected {
		if returnedValue, exists := output[key]; exists {
			if returnedValue != expectedValue {
				t.Fatalf("Update(%v) did not return ouput[\"%s\"] = \"%s\"", input, key, expectedValue)
			}
		} else {
			t.Fatalf("Update(%v) ouput did not contain key \"%s\"", input, key)
		}
	}
}
