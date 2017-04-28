package DependencyInjection

import (
	"reflect"
	//"fmt"
)

const TAG_NAME  = "autoinject"

type constructor func (*Container) interface{}

type Config map[string]interface{}
type ServiceMap map[string]interface{}
type ServiceDefinitionMap map[string](constructor)


type Container struct {
	services ServiceMap
	serviceDefinitions ServiceDefinitionMap
	parameters Config
}

func NewContainer() *Container {
	return &Container{
		services: make(ServiceMap),
		serviceDefinitions: make(ServiceDefinitionMap),
		parameters: make(Config),
	}
}

func (this *Container) Register(key string, serviceConstructor constructor) (*Container) {
	this.serviceDefinitions[key] = serviceConstructor

	return this
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

func (this *Container) GetParameter(key string) interface{} {
	val, ok := this.parameters[key]

	if !ok {
		panic("Parameter not in list")
	}

	return val
}

func (this *Container) AddParameter(key string, value interface{}) (*Container) {
	this.parameters[key] = value

	return this
}

func (this *Container) AutoInject(object interface{}) interface{} {
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

		var constructedService interface{}

		/* automatically resolve type */
		if tagValue == "-" {
			constructedService = this.Get(field.Type().Elem().Name())
		} else {
			constructedService = this.GetParameter(tagValue)
		}


		field.Set(reflect.ValueOf(constructedService))
	}

	return object
}

func (this *Container) has(key string) bool {
	_, ok := this.serviceDefinitions[key]
	return ok
}
