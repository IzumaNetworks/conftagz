//go:build debugconftagz
// +build debugconftagz

package conftagz

import "fmt"

func debugf(format string, args ...interface{}) {
	fmt.Printf("conftags: "+format, args...)
}

// func debugf(format string, args ...interface{}) {

// }
