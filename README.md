# Galidator

Galidator is a package validator which implements struct or map field validations.

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

### Example-1

```go
import (
	"fmt"

	"github.com/golodash/galidator"
)

func main() {
	g := galidator.New()
	schema := g.Generate(galidator.Items{
		"number": g.Item().Int(),
		"id":     g.Item().Float(),
	}, nil)

	output := map[string]string{
		"number": "15", // Valid int data
		"id":     "125s", // Invalid float data
	}
	fmt.Println(schema.Validate(output))
}
```

Output:
```
map[id:[id is not float]]
```

### Example-2

```go
import (
	"fmt"

	"github.com/golodash/galidator"
)

type testStruct struct {
	Number int
	ID     string
}

func main() {
	g := galidator.New()
	schema := g.Generate(galidator.Items{
		"Number": g.Item().Int(),
		"ID":     g.Item().Float(),
	}, nil)

	output := testStruct{
		Number: 15,        // Valid int data
		ID:     "125.23d", // Invalid float data
	}
	fmt.Println(schema.Validate(output))
}
```

Output:
```
map[ID:[ID is not float]]
```
