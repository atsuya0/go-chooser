# go-chooser
![go-chooser](https://user-images.githubusercontent.com/37957375/81494359-3d2fdd80-92e3-11ea-87a3-ed5df76b8da8.gif)

# Key bindings
| key | description |
| --- | ----------- |
| Enter | Returns the character string under the cursor. Or returns the selected character string. |
| Tab | Select the character string under the cursor. |
| Backspace | Delete a character before cursor. |
| Delete | Delete a character under cursor. |
| C-d | Delete a character under cursor. |
| C-h | Delete a character before cursor. |
| C-f | Move forward a character. |
| C-b | Move backward a character. |
| C-a | Go to the beginning of the line. |
| C-e | Go to the end of the line. |
| C-u | Kill characters from cursor current position to the beginning of the line. |
| C-k | Kill characters from cursor current position to the end of the line. |
| C-w | Delete before a word. |
| C-n | Move the cursor to the next line. |
| C-p | Move the cursor to the previous line. |
| ? | Display the key bindings. |

# Examples
```go
package main

import (
	"fmt"
	"log"

	"github.com/tayusa/go-chooser"
)

func main() {
	chooser, err := chooser.NewChooser(
		[]string{
			"about five hundred yen",
			"get to the airport",
			"be angry with sb",
			"play baseball",
			"bring money",
			"listen to Beatle’s CD",
			"fifty cents",
			"be dead",
			"be not far",
			"use a fork",
			"pull a door",
			"listen to music",
			"put salt"})
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	/*
		index, string, err := chooser.SingleRun()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		fmt.Printf("%#v\n", index)
		fmt.Printf("%#v\n", string)
	*/

	indexes, strings, err := chooser.Run()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("%#v\n", indexes)
	fmt.Printf("%#v\n", strings)
}
```
