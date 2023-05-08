package main

import (
	"budgetting/api"
	"budgetting/config"
	"fmt"
)

func main() {
	cfg := config.Load()
	err := api.ListenAndServe(cfg)
	if err != nil {
		fmt.Println(err)
	}
}
