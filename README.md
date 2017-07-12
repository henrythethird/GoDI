# AutoInject

AutoInject is a lightweight golang dependency injection framework

**Build Status:** [![Build Status](https://travis-ci.org/owtorg/autoinject.svg?branch=master)](https://travis-ci.org/owtorg/autoinject)

### Basic Usage - AutoInjection

```go
import "github.com/henrythethird/autoinject"
```

Given a service you'd like to inject, register it with the fully qualified struct name to the container:
```go
type Name struct {}

container := NewContainer()
container.Register("mypackage.Name", func() interface{} {
    return new(Name)
})
```

Now we need to set up the receiving container with the necessary tags. Note: Make sure the service is a pointer
```go
type SomeStruct struct {
    Service *Name `autoinject="-"`
}
```

Once we have done this, we can tell the container to automatically inject the described services
```go
myService := new(SomeStruct)

container.AutoInject(myService)
```

### Basic Usage - Parameter Injection
The injectable named parameters can be of any type. Here a basic example

```go
container := NewContainer()
container.AddParameter("theAnswerToEverything", 42)

// We can simply get the parameter out again:
container.GetParameter("theAnswerToEverything")
/* This requires the parameter to be cast, hence is a bit ugly.
 * Another solution is to create a struct and automatically inject the named parameter:
 */

type SomeStruct struct {
    Answer int `autoinject="theAnswerToEverything"`
}

myStruct := new(SomeStruct)
container.AutoInject(myStruct)
// Now the parameter is injected
```

