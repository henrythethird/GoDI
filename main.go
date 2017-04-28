package main

import (
	"fmt"
	. "./dependencyinjection"
)

type YetAnotherService struct {
	Age 		int
}

type OtherService struct {
	Value 		int
}

type SomeService struct {
	OtherService 	*OtherService 		`autoinject:"-"`
	AnotherService 	*YetAnotherService	`autoinject:"-"`
	InternalValue	int
	Parameter 	string			`autoinject:"param_value"`
}

type ValueService struct {
	SomeService	*SomeService	`autoinject:"-"`
}

func main() {
	container := NewContainer()

	container.
		AddParameter("param_value", "There once was a man").
		Register("OtherService", func(container *Container) interface{} {
			return &OtherService{Value:100}
		}).
		Register("YetAnotherService", func(container *Container) interface{} {
			return &YetAnotherService{Age:200}
		}).
		Register("SomeService", func(container *Container) interface{} {
			return container.AutoInject(&SomeService{})
		}).
		Register("ValueService", func(container *Container) interface{} {
			return container.AutoInject(&ValueService{})
		})

	valueServiceInstance := container.Get("ValueService")
	someServiceInstance := container.Get("SomeService")

	someServiceInstance.(*SomeService).InternalValue = 55

	fmt.Println(
		someServiceInstance.(*SomeService).InternalValue,
		valueServiceInstance.(*ValueService).SomeService.InternalValue,
		someServiceInstance,
	)
}

