package main

import "os"

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(port)

}
