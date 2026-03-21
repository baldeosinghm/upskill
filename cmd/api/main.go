package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("upskill-api starting...")
	http.ListenAndServe(":8080", nil)
}
