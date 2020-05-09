# go-chooser
![go-chooser](https://user-images.githubusercontent.com/37957375/81478817-a0296200-925a-11ea-8fd9-29cda0eaa79f.gif)
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
