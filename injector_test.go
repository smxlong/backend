package backend

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_that_Injector_can_register_a_singleton(t *testing.T) {
	i := NewInjector()
	instance := "test"
	err := i.AddSingleton(instance)
	require.Nil(t, err)
}

func Test_that_Injector_returns_error_when_registering_same_type_as_singleton_twice(t *testing.T) {
	i := NewInjector()
	instance := "test"
	err := i.AddSingleton(instance)
	require.Nil(t, err)
	err = i.AddSingleton(instance)
	require.Equal(t, ErrTypeAlreadyRegistered, err)
}

func Test_that_Injector_can_register_a_transient(t *testing.T) {
	i := NewInjector()
	factory := func() string { return "test" }
	err := i.AddTransient(factory)
	require.Nil(t, err)
}

func Test_that_Injector_returns_error_when_registering_same_type_as_transient_twice(t *testing.T) {
	i := NewInjector()
	factory1 := func() string { return "test" }
	err := i.AddTransient(factory1)
	require.Nil(t, err)
	factory2 := func() string { return "test2" }
	err = i.AddTransient(factory2)
	require.Equal(t, ErrTypeAlreadyRegistered, err)
}

func Test_that_Injector_returns_error_when_registering_invalid_factory(t *testing.T) {
	i := NewInjector()
	factory := "test"
	err := i.AddTransient(factory)
	require.Equal(t, ErrInvalidFactory, err)
}

func Test_that_Injector_returns_error_when_registering_same_type_as_singleton_and_transient(t *testing.T) {
	i := NewInjector()
	instance := "test"
	err := i.AddSingleton(instance)
	require.Nil(t, err)
	factory := func() string { return "test" }
	err = i.AddTransient(factory)
	require.Equal(t, ErrTypeAlreadyRegistered, err)
}

func Test_that_Injector_returns_error_when_getting_instance_of_unregistered_type(t *testing.T) {
	i := NewInjector()
	ty := reflect.TypeOf("test")
	_, _, err := i.GetInstance(ty)
	require.Equal(t, ErrTypeNotRegistered, err)
}

type simpleObject struct {
	name  string
	value int
}

func Test_that_Injector_returns_singleton_instance(t *testing.T) {
	i := NewInjector()
	instance := simpleObject{name: "test", value: 42}
	err := i.AddSingleton(instance)
	require.Nil(t, err)
	ty := reflect.TypeOf(instance)
	instance2, _, err := i.GetInstance(ty)
	require.Nil(t, err)
	require.True(t, instance == instance2) // Don't use Equal because it performs deep comparison
}

func Test_that_Injector_returns_transient_instance(t *testing.T) {
	i := NewInjector()
	factory := func() *simpleObject { return &simpleObject{name: "test", value: 42} }
	err := i.AddTransient(factory)
	require.Nil(t, err)
	ty := reflect.TypeOf(&simpleObject{})
	instance, _, err := i.GetInstance(ty)
	require.Nil(t, err)
	require.Equal(t, &simpleObject{name: "test", value: 42}, instance)
}

func Test_that_Injector_returns_independent_transient_instances(t *testing.T) {
	i := NewInjector()
	factory := func() *simpleObject { return &simpleObject{name: "test", value: 42} }
	err := i.AddTransient(factory)
	require.Nil(t, err)
	ty := reflect.TypeOf(&simpleObject{})
	instance1, _, err := i.GetInstance(ty)
	require.Nil(t, err)
	instance2, _, err := i.GetInstance(ty)
	require.Nil(t, err)
	require.True(t, instance1 != instance2) // Don't use NotEqual because it performs deep comparison
}

func Test_that_Invoke_passes_initial_arguments(t *testing.T) {
	i := NewInjector()
	f := func(a, b int) int { return a + b }
	result, err := i.Invoke(f, 1, 2)
	require.Nil(t, err)
	require.Equal(t, []interface{}{3}, result)
}

func Test_that_Invoke_returns_error_when_function_is_not_invokable(t *testing.T) {
	i := NewInjector()
	f := "test"
	_, err := i.Invoke(f)
	require.Equal(t, ErrNotInvokable, err)
}

func Test_that_Invoke_injects_instances(t *testing.T) {
	i := NewInjector()
	instance := simpleObject{name: "test", value: 42}
	err := i.AddSingleton(instance)
	require.Nil(t, err)
	f := func(o simpleObject) int { return o.value }
	result, err := i.Invoke(f)
	require.Nil(t, err)
	require.Equal(t, []interface{}{42}, result)
}

func Test_that_Invoke_injects_instances_and_passes_initial_arguments(t *testing.T) {
	i := NewInjector()
	instance := simpleObject{name: "test", value: 42}
	err := i.AddSingleton(instance)
	require.Nil(t, err)
	f := func(a int, o simpleObject) int { return a + o.value }
	result, err := i.Invoke(f, 1)
	require.Nil(t, err)
	require.Equal(t, []interface{}{43}, result)
}

type closableObject struct {
	closed *bool
}

func (c *closableObject) Close() error {
	*c.closed = true
	return nil
}

func Test_that_Invoke_calls_Close_on_transient_instance(t *testing.T) {
	i := NewInjector()
	closed := false
	factory := func() *closableObject { return &closableObject{closed: &closed} }
	err := i.AddTransient(factory)
	require.Nil(t, err)
	f := func(o *closableObject) {}
	_, err = i.Invoke(f)
	require.Nil(t, err)
	require.True(t, closed)
}

func Test_that_Invoke_returns_an_error_if_GetInstance_returns_an_error(t *testing.T) {
	i := NewInjector()
	f := func(o simpleObject) {}
	_, err := i.Invoke(f)
	require.Equal(t, ErrTypeNotRegistered, err)
}
