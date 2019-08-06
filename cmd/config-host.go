/*
 * Minio Client (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"encoding/json"

	"github.com/minio/cli"
	"github.com/didiyun/mc/pkg/console"
	"github.com/didiyun/mc/pkg/probe"
)

var configHostCmd = cli.Command{
	Name:   "host",
	Usage:  "list, modify and remove hosts in configuration file",
	Action: mainConfigHost,
	Before: setGlobalsFromContext,
	Flags:  globalFlags,
	Subcommands: []cli.Command{
		configHostAddCmd,
		configHostRemoveCmd,
		configHostListCmd,
	},
	HideHelpCommand: true,
}

// mainConfigHost is the handle for "mc config host" command.
func mainConfigHost(ctx *cli.Context) error {
	cli.ShowCommandHelp(ctx, ctx.Args().First())
	return nil
	// Sub-commands like "remove", "list" have their own main.
}

// hostMessage container for content message structure
type hostMessage struct {
	op          string
	prettyPrint bool
	Status      string `json:"status"`
	Alias       string `json:"alias"`
	URL         string `json:"URL"`
	AccessKey   string `json:"accessKey,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
	API         string `json:"api,omitempty"`
	Lookup      string `json:"lookup,omitempty"`
}

// Print the config information of one alias, when prettyPrint flag
// is activated, fields contents are cut and '...' will be added to
// show a pretty table of all aliases configurations
func (h hostMessage) String() string {
	switch h.op {
	case "list":
		urlFieldMaxLen, apiFieldMaxLen := -1, -1
		accessFieldMaxLen, secretFieldMaxLen := -1, -1
		lookupFieldMaxLen := -1
		// Set cols width if prettyPrint flag is enabled
		if h.prettyPrint {
			urlFieldMaxLen = 30
			accessFieldMaxLen = 20
			secretFieldMaxLen = 40
			apiFieldMaxLen = 5
			lookupFieldMaxLen = 5
		}

		// Create a new pretty table with cols configuration
		t := newPrettyTable("  ",
			Field{"Alias", -1},
			Field{"URL", urlFieldMaxLen},
			Field{"AccessKey", accessFieldMaxLen},
			Field{"SecretKey", secretFieldMaxLen},
			Field{"API", apiFieldMaxLen},
			Field{"Lookup", lookupFieldMaxLen},
		)
		return t.buildRow(h.Alias, h.URL, h.AccessKey, h.SecretKey, h.API, h.Lookup)
	case "remove":
		return console.Colorize("HostMessage", "Removed `"+h.Alias+"` successfully.")
	case "add":
		return console.Colorize("HostMessage", "Added `"+h.Alias+"` successfully.")
	default:
		return ""
	}
}

// JSON jsonified host message
func (h hostMessage) JSON() string {
	h.Status = "success"
	jsonMessageBytes, e := json.Marshal(h)
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(jsonMessageBytes)
}
