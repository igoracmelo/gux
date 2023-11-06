package main

import (
	"fmt"
	"net/http"

	"github.com/igoracmelo/gux"
)

func main() {
	mux := gux.New()

	mux.Get("/users", handleFindUsers)
	mux.Post("/users", handleCreateUser)
	mux.Get("/users/:id", handleFindUserByID)
	mux.Put("/users/:id", handleUpdateUserByID)
	mux.Delete("/users/:id", handleDeleteUserByID)

	http.ListenAndServe(":3000", mux)
}

func handleFindUsers(c *gux.Ctx) {
	fmt.Fprintln(c.W, `[{"id": 123, "name": "someone"}]`)
}

func handleCreateUser(c *gux.Ctx) {
	c.W.WriteHeader(http.StatusCreated)
}

func handleFindUserByID(c *gux.Ctx) {
	fmt.Fprintln(c.W, `{"id": `+c.Vars["id"]+`, "name": "someone"}`)
}

func handleUpdateUserByID(c *gux.Ctx) {
	gux.NotFound(c)
}

func handleDeleteUserByID(c *gux.Ctx) {
	c.W.WriteHeader(http.StatusOK)
}
