package DependencyInjection

import (
	"reflect"
)

const TAG_NAME  = "autoinject"

type constructor func (*Container) interface{}

type Container struct {
	services map[string]interface{}
	serviceDefinitions map[string](constructor)
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
		serviceDefinitions: make(map[string](constructor)),
	}
}

func (this *Container) Register(key string, serviceConstructor constructor)  {
	this.serviceDefinitions[key] = serviceConstructor
}

func (this *Container) Get(key string) interface{} {
	if !this.has(key) {
		panic("Not implemented")
	}

	if _, ok := this.services[key]; !ok {
		serviceConstructor := this.serviceDefinitions[key]
		this.services[key] = serviceConstructor(this)
	}

	return this.services[key]
}

func (this *Container) AutoInject(object interface{}) interface{} {
	value := reflect.ValueOf(object).Elem()
	vType := reflect.TypeOf(object).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		vField := vType.Field(i)

		if !field.CanSet() {
			continue
		}

		if _, ok := vField.Tag.Lookup(TAG_NAME); !ok {
			continue
		}

		constructedService := this.Get(field.Type().Elem().Name())
		field.Set(reflect.ValueOf(constructedService))
	}

	return object
}

func (this *Container) has(key string) bool {
	_, ok := this.serviceDefinitions[key]
	return ok
}
