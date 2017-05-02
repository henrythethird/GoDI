package autoinject

import (
	"reflect"
	"fmt"
	"errors"
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

func (c *Container) Get(key string) (interface{}, error) {
	if !c.has(key) {
		return nil, errors.New(fmt.Sprintf("Unregistered service: \"%s\"", key))
	}

	if _, ok := c.services[key]; !ok {
		serviceConstructor := c.serviceDefinitions[key]
		service, err := c.AutoInject(serviceConstructor())

		if nil != err {
			return nil, err
		}

		c.services[key] = service
	}

	return c.services[key], nil
}

func (c *Container) GetParameter(key string) (interface{}, error) {
	val, ok := c.parameters[key]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Undefined parameter: \"%s\"", key))
	}

	return val, nil
}

func (c *Container) AddParameter(key string, value interface{}) (*Container) {
	c.parameters[key] = value

	return c
}

func (c *Container) AutoInject(object interface{}) (interface{}, error) {
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

		resolved, err := c.resolveTag(tagValue, field.Type())

		if nil != err {
			return nil, err
		}

		field.Set(reflect.ValueOf(resolved))
	}

	return object, nil
}

func (c *Container) resolveTag(tagValue string, fieldType reflect.Type) (interface{}, error) {
	if tagValue == "-" {
		if fieldType.Kind().String() != "ptr" {
			return nil, errors.New(fmt.Sprintf("Expected pointer type got: \"%s\"", fieldType))
		}

		return c.Get(fieldType.Elem().String())
	}

	return c.GetParameter(tagValue)
}

func (c *Container) has(key string) bool {
	_, ok := c.serviceDefinitions[key]
	return ok
}
