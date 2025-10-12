package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

func main() {
	name := flag.String("name", "world", "Name to greet")
	flag.Parse()
	env := os.Getenv("APP_ENV")
	cfg := Config{
		Name: *name,
		Env:  env,
	}
	out, _ := json.MarshalIndent(cfg, "", " ")
	fmt.Println(string(out))
}
