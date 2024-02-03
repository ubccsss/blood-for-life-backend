package main

import (
	"blood-for-life-backend/store"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	db, err := sqlx.Connect() // configuration needed

	if err != nil {
		panic(err)
	}

	eventStore := store.NewPGEventStore(db)
	bind(e, eventStore)

	e.Logger.Fatal(e.Start(":1323"))

}

func bind(e *echo.Echo, eventStore store.EventStore) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/api/create", func(c echo.Context) error {
		event := new(store.Event)

		// parse request body
		bindErr := c.Bind(event)

		if bindErr != nil {
			return c.JSON(http.StatusBadRequest, bindErr)
		}

		_, err := eventStore.Create(c.Request().Context(), *event)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, event)
	})
}
