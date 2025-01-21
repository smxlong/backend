// backend
package backend

import (
	"io"
	"reflect"
)

// Injector is a dependency injecting function invoker.
type Injector struct {
	// singletons holds instances of the given types which will be directly
	// returned when requested.
	singletons map[reflect.Type]interface{}
	// transients holds instances of the given types which will be created
	// using the given factory function when requested.
	transients map[reflect.Type]interface{}
}

// NewInjector creates a new Injector.
func NewInjector() *Injector {
	return &Injector{
		singletons: make(map[reflect.Type]interface{}),
		transients: make(map[reflect.Type]interface{}),
	}
}

// typeIsRegistered returns true if the given type is registered in the injector.
func (i *Injector) typeIsRegistered(t reflect.Type) bool {
	_, singleton := i.singletons[t]
	_, transient := i.transients[t]
	return singleton || transient
}

// AddSingleton registers the given instance as a singleton for the given type.
func (i *Injector) AddSingleton(instance interface{}) error {
	t := reflect.TypeOf(instance)
	if i.typeIsRegistered(t) {
		return ErrTypeAlreadyRegistered
	}
	i.singletons[t] = instance
	return nil
}

// AddTransient registers the given factory function as a transient for the given type.
func (i *Injector) AddTransient(factory interface{}) error {
	f := reflect.TypeOf(factory)
	if f.Kind() != reflect.Func {
		return ErrInvalidFactory
	}
	t := f.Out(0)
	if i.typeIsRegistered(t) {
		return ErrTypeAlreadyRegistered
	}
	i.transients[t] = factory
	return nil
}

// GetInstance returns an instance of the given type. It returns the instance,
// a releaser function, and an error. The releaser function should be called
// when the instance is no longer needed.
func (i *Injector) GetInstance(t reflect.Type) (interface{}, func(), error) {
	if instance, ok := i.singletons[t]; ok {
		return instance, func() {}, nil
	}
	if factory, ok := i.transients[t]; ok {
		values := reflect.ValueOf(factory).Call(nil)
		i := values[0].Interface()
		// If the instance is an io.Closer, return a releaser function that
		// closes the instance.
		if _, ok := i.(io.Closer); ok {
			return i, func() { i.(io.Closer).Close() }, nil
		}
		return i, func() {}, nil
	}
	return nil, nil, ErrTypeNotRegistered
}

// Invoke calls the given function, passing the initial arguments and then
// injecting instances for the remaining arguments. It returns the results of
// the function and an error. Invoke handles collecting the releaser functions
// and calling them after the function has been called.
func (i *Injector) Invoke(fn interface{}, args ...interface{}) ([]interface{}, error) {
	f := reflect.TypeOf(fn)
	if f.Kind() != reflect.Func {
		return nil, ErrNotInvokable
	}
	// Collect the initial arguments.
	values := make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}
	// Collect the remaining arguments.
	for j := len(args); j < f.NumIn(); j++ {
		t := f.In(j)
		instance, releaser, err := i.GetInstance(t)
		if err != nil {
			return nil, err
		}
		values = append(values, reflect.ValueOf(instance))
		defer releaser()
	}
	// Call the function.
	results := reflect.ValueOf(fn).Call(values)
	// Collect the results.
	r := make([]interface{}, len(results))
	for i, result := range results {
		r[i] = result.Interface()
	}
	return r, nil
}
