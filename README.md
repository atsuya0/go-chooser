# go-choice

chooser.Run([]string) -> string

```go
package main

import (
	"fmt"
	"github.com/tayusa/go-choice"
)

func main() {
	chooser := choice.NewChooser([]string{"test9", "te11", "jfdls", "fdsaf", "daj", "fdsie", "feafii", "fdiaoeioa", "feiaofjl"})
	fmt.Print(chooser.Run())
}
```
