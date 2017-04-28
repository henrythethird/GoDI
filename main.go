package main

import (
	"fmt"
	. "./dependencyinjection"
)

type YetAnotherService struct {
	Age 		int
	Name		string			`autoinject:"service_param"`
}

type OtherService struct {
	Value int
}

type SomeService struct {
	OtherService 	*OtherService 		`autoinject:"-"`
	AnotherService 	*YetAnotherService	`autoinject:"-"`
	InternalValue	int
	Parameter 	string			`autoinject:"param_value"`
}

type ValueService struct {
	SomeService *SomeService                `autoinject:"-"`
}

func main() {
	container := NewContainer()

	container.
	AddParameter("param_value", "There once was a man").
		AddParameter("service_param", "My Name").
		Register("main.OtherService", func() interface{} {
			return &OtherService{Value: 100}
		}).
		Register("main.YetAnotherService", func() interface{} {
			return &YetAnotherService{Age: 200}
		}).
		Register("main.SomeService", func() interface{} {
			return &SomeService{}
		}).
		Register("main.ValueService", func() interface{} {
			return &ValueService{}
		})

	valueServiceInstance := container.Get("main.ValueService")
	someServiceInstance := container.Get("main.SomeService")

	someServiceInstance.(*SomeService).InternalValue = 55

	fmt.Println(
		someServiceInstance.(*SomeService).AnotherService.Name,
		valueServiceInstance.(*ValueService).SomeService.InternalValue,
		someServiceInstance,
	)
}
