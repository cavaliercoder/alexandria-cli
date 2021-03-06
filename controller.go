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
	"net/url"
	"os"
	"strings"
)

type Controller interface {
	Init(app *cli.App) error
}

type controller struct {
	app *cli.App
}

func (c *controller) ApiRequest(method string, path string, body io.Reader) (*http.Response, error) {
	context := GetContext()
	url := context.GlobalString("url")
	apiKey := context.GlobalString("api-key")
	format := context.GlobalString("format")

	// Validate URL and API Key
	if url == "" {
		Die("API base URL not specified")
	}
	if apiKey == "" {
		Die("API authentication key not specified")
	}

	// Formulate request URL
	if strings.Contains(path, "?") {
		url = fmt.Sprintf("%s%s&format=%s&pretty=true", url, path, format)
	} else {
		url = fmt.Sprintf("%s%s?format=%s&pretty=true", url, path, format)
	}
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
	Dief("Unexpected response: %s", res.Status)
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
	if err != nil {
		Die(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusCreated:
		fmt.Printf("Created %s\n", res.Header.Get("Location"))
	case http.StatusConflict:
		Die("Resource conflicts with an existing resource")
	default:
		c.ApiError(res)
	}
}

func (c *controller) UpdateResource(path string, resource string) {

	// Decode the resource from STDIN or from the first command argument?
	var input io.Reader
	if resource == "" {
		NotifyStdin()
		input = os.Stdin
	} else {
		input = strings.NewReader(resource)
	}

	res, err := c.ApiRequest("PUT", path, input)
	if err != nil {
		Die(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusNoContent:
		fmt.Printf("Updated %s\n", path)
	case http.StatusMovedPermanently:
		fmt.Printf("Updated %s (relocated to %s)\n", path, res.Header.Get("Location"))
	case http.StatusNotFound:
		Die(fmt.Sprintf("No such resource found at %s", path))
	default:
		c.ApiError(res)
	}
}

func (c *controller) GetResource(sel string, format string, a ...interface{}) {
	// Get requested resource ID from first command argument
	var err error
	var res *http.Response

	path := strings.TrimRight(fmt.Sprintf(format, a...), "/")

	// Include field selector
	if sel != "" {
		path = fmt.Sprintf("%s?select=%s", path, url.QueryEscape(sel))
	}

	// Fetch the request
	res, err = c.ApiRequest("GET", path, nil)
	if err != nil {
		Die(err)
	}

	// Parse the result
	switch res.StatusCode {
	case http.StatusOK:
		c.ApiResult(res)
	case http.StatusNotFound:
		Die(fmt.Sprintf("No such resource found at %s", path))
	default:
		c.ApiError(res)
	}
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
