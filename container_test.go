package autoinject

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


type TestContainer_Service struct {

}

type TestContainer_InjectTest struct {
	Foo	string			`autoinject:"test.string.foo"`
	Service *TestContainer_Service	`autoinject:"-"`
}

type TestContainer_Invalid struct {
	Service TestContainer_Service	`autoinject:"-"`
}

type TestContainer_MultiLevel struct {
	Service *TestContainer_InjectTest	`autoinject:"-"`
}

func TestContainer_AddParameter(t *testing.T) {
	container := NewContainer()

	container.AddParameter("test.randomInt", 42)

	assert.Equal(t, 42, container.GetParameter("test.randomInt"))
}

func TestContainer_GetParameter(t *testing.T) {
	container := NewContainer()

	assert.Panics(t, func() {
		container.GetParameter("test.invalid")
	})
}

func TestContainer_Register(t *testing.T) {
	container := NewContainer()

	container.Register("service", func() interface{} {
		return &struct {}{}
	})

	assert.NotPanics(t, func() {
		container.Get("service")
	})
}

func TestContainer_Get(t *testing.T) {
	container := NewContainer()

	assert.Panics(t, func() {
		container.Get("invalid")
	})
}

func TestContainer_AutoInject(t *testing.T) {
	container := NewContainer()

	container.AddParameter("test.string.foo", "Foo is a string")
	container.Register("autoinject.TestContainer_Service", func() interface{} {
		return new(TestContainer_Service)
	})

	testObj := new(TestContainer_InjectTest)

	assert.NotPanics(t, func() {
		container.AutoInject(testObj)
	})
}

func TestContainer_AutoInject2(t *testing.T) {
	container := NewContainer()

	container.AddParameter("test.string.foo", "Foo is a string")
	container.Register("autoinject.TestContainer_Service", func() interface{} {
		return new(TestContainer_Service)
	})
	container.Register("autoinject.TestContainer_InjectTest", func() interface{} {
		return new(TestContainer_InjectTest)
	})

	testObj := new(TestContainer_MultiLevel)

	assert.NotPanics(t, func() {
		container.AutoInject(testObj)
	})
}

func TestContainer_AutoInject_PanicsOnNonPointer(t *testing.T) {
	container := NewContainer()

	container.Register("autoinject.TestContainer_Service", func() interface{} {
		return new(TestContainer_Service)
	})

	testObj := new(TestContainer_Invalid)

	assert.Panics(t, func() {
		container.AutoInject(testObj)
	})
}

