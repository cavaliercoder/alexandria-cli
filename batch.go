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
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"strings"
)

type BatchCommand struct {
	Name         string
	Cmdb         string
	ResourceType string `json:"type"`
	CIType       string
	Action       string
	Target       string
	Content      json.RawMessage
	SubCommands  []BatchCommand
}

type BatchScript []BatchCommand

type BatchController struct {
	controller
}

func (c *BatchController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:   "batch",
			Usage:  "Batch API operations",
			Action: c.Execute,
		},
	}...)

	return nil
}

func (c *BatchController) Execute(context *cli.Context) {
	// Decode the resource from STDIN or from the first command argument?
	var err error
	var filename string = context.Args().First()
	var input io.ReadCloser
	if filename == "" {
		NotifyStdin()
		input = os.Stdin
	} else {
		input, err = os.Open(filename)
		if err != nil {
			Die(err)
		}
	}

	// Decode batch script
	var script BatchScript
	err = json.NewDecoder(input).Decode(&script)
	if err != nil {
		Die(err)
	}
	input.Close()

	cmdb := context.GlobalString("cmdb")
	c.ExecuteCommands(script, cmdb, "", 0)
}

func (c *BatchController) ExecuteCommands(commands []BatchCommand, cmdb string, path string, depth int) error {
	count := len(commands)
	indent := strings.Repeat("    ", depth)
	for i, command := range commands {
		fmt.Printf("%s==> Executing task %d/%d [%s%s]\n", indent, i+1, count, path, command.Name)

		// Switch CMDB if specified
		if command.Cmdb != "" {
			cmdb = command.Cmdb
		}

		// Determine the URL
		var url string
		switch command.ResourceType {
		case "user":
			url = "/users"

		case "tenant":
			url = "/tenants"

		case "cmdb":
			url = "/cmdbs"

		case "citype":
			url = fmt.Sprintf("/cmdbs/%s/citypes", cmdb)

		case "ci":
			if command.CIType == "" {
				Dief("No CI Type specified in '%v' command %v", command.Action, i+1)
			}

			url = fmt.Sprintf("/cmdbs/%s/%s", cmdb, command.CIType)
		default:
			Dief("Unsupported resource type '%v' in '%v' command %v", command.ResourceType, command.Action, i+1)
		}

		// Append target for commands that need one
		if command.Action != "add" {
			if command.Target == "" {
				Dief("No target resource specified in '%s' command %d [%s]", command.Action, i+1, command.Name)
			}

			url = fmt.Sprintf("%s/%s", url, command.Target)
		}

		// Perform the action
		fmt.Printf("%s    ", indent)
		switch command.Action {
		case "add":
			c.AddResource(url, string(command.Content))

		case "update":
			c.UpdateResource(url, string(command.Content))

		case "delete":
			c.DeleteResource(url)

		default:
			Dief("Unsupport action in command %v: %s", i+1, command.Action)
		}

		// Execute subcommands
		if 0 < len(command.SubCommands) {
			fmt.Printf("%s    Subcommands:\n", indent)
			c.ExecuteCommands(command.SubCommands, cmdb, fmt.Sprintf("%s%s/", path, command.Name), depth+1)
		}
	}

	return nil
}
