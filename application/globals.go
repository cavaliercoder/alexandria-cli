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
package application

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
)

var context *cli.Context

func SetContext(c *cli.Context) error {
	context = c
	return nil
}

func Die(message interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}

func DPrint(message interface{}) {
	if context.GlobalBool("debug") {
		fmt.Fprintf(os.Stderr, "DEBUG: %s\n", message)
	}
}

func DPipeToStderr(reader io.Reader) {
	if context.GlobalBool("debug") {
		buf := bufio.NewReader(reader)

		var line []byte
		var err error = nil
		for err == nil {
			line, _, err = buf.ReadLine()

			DPrint(line)
		}
	}
}

func PipeToFile(reader io.Reader, file *os.File) {
	io.Copy(file, reader)
}
