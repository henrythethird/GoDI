package main

import (
	"fmt"
	. "./DependencyInjection"
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
}

type ValueService struct {
	SomeService	*SomeService	`autoinject:"-"`
}

func main() {
	container := NewContainer()

	container.Register("OtherService", func(container *Container) interface{} {
		return &OtherService{Value:100}
	})

	container.Register("YetAnotherService", func(container *Container) interface{} {
		return &YetAnotherService{Age:200}
	})

	container.Register("SomeService", func(container *Container) interface{} {
		return container.AutoInject(&SomeService{})
	})

	container.Register("ValueService", func(container *Container) interface{} {
		return container.AutoInject(&ValueService{})
	})

	valueServiceInstance := container.Get("ValueService")
	someServiceInstance := container.Get("SomeService")

	someServiceInstance.(*SomeService).InternalValue = 55

	fmt.Println(someServiceInstance.(*SomeService).InternalValue, valueServiceInstance.(*ValueService).SomeService.InternalValue)
}

