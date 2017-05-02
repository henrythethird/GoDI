package autoinject

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


type TestContainer_Service struct {

}

type TestContainer_InjectTest struct {
	private		string
	OtherTag	string			`json:"other_tag"`
	Foo		string			`autoinject:"test.string.foo"`
	Service 	*TestContainer_Service	`autoinject:"-"`
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

	param, _ := container.GetParameter("test.randomInt")
	assert.Equal(t, 42, param)
}

func TestContainer_GetParameter(t *testing.T) {
	container := NewContainer()

	_, err := container.GetParameter("test.invalid")
	assert.Error(t, err)
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

	_, err := container.Get("invalid")
	assert.Error(t, err)
}

func TestContainer_AutoInject(t *testing.T) {
	container := NewContainer()

	container.AddParameter("test.string.foo", "Foo is a string")
	container.Register("autoinject.TestContainer_Service", func() interface{} {
		return new(TestContainer_Service)
	})

	testObj := new(TestContainer_InjectTest)

	_, err := container.AutoInject(testObj)

	assert.NoError(t, err)
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

	_, err := container.AutoInject(testObj)
	assert.NoError(t, err)
}

func TestContainer_AutoInject_ErrorOnNonPointer(t *testing.T) {
	container := NewContainer()

	container.Register("autoinject.TestContainer_Service", func() interface{} {
		return new(TestContainer_Service)
	})

	testObj := new(TestContainer_Invalid)

	_, err := container.AutoInject(testObj)
	assert.Error(t, err)
}

