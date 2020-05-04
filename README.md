# go-choice
## examples
### []string
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
	fmt.Println(choice.ToString(chooser.Run()))
}
```

### []fmt.Stringer
```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tayusa/go-choice"
)

type file struct {
	path string
	info os.FileInfo
}

func (f file) String() string {
	return f.path
}

func getFiles() ([]fmt.Stringer, error) {
	wd, err := os.Getwd()
	if err != nil {
		return make([]fmt.Stringer, 0), err
	}

	var files []fmt.Stringer
	if err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, file{info: info, path: path})

		return nil
	}); err != nil {
		return make([]fmt.Stringer, 0), err
	}
	return files, nil
}

func main() {
	files, _ := getFiles()
	chooser, err := choice.NewChooser(files)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Println(chooser.Run())
}
```
