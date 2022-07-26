package registrable

import (
	"reflect"
	"sort"
)

// Registration holds info to be registered and used within your code.
// Good place to store any information required by the Registrable.
type Registration interface{}

// RegisterFn signature of a Registrable method.
type RegisterFn func() Registration

// used to detect return type of methods
var registerFnType RegisterFn = func() Registration {
	return nil
}

// Registrable make a method Registrable.
// Good place to store any dependencies for use within your Registrable methods.
type Registrable[H Registration] interface {
	Register(H)
}

// RegisterMethods Register methods on Registrable T to be ran.
// Uses reflection to find methods that return Registration H.
// That Registration is then passed to and called by Registrable T's Register method.
func RegisterMethods[H Registration, T Registrable[H]](registrable T) {
	st := reflect.TypeOf(registrable)
	sv := reflect.ValueOf(registrable)
	for i := 0; i < st.NumMethod(); i++ {
		mtype := st.Method(i)
		mvalue := sv.Method(i)

		register := mvalue.Type().AssignableTo(reflect.TypeOf(registerFnType))
		if register {
			sh := mtype.Func.Call([]reflect.Value{sv})
			registrable.Register(sh[0].Interface().(H))
		}
	}
}

// OrderedRegisterFn signature of an OrderedRegistrable method.
type OrderedRegisterFn func() (int, Registration)

// used to detect return type of methods
var orderedRegisterFnType OrderedRegisterFn = func() (int, Registration) {
	return 0, nil
}

// RegisterOrderedMethods Same as RegisterMethods uses an OrderedRegistration to order register calls.
func RegisterOrderedMethods[H Registration, T Registrable[H]](registrable T) {
	st := reflect.TypeOf(registrable)
	sv := reflect.ValueOf(registrable)

	keys := make([]int, 0)
	buckets := make(map[int][]Registration, 0)
	for i := 0; i < st.NumMethod(); i++ {
		mtype := st.Method(i)
		mvalue := sv.Method(i)

		register := mvalue.Type().AssignableTo(reflect.TypeOf(orderedRegisterFnType))
		if register {
			sh := mtype.Func.Call([]reflect.Value{sv})

			ordering := sh[0].Interface().(int)
			registration := sh[1].Interface().(H)

			if _, ok := buckets[ordering]; !ok {
				keys = append(keys, ordering)
				buckets[ordering] = make([]Registration, 0)
			}

			buckets[ordering] = append(buckets[ordering], registration)
		}
	}

	sort.Ints(keys)
	for _, key := range keys {
		for _, registration := range buckets[key] {
			registrable.Register(registration.(H))
		}
	}
}
