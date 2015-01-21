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
	"fmt"
	"github.com/codegangsta/cli"
)

type AddController struct {
	controller
}

func (c *AddController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:  "add",
			Usage: "Create API resources",
			Subcommands: []cli.Command{
				{
					Name:  "user",
					Usage: "add a user",
					Action: func(context *cli.Context) {
						c.AddResource("/users", context.Args().First())
					},
				},
				{
					Name:  "tenant",
					Usage: "get a tenant",
					Action: func(context *cli.Context) {
						c.AddResource("/tenants", context.Args().First())
					},
				},
				{
					Name:  "cmdb",
					Usage: "add a CMDB",
					Action: func(context *cli.Context) {
						c.AddResource("/cmdbs", context.Args().First())
					},
				},
				{
					Name:  "citype",
					Usage: "add a CI Types",
					Action: func(context *cli.Context) {
						c.AddResource(fmt.Sprintf("/cmdbs/%s/citypes", context.GlobalString("cmdb")), context.Args().First())
					},
				},
				{
					Name:  "ci",
					Usage: "add a CI",
					Action: func(context *cli.Context) {
						citype := context.Args().First()
						body := context.Args().Get(1)

						fmt.Printf("CI Type: %v\nBody: %#v\n", citype, body)
						c.AddResource(fmt.Sprintf("/cmdbs/%s/%s", context.GlobalString("cmdb"), citype), body)
					},
				},
			},
		},
	}...)

	return nil
}
