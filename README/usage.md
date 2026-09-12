```go
package main

import (
    "fmt"

    "github.com/MarkRosemaker/openapi"
)

func main() {
    doc, err := openapi.LoadFromFile("path/to/openapi.json") // or openapi.yaml
    if err != nil {
        fmt.Println("Error parsing spec:", err)
        return
    }

    if err := doc.Validate(); err != nil {
        fmt.Println("Error validating spec:", err)
        return
    }

    // sort keys of each component in alphabetical order
    doc.Components.SortMaps()

	// write an improved version of your spec
    if err := doc.WriteToFile("path/to/openapi.json"); err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }
}
```
