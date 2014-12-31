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

type GetController struct {
	controller
}

func (c *GetController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:  "get",
			Usage: "Retrieve API resources",
			Subcommands: []cli.Command{
				{
					Name:  "info",
					Usage: "get API information",
					Action: func(context *cli.Context) {
						c.GetResource("/info")
					},
				},
				{
					Name:  "users",
					Usage: "get users",
					Action: func(context *cli.Context) {
						c.GetResource("/users/%s", context.Args().First())
					},
				},
				{
					Name:  "tenants",
					Usage: "get tenants",
					Action: func(context *cli.Context) {
						c.GetResource("/tenants/%s", context.Args().First())
					},
				},
				{
					Name:  "cmdbs",
					Usage: "get CMDBs",
					Action: func(context *cli.Context) {
						c.GetResource("/cmdbs/%s", context.Args().First())
					},
				},
				{
					Name:  "citypes",
					Usage: "get CI Types",
					Action: func(context *cli.Context) {
						c.GetResource("/cmdbs/%s/citypes/%s", context.GlobalString("cmdb"), context.Args().First())
					},
				},
			},
		},
	}...)

	return nil
}
