package intergration

import (
	"fmt"
	"os"
)

func fatalF(context string, args ...interface{}) {
	fmt.Printf(context, args...)
	os.Exit(1)
}
