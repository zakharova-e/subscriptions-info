package main

import (
	_ "github.com/zakharova-e/subscriptions-info/internal/connections"
	"github.com/zakharova-e/subscriptions-info/internal/web"
)

func main() {
	web.Run()
}
