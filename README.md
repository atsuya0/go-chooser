# go-choice

chooser.Run([]string) -> string

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/tayusa/go-choice"
)

func main() {
	chooser, err := choice.NewChooser([]string{"test9", "te11", "jfdls", "fdsaf", "daj", "fdsie", "feafii", "fdiaoeioa", "feiaofjl"})
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Println(chooser.Run())
}
```
