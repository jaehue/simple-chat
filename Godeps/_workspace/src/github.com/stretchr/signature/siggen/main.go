package main

import (
	"fmt"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/objx"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/signature"
	"github.com/stretchr/commander"
)

// Generates a 32 character random signature
func main() {
	commander.Go(func() {

		commander.Map(commander.DefaultCommand, "", "",
			func(args objx.Map) {
				fmt.Println(signature.RandomKey(32))
			})

		commander.Map("len length=(int)", "Key length",
			"Specify the length of the generated key",
			func(args objx.Map) {
				length := args.Get("length").Int()
				fmt.Println(signature.RandomKey(length))
			})

	})
}
