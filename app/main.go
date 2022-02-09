package main

import (
	api "boilerplate-golang-v2/api"
	"boilerplate-golang-v2/app/modules"
	"boilerplate-golang-v2/config"
	"boilerplate-golang-v2/util"
	"context"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

var banner = `

██████╗  ██████╗ ███████╗███████╗██╗██████╗  ██████╗ ███╗   ██╗
██╔══██╗██╔═══██╗██╔════╝██╔════╝██║██╔══██╗██╔═══██╗████╗  ██║
██████╔╝██║   ██║███████╗█████╗  ██║██║  ██║██║   ██║██╔██╗ ██║
██╔═══╝ ██║   ██║╚════██║██╔══╝  ██║██║  ██║██║   ██║██║╚██╗██║
██║     ╚██████╔╝███████║███████╗██║██████╔╝╚██████╔╝██║ ╚████║
╚═╝      ╚═════╝ ╚══════╝╚══════╝╚═╝╚═════╝  ╚═════╝ ╚═╝  ╚═══╝
v1.0.0-alpha                                                        
`

func main() {
	stdLog.Println(banner)

	//load config if available or set to default
	config := config.GetConfig()

	//initialize database connection based on given config
	dbCon := util.NewDatabaseConnection(config)

	//initiate item repository
	controllers := modules.RegisterController(dbCon)

	//create echo http
	e := echo.New()
	e.HideBanner = true
	// index route
	e.GET("/", func(c echo.Context) error {
		message := `Aku Ingin

Aku ingin mencintaimu dengan sederhana
dengan kata yang tak sempat diucapkan
kayu kepada api yang menjadikannya abu

Aku ingin mencintaimu dengan sederhana
dengan isyarat yang tak sempat disampaikan
awan kepada hujan yang menjadikannya tiada

-- Sapardi Djoko Damono`
		return c.String(http.StatusOK, message)
	})

	//register API path and handler
	api.RegisterPath(e, controllers)

	// run server
	go func() {
		address := fmt.Sprintf("localhost:%d", config.App.Port)
		if err := e.Start(address); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//close db
	defer dbCon.CloseConnection()

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
