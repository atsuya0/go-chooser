# go-chooser
![go-chooser](https://user-images.githubusercontent.com/37957375/81494359-3d2fdd80-92e3-11ea-87a3-ed5df76b8da8.gif)
# examples
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
