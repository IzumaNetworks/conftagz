//go:build debug
// +build debug

package conftagz

import "fmt"

func debugf(format string, args ...interface{}) {
	fmt.Printf("conftags: "+format, args...)
}

// func debugf(format string, args ...interface{}) {

// }
