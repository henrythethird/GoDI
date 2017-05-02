package autoinject

import (
	"reflect"
	"fmt"
)

const TAG_NAME = "autoinject"

type constructor func() interface{}

type Config map[string]interface{}
type ServiceMap map[string]interface{}
type ServiceDefinitionMap map[string](constructor)

type Container struct {
	services           ServiceMap
	serviceDefinitions ServiceDefinitionMap
	parameters         Config
}

func NewContainer() *Container {
	return &Container{
		services:           make(ServiceMap),
		serviceDefinitions: make(ServiceDefinitionMap),
		parameters:         make(Config),
	}
}

func (c *Container) Register(key string, serviceConstructor constructor) (*Container) {
	c.serviceDefinitions[key] = serviceConstructor

	return c
}

func (c *Container) Get(key string) interface{} {
	if !c.has(key) {
		panic(fmt.Sprintf("Unregistered service: \"%s\"", key))
	}

	if _, ok := c.services[key]; !ok {
		serviceConstructor := c.serviceDefinitions[key]
		c.services[key] = c.AutoInject(serviceConstructor())
	}

	return c.services[key]
}

func (c *Container) GetParameter(key string) interface{} {
	val, ok := c.parameters[key]

	if !ok {
		panic(fmt.Sprintf("Undefined parameter: \"%s\"", key))
	}

	return val
}

func (c *Container) AddParameter(key string, value interface{}) (*Container) {
	c.parameters[key] = value

	return c
}

func (c *Container) AutoInject(object interface{}) interface{} {
	value := reflect.ValueOf(object).Elem()
	vType := reflect.TypeOf(object).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := vType.Field(i).Tag

		if !field.CanSet() {
			continue
		}

		tagValue, ok := tag.Lookup(TAG_NAME)

		if !ok {
			continue
		}

		field.Set(reflect.ValueOf(
			c.resolveTag(tagValue, field.Type()),
		))
	}

	return object
}

func (c *Container) resolveTag(tagValue string, fieldType reflect.Type) interface{} {
	if tagValue == "-" {
		if fieldType.Kind().String() != "ptr" {
			panic(fmt.Sprintf("Expected pointer type got: \"%s\"", fieldType))
		}

		return c.Get(fieldType.Elem().String())
	}

	return c.GetParameter(tagValue)
}

func (c *Container) has(key string) bool {
	_, ok := c.serviceDefinitions[key]
	return ok
}
