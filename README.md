# Galidator

Galidator is a package validator which implements struct or map validation.\
At this point, package is at beta phase and I'm still testing the core of galidator.

## Installation

Just use `go get` for installation:

```
go get github.com/golodash/galidator
```

And then just import the package into your own code.

```go
import (
	"github.com/golodash/galidator"
)
```

## Example Usage

### Simple Usage(Register a User)

Lets validate a register form:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Username": g.R("username").Required().Min(3).Max(32),
		"Password": g.R("password").Required().Password(),
		"Email":    g.R("email").Required().Email(),
	}))

	userInput := map[string]string{
		"username": "DoctorMK",
		"password": "123456789",
		"email":    "DoctorMK@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[password:[password must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character]]
false
```

We can even validate a struct and get the same result:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Person struct {
	Username string
	Password string
	Email    string
}

func main() {
	g := galidator.New()
	validator := g.Validator(galidator.Rules{
		"Username": g.R("username").Required().Min(3).Max(32),
		"Password": g.R("password").Required().Min(5).Password(),
		"Email":    g.R("email").Required().Email(),
	})

	userInput := Person{
		Username: "DoctorMK",
		Password: "123456789",
		Email:    "DoctorMK@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Or we can even create our validator by defining struct tags.

**Note**: This is not a complete feature and works just in a single layer.

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Person struct {
	Username string `json:"username" g:"required,min=3,max=32"`
	Password string `json:"password" g:"required,min=5,password"`
	Email    string `json:"email" g:"required,email"`
}

func main() {
	g := galidator.New()
	validator := g.ValidatorFromStruct(Person{})

	userInput := Person{
		Username: "DoctorMK",
		Password: "123456789",
		Email:    "DoctorMK@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

### Lists

Lets assume we need to receive a list of orders, in this case:

**Note**: This example won't work in a tag based validator. (Under Development...)

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Order struct {
	ObjectName string
	Amount     int
	Price      float64
}

func main() {
	g := galidator.New()
	// Implementing a slice of Orders = []Order
	ordersValidator := g.Validator(g.R().Children(
		g.R().Complex(galidator.Rules{
			"ObjectName": g.R("object_name").Min(3),
			"Amount":     g.R("Amount").Min(1).Max(10),
			"Price":      g.R("Price").Min(1).Max(500),
		}),
	))

	userInput := []Order{
		{
			ObjectName: "e",
			Amount:     3,
			Price:      2,
		}, {
			ObjectName: "Bathroom cleaner",
			Amount:     -2,
			Price:      5,
		},
	}

	errors := ordersValidator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[0:map[object_name:[object_name's length must be higher equal to 3]] 1:map[Amount:[Amount's length must be higher equal to 1]]]
false
```

### Custom Validators

Lets go back to previous example of sign up a user.\
In this example, we need to check if the user is inside the database but in this example we just use a list of names instead of a database.

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Person struct {
	Username string
	Password string
	Email    string
}

var usernames = []string{
	"erfan",
	"ali",
	"asghar",
	"mohammad",
	"james",
}

func usernameDuplicateChecker(input interface{}) bool {
	for _, name := range usernames {
		if name == input.(string) {
			return false
		}
	}
	return true
}

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Username": g.R("username").Required().Min(3).Max(32).Custom(galidator.Validators{"DuplicateUsername": usernameDuplicateChecker}),
		"Password": g.R("password").Required().Password(),
		"Email":    g.R("email").Required().Email(),
	}), galidator.Messages{
		"DuplicateUsername": "$value already exists",
	})

	userInput := map[string]string{
		"username": "mohammad",
		"password": "123456789",
		"email":    "james@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[Password:[Password must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character] Username:[mohammad already exists]]
false
```

### OR, XOR Rule

Or operator checks if at least one of the passed rules pass.\
Note: XOR usage is the same but results are based on XOR but not based on OR operation.

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Request struct {
	Username string
	Password string
}

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Username": g.R().Required().OR(g.R().Email(), g.R().Phone()),
		"Password": g.R().Password().Min(8).Max(100),
	}), galidator.Messages{
		"OR": "$field should be a valid email or phone number",
	})

	userInput := &Request{
		Username: "not an email or phone number",
		Password: "12345678Aa!",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[Username:[ruleSets in Username did not pass based on or logic]]
false
```

### Choices

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Request struct {
	Username string
	Password string
	Method   string
}

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Username": g.R("username").Required(),
		"Password": g.R("password").Password().Min(8).Max(100),
		"Method":   g.R("method").Required().Choices([]string{"session", "jwt"}),
	}))

	userInput := &Request{
		Username: "randomEmail@gmail.com",
		Password: "12345678Aa!",
		Method:   "invalid method",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}

```

Output:
```
map[method:[invalid method does not include in allowed choices: [session, jwt]]]
false
```

### WhenExistAll

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Request struct {
	Option1 string
	Option2 string
	Option3 string
}

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Option1": g.R("option_1").WhenExistAll("Option2", "Option3"),
		"Option2": g.R("option_2"),
		"Option3": g.R("option_3"),
	}), galidator.Messages{
		"OR": "$field should be a valid email or phone number",
	})

	userInput := &Request{
		Option1: "",
		Option2: "data",
		Option3: "data",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}

```

Output:
```
map[option_1:[option_1 is required because all of [Option2, Option3] fields are not nil, empty or zero(0, "", '')]]
false
```

### WhenExistOne

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Request struct {
	Option1 string
	Option2 string
	Option3 string
}

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Complex(galidator.Rules{
		"Option1": g.R("option_1").WhenExistOne("Option2", "Option3"),
		"Option2": g.R("option_2"),
		"Option3": g.R("option_3"),
	}), galidator.Messages{
		"OR": "$field should be a valid email or phone number",
	})

	userInput := &Request{
		Option1: "",
		Option2: "",
		Option3: "data",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[option_1:[option_1 is required because at least one of [Option2, Option3] fields are not nil, empty or zero(0, "", '')]]
false
```
