package print

import (
	"encoding/json"
	"fmt"
)

func Pretty(i interface{}) {
	j, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(j))
}
