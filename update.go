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

type UpdateController struct {
	controller
}

func (c *UpdateController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:  "update",
			Usage: "Update API resources",
			Subcommands: []cli.Command{
				{
					Name:  "user",
					Usage: "update a user",
					Action: func(context *cli.Context) {
						c.UpdateResource(fmt.Sprintf("/users/%s", context.Args().First()), context.Args().Get(2))
					},
				},
				{
					Name:  "tenant",
					Usage: "update a tenant",
					Action: func(context *cli.Context) {
						c.UpdateResource(fmt.Sprintf("/tenants/%s", context.Args().First()), context.Args().Get(2))
					},
				},
				{
					Name:  "cmdb",
					Usage: "update a CMDB",
					Action: func(context *cli.Context) {
						c.UpdateResource(fmt.Sprintf("/cmdbs/%s", context.Args().First()), context.Args().Get(2))
					},
				},
				{
					Name:  "citype",
					Usage: "update a CI Type",
					Action: func(context *cli.Context) {
						c.UpdateResource(fmt.Sprintf("/cmdbs/%s/citypes/%s", context.GlobalString("cmdb"), context.Args().First()), context.Args().Get(2))
					},
				},
			},
		},
	}...)

	return nil
}
