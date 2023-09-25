package main

import "gin_learn/webook/internal/web"

func main() {
	server := web.RegisterRoutes()
	err := server.Run(":8000")
	if err != nil {
		return
	}
}
