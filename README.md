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
import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(galidator.Rules{
		"Username": g.RuleSet().Required().Min(3).Max(32),
		"Password": g.RuleSet().Required().Password(),
		"Email":    g.RuleSet().Required().Email(),
	})

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
map[Password:[password must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character]]
false
```

We can even validate a struct and get the same result:

```go
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
		"Username": g.RuleSet().Required().Min(3).Max(32),
		"Password": g.RuleSet().Required().Min(5).Password(),
		"Email":    g.RuleSet().Required().Email(),
	})

	userInput := Person{
		Username: "DoctorMK",
		Password: "123456789",
		Email:    "DoctorMK@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == g.RuleSet().NonNil())
}
```

**Note**: Pay attention that you need to match exact key names of Person with keys of Rules in Validator when validator creation is happening.

### Lists

Lets assume we need to receive a list of orders, in this case:

```go
import (
	"fmt"

	"github.com/golodash/galidator"
)

type Order struct {
	ObjectName string
	Amount     int
	Price      float64
}

type UserInput struct {
	Orders []Order
}

func main() {
	g := galidator.New()
	ordersValidator := g.Validator(galidator.Rules{
		"Orders": g.RuleSet().Slice().Complex(g.Validator(galidator.Rules{
			"ObjectName": g.RuleSet().Required().Min(3),
			"Amount":     g.RuleSet().Min(1).Max(10),
			"Price":      g.RuleSet().Required().Min(1).Max(500),
		}),
		),
	})

	userInput := UserInput{
		Orders: []Order{{
			ObjectName: "e",
			Amount:     3,
			Price:      2,
		}, {
			ObjectName: "Bathroom cleaner",
			Amount:     -2,
			Price:      5,
		}},
	}

	errors := ordersValidator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == g.RuleSet().NonNil())
}
```

Output:
```
map[Orders:map[0:map[ObjectName:[object_name's length must be higher equal to 3]] 1:map[Amount:[amount's length must be higher equal to 1]]]]
false
```

### Custom Validators

Lets go back to previous example of sign up a user.\
In this example, we need to check if the user is inside the database but in this example we just use a list of names instead of a database.

```go
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
	validator := g.Validator(galidator.Rules{
		"Username": g.RuleSet().Required().Min(3).Max(32).Custom(galidator.Validators{"DuplicateUsername": usernameDuplicateChecker}),
		"Password": g.RuleSet().Required().Password(),
		"Email":    g.RuleSet().Required().Email(),
	}, galidator.Messages{
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
map[Password:[password must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character] Username:[mohammad already exists]]
false
```
