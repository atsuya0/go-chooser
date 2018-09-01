# selector

selector.Run([]string) -> string

```go
package main

import (
	"fmt"
	"github.com/tayusa/selector"
)

func main() {
	s := selector.NewSelector([]string{"test9", "te11", "jfdls", "fdsaf", "daj", "fdsie", "feafii", "fdiaoeioa", "feiaofjl"})
	fmt.Print(s.Run())
}
```
