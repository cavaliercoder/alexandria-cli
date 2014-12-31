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
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
	"regexp"
	"strings"
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
			Value:  "http://localhost:3000/api/v1",
			Usage:  "specify the API base URL",
			EnvVar: "ALEX_API_URL",
		},
		cli.StringFlag{
			Name:   "k, api-key",
			Usage:  "specify the API authentication key",
			EnvVar: "ALEX_API_KEY",
		},
		cli.StringFlag{
			Name:   "c, cmdb",
			Value:  "default",
			Usage:  "specify the CMDB to use for CI queries",
			EnvVar: "ALEX_CMDB",
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

	// Make context available globally
	app.Before = SetContext

	// Load config from ~/.alexrc
	LoadRc()

	// Add controllers
	var err error
	controllers := []Controller{
		&GetController{},
		&AddController{},
		&DeleteController{},
	}

	for _, controller := range controllers {
		err = controller.Init(app)
		if err != nil {
			Die(err)
		}
	}

	app.Run(os.Args)
}

// LoadRc loads configuration settings from ~/.alexrc using a simple
// KEY="value" format. The KEY field should reflect an available
// environment variable setting as described with --help
func LoadRc() {
	usr, _ := user.Current()
	path := fmt.Sprintf("%s/.alexrc", usr.HomeDir)

	// Open ~/.alexrc
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Read each line
	lineNo := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()

		// Skip comments and blank lines
		if match, _ := regexp.MatchString("^(#.*|\\s*)$", line); match {
			continue
		}

		// Parse key/vals
		r := regexp.MustCompile(`^(\w+)=(.+)$`)
		matches := r.FindStringSubmatch(line)
		if matches == nil {
			Warnf("%s:%d Invalid declaration: %s", path, lineNo, line)
			continue
		}

		key := matches[1]
		val := strings.Trim(matches[2], "\" ")

		// Set environment variables locally
		os.Setenv(key, val)
	}
}
