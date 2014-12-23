/*
 * Alexandria CMDB - Open source configuration management database
 * Copyright (C) 2014  Ryan Armstrong <ryan@cavaliercoder.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * package controllers
 */
package main

import (
	. "alexandria/cli/application"
	"alexandria/cli/controllers"
	"github.com/codegangsta/cli"
	"os"
)

var app *cli.App

func main() {
	app = cli.NewApp()
	app.Name = "alex"
	app.Usage = "Alexandria CMDB CLI"
	app.Version = "1.0.0"
	app.Author = "Ryan Armstrong"
	app.Email = "ryan@cavaliercoder.com"

	// Global args
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "u, url",
			Value:  "http://localhost:3000",
			Usage:  "specify the API base URL",
			EnvVar: "ALEX_API_URL",
		},
		cli.StringFlag{
			Name:   "k, api-key",
			Usage:  "specify the API authentication key",
			EnvVar: "ALEX_API_KEY",
		},
		cli.BoolFlag{
			Name:  "i, stdin",
			Usage: "read request body from stdin",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Show extra verbose output",
			EnvVar: "ALEX_DEBUG",
		},
	}

	app.Before = SetContext

	// Add controllers
	var err error
	controllers := []controllers.Controller{
		&controllers.CITypeController{},
		&controllers.ConfigController{},
		&controllers.ResourceController{},
		&controllers.TenantController{},
		&controllers.UserController{},
	}

	for _, controller := range controllers {
		err = controller.Init(app)
		if err != nil {
			Die(err)
		}
	}

	app.Run(os.Args)
}
