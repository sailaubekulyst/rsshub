package main

import (
	"fmt"
	"os"
	"rsshub/internal/application"
)

func main() {
	app, err := application.GetNewApp(os.Args[1:])
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = app.Run()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	
}
