# Galidator

Galidator provides general use case for validation purpose.\
Just simply create a validator and validate your data with it.\
Either it returns `nil` which means that data is valid and your rules all passed or not which means that there
is/are problem/problems with passed data and validation has failed.

## Installation

Just use `go get` for installation:

```
go get github.com/golodash/galidator
```

And then just import the package into your own code.

```go
import "github.com/golodash/galidator"
```

## Generator

### What is a Generator?

Generator is a validator generator which creates them under some common circumstances.\
For example, with `generator.CustomMessages` method you can change default error messages of rules for every other validator that is gonna get created with this generator.

With this mechanism, you can change some default behavior even before attempting to create your validator.

### How to Create A Generator?

```go
var g = galidator.NewGenerator()
or
var g = galidator.New()
or
var g = galidator.G()
```

All three does the same.

## Validator

### What is a Validator?

Validator is a `galidator.Validator` interface that gets created under some circumstances that its generator defines + common environmental variables of galidator defines + your specified rules.

### How to Create a Validator?

To create a validator, first you need to create a unique generator instance and then use generator to call a method to create your validator.

1. `Validator(input interface{}, messages ...Messages) Validator`:
   1. `input` can be a `ruleSet`. (which gets created by `generator.R() or generator.RuleSet()`)\
	Example down here accepts just an email string:
	```go
	var g = galidator.G()
	var validator = g.Validator(g.R().Email())
	```
   2. `input` can be a `struct instance` with tags that define rules on every field.\
    Example down here accepts a map or a struct which has a `Email` key and a value which is a valid email:
	```go
	type user struct {
		Email string `galidator:"email"` // instead of `galidator:"email"`, `g:"email"` can be used
	}
	var g = galidator.G()
	var validator = g.Validator(user{})
	```
   3. `messages` input type is obvious and is used to replace common error messages on rules for just this validator.

2. `ComplexValidator(rules Rules, messages ...Messages) Validator`:
   - Generates a struct/map validator.\
	Mostly used for complex scenarios that start with a struct or map view.
	```go
	var g = galidator.G()
	var validator = g.ComplexValidator(galidator.Rules{
		"Name":        g.RuleSet("name").Required(),
		"Description": g.RuleSet("description").Required(),
	})
	```

### How to Validate?

Simply just create your validator with one of the top discussed methods and then:

```go
var g = galidator.G()
var emailValidator = g.Validator(g.R().Email())

func main() {
	input1 := "valid@email.com"
	input2 := "invalidEmail.com"

	output1 := emailValidator.Validate(input1)
	output2 := emailValidator.Validate(input2)

	fmt.Println(output1)
	fmt.Println(output2)
}
```

output:
```
<nil>
[not a valid email address]
```

And that's it, just to get better, see more examples down here.

# Just For [Gin](https://github.com/gin-gonic/gin) Users

## 1. Use Galidator Just to Customize [Gin](https://github.com/gin-gonic/gin)'s Bind Method Error Outputs

You can choose not to use galidator and it's validation process but instead use `Bind` method of [Gin](https://github.com/gin-gonic/gin) or other acronym's for `Bind` like: `BindJson` for validation process and just use galidator to change output error messages of it.

Example of using galidator inside a gin project:

```go
type login struct {
	Username string `json:"username" binding:"required" required:"$field is required"`
	Password string `json:"password"`
}

var (
	g = galidator.G()
	validator = g.Validator(login{})
)

func test(c *gin.Context) {
	req := &login{}

	// Parse json
	if err := c.BindJSON(req); err != nil {
		c.JSON(400, gin.H{
			// This is the part which generates that output
			"message": validator.DecryptErrors(err),
		})
		return
	}

	c.JSON(200, gin.H{
		"good": 200,
	})
}

func main() {
	r := gin.Default()
	r.POST("/", test)
	r.Run("127.0.0.1:3000")
}
```

If you don't send `username` or send it empty in json request body, this message returns:
```
{"message":{"username":"username is required"}}
```

**Note**: In cases which there is a conflict of names for rule messages(like `json`, `json` tag is used for output in json and if you want to add error message of json for it, you can't use `json` tag, so the solution is to use `_json` tag)

Example:

```go
type login struct {
	Field string `json:"field" binding:"required,json" _json:"$field is not json"`
}
```

## 2. Translate Error Output to Different Languages in [Gin]((https://github.com/gin-gonic/gin))

If you need to translate output error messages for different languages in a gin project, use this template:

```go
type login struct {
	Username string `json:"username" g:"required" required:"$field is required"`
	Password string `json:"password"`
}

var (
	g            = galidator.New()
	validator    = g.Validator(login{})
	// Persian Language Dictionary
	faDictionary = map[string]string{
		"$field is required": "$field نمیتواند خالی باشد",
	}
)

// Persian Language Translator
func PersianTranslator(input string) string {
	if translated, ok := faDictionary[input]; ok {
		return translated
	}
	return input
}

// Middleware that assigns a translator requested by user
func customizeTranslator(c *gin.Context) {
	languageCode := c.GetHeader("Accept-Language")
	if languageCode == "fa" {
		c.Set("translator", PersianTranslator)
	} else {
		c.Set("translator", func(input string) string { return input })
	}
	c.Next()
}

// Main Handler
func loginHandler(c *gin.Context) {
	req := &login{}
	translator := c.MustGet("translator").(func(string) string)

	// Parse json
	if err := c.BindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"message": "bad json",
		})
		return
	}

	// Validation
	if errors := validator.Validate(req, translator); errors != nil {
		c.JSON(400, gin.H{
			"errors":  errors,
			"message": "bad inputs",
		})
		return
	}

	c.JSON(200, gin.H{
		"good": 200,
	})
}

func main() {
	r := gin.Default()
	groupWithMiddleware := r.Group("/", customizeTranslator)
	groupWithMiddleware.POST("/", loginHandler)
	r.Run("127.0.0.1:3000")
}
```

Now if you make a post request to http://127.0.0.1:3000 url when having `Accept-Language` header with `fa` value assigned to it, and in request body do not specify username field or specify it's value as an empty string, this will be the output:

```
{"errors":{"username":["username نمیتواند خالی باشد"]},"message": "bad inputs"}
```

## 3. Making Ease of PATCH method with Galidator in [Gin]((https://github.com/gin-gonic/gin))

1. Make all your field types pointer. (string -> *string)
2. Use `SetDefaultOnNil` method which is accessible from a `Validator` instance.
   - Have in mind to pass pointer to a struct variable into `SetDefaultOnNil` method.
3. Done... items are nil if user did not send them to api and will be filled with default values which programmer passed them.

```go
type article struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

var (
	title    = "This is first"
	content  = "This is the content"
	articles = []article{
		{Title: &title, Content: &content},
	}
	g         = galidator.G()
	validator = g.Validator(article{})
)

func patchArticle(c *gin.Context) {
	req := &article{}
	id, _ := strconv.Atoi(c.Param("id"))
	defaults := articles[id]

	c.BindJSON(req)
	if err := validator.Validate(req); err == nil {
		// This is the part to set default value for nil fields
		validator.SetDefaultOnNil(req, defaults)
		// This is update action
		articles[id] = *req
	} else {
		c.JSON(400, gin.H{
			"message": "error in validation",
		})
	}

	c.JSON(200, gin.H{
		"good": 200,
		"data": req,
	})
}

func main() {
	r := gin.Default()
	r.PATCH("/comments/:id", patchArticle)
	r.Run("127.0.0.1:3000")
}
```

# Examples

## Simple Usage(Register a User)

Lets validate a register form:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.ComplexValidator(galidator.Rules{
		"Username": g.R("username").Required().Min(3).Max(32),
		"Password": g.R("password").Required().Password(),
		"Email":    g.R("email").Required().Email(),
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
map[password:[password must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character]]
false
```

We can even validate a struct by the same validator and get the same result:

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

Or we can even create the same validator by defining some struct tags.

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
	validator := g.Validator(Person{})

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

## Receive a list of users

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type Person struct {
	Username string `json:"username" g:"required,min=3,max=32"`
	Password string `json:"password" g:"required,min=5,password" password:"$field failed"`
	Email    string `json:"email" g:"required,email"`
}

func main() {
	g := galidator.New()
	validator := g.Validator([]Person{})

	userInput := []Person{
		{
			Username: "DoctorMK",
			Password: "123456789",
			Email:    "DoctorMK@gmail.com",
		},
		{
			Username: "Asghar",
			Password: "123456789mH!@",
			Email:    "Doctors@gmail.com",
		},
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[0:map[password:[password failed]]]
false
```

We can create the same validator without tags too:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Children(
		g.R().Complex(galidator.Rules{
			"Username": g.R("username").Required().Min(3).Max(32),
			"Password": g.R("password").Required().Password().SpecificMessages(galidator.Messages{"password": "$field failed"}),
			"Email":    g.R("email").Required().Email(),
		})))

	userInput := []map[string]string{
		{
			"username": "DoctorMK",
			"password": "123456789",
			"email":    "DoctorMK@gmail.com",
		},
		{
			"username": "Asghar",
			"password": "123456789mH!@",
			"email":    "Doctors@gmail.com",
		},
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

## OR

In this example, input has to be either an email address or just a string longer equal to 5 characters or both.

This example OR operator in struct tags can be used like: `g:"required,or=email|string+min=5"`

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Required().OR(g.R().Email(), g.R().String().Min(5)))

	input := "m@g.com"
	errors := validator.Validate(input)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
<nil>
true
```

## XOR

In this example, input has to be either an email address or phone number.

This example XOR operator in struct tags can be used like: `g:"required,xor=email|phone"`

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().Required().XOR(g.R().Email(), g.R().Phone()))

	input := "m@g.com"
	errors := validator.Validate(input)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
<nil>
true
```

## WhenExistAll - WhenExistOne

In this example, if two other struct fields(`Username` and `Password`) are not empty,
nil or zero, field will act as a required field and all of its rules will get checked.\
Otherwise, if empty, nil or zero value get passed, because by default fields are
optional, it does not check other defined rules and assume it passed.\
You can use this in tags like: `g:"when_exist_all=Username&Password,string" when_exist_all:"when_exist_all failed"`

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	v := g.ComplexValidator(galidator.Rules{
		"Username": g.R("username").String(),
		"Password": g.R("password").String(),
		"Data":     g.R("data").WhenExistAll("Username", "Password").String().SpecificMessages(galidator.Messages{"when_exist_all": "when_exist_all failed"}),
	})

	errors := v.Validate(map[string]string{
		"Username": "username",
		"Password": "password",
		"Data":     "",
	})

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[data:[when_exist_all failed]]
false
```

## Custom Validator

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
	"github.com/golodash/godash/slices"
)

type Person struct {
	Username string `json:"username" g:"required,min=3,max=32,duplicate_check" duplicate_check:"$field is duplicate"`
	Password string `json:"password" g:"required,min=5,password"`
	Email    string `json:"email" g:"required,email"`
}

var users = []string{
	"ali",
	"james",
	"john",
}

func duplicate_check(input interface{}) bool {
	return slices.FindIndex(users, input) == -1
}

func main() {
	g := galidator.New().CustomValidators(galidator.Validators{"duplicate_check": duplicate_check})
	validator := g.Validator(Person{})

	userInput := Person{
		Username: "ali",
		Password: "123456789mH!",
		Email:    "DoctorMK@gmail.com",
	}

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
map[username:[username is duplicate]]
false
```

## LenRange

LenRange can be used in a struct tag like: `g:"len_range=3&5" len_range="len_range failed"`

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	validator := g.Validator(g.R().LenRange(3, 5).SpecificMessages(galidator.Messages{"len_range": "len_range failed"}))

	userInput := 3

	errors := validator.Validate(userInput)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

Output:
```
[len_range failed]
false
```

## Changing Default Error Messages

1. Changing default error messages in generator layer:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New().CustomMessages(galidator.Messages{
		"string": "$value is not string",
	})
	validator := g.Validator(g.R().String())

	errors := validator.Validate(1)

	fmt.Println(errors)
	fmt.Println(errors == nil)
}
```

output:
```
[1 is not string]
```

2. Changing default error messages in validator layer (This layer overrides generator error messages-if same value that is defined in other layers gets defined in this layer too):

```go
g := galidator.New().CustomMessages(galidator.Messages{
   "string": "$value is not string",
})
validator := g.Validator(g.R().String(), galidator.Messages{
   "string": "not",
})
```

output:
```
[not]
```

3. Changing default error messages in ruleSet layer (This layer overrides generator and validator error messages-if same value that is defined in other layers gets defined in this layer too):

```go
g := galidator.New().CustomMessages(galidator.Messages{
   "string": "$value is not string",
})
validator := g.Validator(g.R().String().SpecificMessages(galidator.Messages{
		"string": "not valid",
	}), galidator.Messages{
	"string": "not",
})
```

output:
```
[not valid]
```

## Defining `ruleSet` for Children of a Slice in Struct Tags

If you need to define a rule for children of a slice in struct tags, you should use
some proper prefix for those rules like: `c.` or `child.`\
And have in mind that with adding two or more of these prefixes, you keep digging in
deeper layers. like: `child.child.child.min` or `c.c.c.min` or `c.child.c.min` or... means go deep three slices and add `min` rule to children of the last slice.

example:

```go
package main

import (
	"fmt"

	"github.com/golodash/galidator"
)

type numbers struct {
	Numbers []int `g:"c.min=1,c.max=5" c.max:"$value is not <= $max" c.min:"$value is not >= $min"`
}

func main() {
	g := galidator.New()
	validator := g.Validator(numbers{})

	fmt.Println(validator.Validate(numbers{
		Numbers: []int{
			1,
			0,
			5,
			35,
		},
	}))
}
```

output:
```
map[Numbers:map[1:[0 is not >= 1] 3:[35 is not <= 5]]]
```

## Translator

When calling `Validator.Validate` method with your data, you can pass a translator function to translate output of error messages to your desired language.

For example:

```go
var (
	g = galidator.G()
	validator = g.Validator(g.R().Required())
	translates = map[string]string{
    	"required": "this is required and it is translated",
	}
)

func translator(input string) string {
	if out, ok := translates[input]; ok {
		return out
	}
	return input
}

func main() {
	fmt.Println(validator.Validate(nil, translator))
}
```

output:
```go
[this is required and it is translated]
```

# Star History

[![Star History Chart](https://api.star-history.com/svg?repos=golodash/galidator&type=Date)](https://star-history.com/#golodash/galidator&Date)
