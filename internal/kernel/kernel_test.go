package kernel

import (
	"testing"
	"time"
)

func TestAddWorker(t *testing.T) {
	// Add a worker to the kernel
	goKernel := NewKernel()
	done := make(chan string)
	goKernel.addWorker("test", "return 1", done)
	activeWorkerCount := goKernel.getActiveCount()
	if activeWorkerCount != 1 {
		t.Fatal("Go Kernel should have 1 active worker but instead has", activeWorkerCount)
	}

	// Stop the worker
	goKernel.stop("test")
	time.Sleep(time.Millisecond)
	activeWorkerCount = goKernel.getActiveCount()
	if activeWorkerCount != 0 {
		t.Fatal("Go Kernel should have 0 active workers but instead has", activeWorkerCount)
	}
}

func checkUpdate(t *testing.T, input, expected map[string]string) {
	// Start the kernel
	goKernel := NewKernel()
	output := goKernel.Update(input)

	// Check the output
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

func TestBasic(t *testing.T) {
	input := make(map[string]string)
	input["a"] = "return 1"
	input["b"] = "return 2"
	expected := make(map[string]string)
	expected["a"] = "1"
	expected["b"] = "2"
	checkUpdate(t, input, expected)
}

func TestImport(t *testing.T) {
	input := make(map[string]string)
	input["math"] = "return math.Pow(2,3)"
	expected := make(map[string]string)
	expected["math"] = "8"
	checkUpdate(t, input, expected)
}
