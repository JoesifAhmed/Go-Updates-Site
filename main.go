package main

import (
	"net/http"

	"go.moc/models"
	"go.moc/routes"
	"go.moc/utils"
)

func main() {
	models.InitClient()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
