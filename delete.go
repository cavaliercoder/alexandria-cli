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
	"github.com/codegangsta/cli"
)

type DeleteController struct {
	controller
}

func (c *DeleteController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:  "delete",
			Usage: "Delete API resources",
			Subcommands: []cli.Command{
				{
					Name:  "user",
					Usage: "delete a user",
					Action: func(context *cli.Context) {
						c.DeleteResourceAction(context, "/users")
					},
				},
				{
					Name:  "tenant",
					Usage: "get a tenant",
					Action: func(context *cli.Context) {
						c.DeleteResourceAction(context, "/tenants")
					},
				},
				{
					Name:  "cmdb",
					Usage: "add a CMDB",
					Action: func(context *cli.Context) {
						c.DeleteResourceAction(context, "/cmdbs")
					},
				},
			},
		},
	}...)

	return nil
}

func (c *controller) DeleteResourceAction(context *cli.Context, baseUrl string) {
	id := context.Args().First()
	if id == "" {
		Die("No resource specified for deletion")
	}
	c.DeleteResource("%s/%s", baseUrl, id)
}
