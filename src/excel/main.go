package main

import (
	"fmt"
	"log"
	registry "golang.org/x/sys/windows/registry"

)



func main() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SAM\XIAO`, registry.ALL_ACCESS)

	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("x")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", s)
}
