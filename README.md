# go-chooser
## examples
```go
package main

import (
	"fmt"
	"log"

	"github.com/tayusa/go-chooser"
)

func main() {
	chooser, err := chooser.NewChooser([]string{"test9", "te11", "jfdls", "fdsaf", "daj", "fdsie", "feafii", "fdiaoeioa", "feiaofjl"})
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Println(chooser.Run())
}
```
