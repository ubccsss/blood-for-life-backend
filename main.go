package main

import (
	"blood-for-life-backend/apimodels"
	"fmt"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Println(envErr.Error())
	}

	dbConnection := os.Getenv("DB")

	db, err := sqlx.Connect("postgres", dbConnection)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully Connected")
	}

	loadSchema(db)

	eventStore := apimodels.NewPGEventStore(db)
	bind(e, eventStore)

	e.Logger.Fatal(e.Start(":1323"))

}

func loadSchema(db *sqlx.DB) {
	file, err := os.ReadFile("./db/schema.sql")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = db.Exec(string(file))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func bind(e *echo.Echo, eventStore apimodels.EventStore) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/api/create", func(c echo.Context) error {
		event := new(apimodels.CreateEventModel)

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
