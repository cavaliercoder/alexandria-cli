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
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"net/http"
	"os"
	"strings"
)

type Controller interface {
	Init(app *cli.App) error
}

type controller struct {
	app     *cli.App
	baseUrl string
}

func (c *controller) ApiRequest(method string, path string, body io.Reader) (*http.Response, error) {
	context := GetContext()
	url := context.GlobalString("url")
	apiKey := context.GlobalString("api-key")

	// Validate URL and API Key
	if url == "" {
		Die("API base URL not specified")
	}
	if apiKey == "" {
		Die("API authentication key not specified")
	}

	// Formulate request URL
	url = fmt.Sprintf("%s%s?pretty=true", url, path)
	Dprintf("API Request: %s %s", method, url)

	// Create a HTTP client that does not follow redirects
	// This allows 'Location' headers to be printed to the CLI
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Never follow redirects")
		},
	}
	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}

	// Add request headers
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Auth-Token", apiKey)
	req.Header.Add("User-Agent", "Alexandria CMDB CLI")

	// Submit the request
	res, err := client.Do(req)

	return res, err
}

func (c *controller) ApiResult(res *http.Response) {
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

func (c *controller) ApiError(res *http.Response) {
	Die(res.Status)
}

func (c *controller) AddResource(path string, resource string) {
	// Decode the resource from STDIN or from the first command argument?
	var input io.Reader
	if resource == "" {
		NotifyStdin()
		input = os.Stdin
	} else {
		input = strings.NewReader(resource)
	}

	res, err := c.ApiRequest("POST", path, input)
	defer res.Body.Close()
	if err != nil {
		Die(err)
	}

	if res.StatusCode == http.StatusCreated {
		fmt.Printf("Created %s\n", res.Header.Get("Location"))
	} else {
		c.ApiError(res)
	}
}

func (c *controller) AddResourceAction(context *cli.Context) {
	c.AddResource(c.baseUrl, context.Args().First())
}

func (c *controller) GetResource(format string, a ...interface{}) {
	// Get requested resource ID from first command argument
	var err error
	var res *http.Response

	path := strings.TrimRight(fmt.Sprintf(format, a...), "/")
	res, err = c.ApiRequest("GET", path, nil)
	if err != nil {
		Die(err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		c.ApiResult(res)
	case http.StatusNotFound:
		Die(fmt.Sprintf("No such resource found at %s", path))
	default:
		c.ApiError(res)
	}
}

func (c *controller) GetResourceAction(context *cli.Context) {
	c.GetResource("%s/%s", c.baseUrl, context.Args().First())
}

func (c *controller) DeleteResource(format string, a ...interface{}) {
	var err error
	var res *http.Response

	path := fmt.Sprintf(format, a...)
	res, err = c.ApiRequest("DELETE", path, nil)
	if err != nil {
		Die(err)
	}

	switch res.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("Deleted %s\n", path)
	case http.StatusNotFound:
		Die(fmt.Sprintf("No such resource found at %s", path))
	default:
		c.ApiError(res)
	}
}

func (c *controller) DeleteResourceAction(context *cli.Context) {
	id := context.Args().First()
	if id == "" {
		Die("No resource specified for deletion")
	}
	c.DeleteResource("%s/%s", c.baseUrl, id)
}
