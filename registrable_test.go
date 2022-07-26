package registrable

import (
	"fmt"
	"testing"
)

type TestRegistration struct {
	Remove func()
}

type TestRegistrable struct {
	ToBeRemoved map[string]bool
}

func (self TestRegistrable) HelloWorld() Registration {
	return TestRegistration{
		Remove: func() {
			delete(self.ToBeRemoved, "Hello World")
		},
	}
}

func (self TestRegistrable) GoodbyeWorld() Registration {
	return TestRegistration{
		Remove: func() {
			delete(self.ToBeRemoved, "Goodbye World")
		},
	}
}

func (self TestRegistrable) Register(tr TestRegistration) {
	tr.Remove()
}

// TestRegistrableStructure Ensure registrable methods run
func TestRegistrableStructure(t *testing.T) {
	toBeRemoved := make(map[string]bool)
	toBeRemoved["Hello World"] = false
	toBeRemoved["Goodbye World"] = false

	tr := &TestRegistrable{toBeRemoved}
	RegisterMethods[TestRegistration](tr)

	if len(toBeRemoved) != 0 {
		t.Fatalf("map should be empty but is not: %v\n", toBeRemoved)
	}
}

type OrderedTestRegistration struct {
	Push func()
}

type OrderedTestRegistrable struct {
	Numbers []int
}

func (self *OrderedTestRegistrable) PushA() (int, Registration) {
	return 3, OrderedTestRegistration{
		Push: func() {
			self.Numbers = append(self.Numbers, 3)
		},
	}
}

func (self *OrderedTestRegistrable) PushB() (int, Registration) {
	return 2, OrderedTestRegistration{
		Push: func() {
			self.Numbers = append(self.Numbers, 2)
		},
	}
}

func (self *OrderedTestRegistrable) PushC() (int, Registration) {
	return 1, OrderedTestRegistration{
		Push: func() {
			self.Numbers = append(self.Numbers, 1)
		},
	}
}

func (self OrderedTestRegistrable) Register(tr OrderedTestRegistration) {
	tr.Push()
}

// TestOrderedRegistrableStructure Ensure registrable methods run
func TestOrderedRegistrableStructure(t *testing.T) {
	tr := &OrderedTestRegistrable{make([]int, 0)}
	RegisterOrderedMethods[OrderedTestRegistration](tr)

	if fmt.Sprintf("%v", tr.Numbers) != "[1 2 3]" {
		t.Fatalf("numbers not equal to [1 2 3], got: %v\n", tr.Numbers)
	}
}
