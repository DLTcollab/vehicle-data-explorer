package main

import controller "biilabs.com/mam-data-explorer/controllers"

func main() {
	r := controller.SetupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
