/*
 * Minio Client (C) 2018 Minio, Inc.
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
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/didiyun/mc/pkg/console"
	"github.com/didiyun/mc/pkg/probe"
)

var adminPolicyAddCmd = cli.Command{
	Name:   "add",
	Usage:  "add new policy",
	Action: mainAdminPolicyAdd,
	Before: setGlobalsFromContext,
	Flags:  globalFlags,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} TARGET POLICYNAME POLICYFILE

POLICYNAME:
  Name of the canned policy on Minio server.

POLICYFILE:
  Name of the policy file associated with the policy name.

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
  1. Add a new canned policy 'writeonly'.
     $ {{.HelpName}} myminio writeonly /tmp/writeonly.json
 `,
}

// checkAdminPolicyAddSyntax - validate all the passed arguments
func checkAdminPolicyAddSyntax(ctx *cli.Context) {
	if len(ctx.Args()) != 3 {
		cli.ShowCommandHelpAndExit(ctx, "add", 1) // last argument is exit code
	}
}

// userPolicyMessage container for content message structure
type userPolicyMessage struct {
	op         string
	Status     string `json:"status"`
	Policy     string `json:"policy,omitempty"`
	PolicyJSON []byte `json:"policyJSON,omitempty"`
}

func (u userPolicyMessage) String() string {
	switch u.op {
	case "list":
		if len(u.PolicyJSON) > 0 {
			return string(u.PolicyJSON)
		}
		policyFieldMaxLen := 20
		// Create a new pretty table with cols configuration
		return newPrettyTable("  ",
			Field{"Policy", policyFieldMaxLen},
		).buildRow(u.Policy)
	case "remove":
		return console.Colorize("PolicyMessage", "Removed policy `"+u.Policy+"` successfully.")
	case "add":
		return console.Colorize("PolicyMessage", "Added policy `"+u.Policy+"` successfully.")
	}
	return ""
}

func (u userPolicyMessage) JSON() string {
	u.Status = "success"
	jsonMessageBytes, e := json.Marshal(u)
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(jsonMessageBytes)
}

// mainAdminPolicyAdd is the handle for "mc admin policy add" command.
func mainAdminPolicyAdd(ctx *cli.Context) error {
	checkAdminPolicyAddSyntax(ctx)

	console.SetColor("PolicyMessage", color.New(color.FgGreen))

	// Get the alias parameter from cli
	args := ctx.Args()
	aliasedURL := args.Get(0)

	policy, e := ioutil.ReadFile(args.Get(2))
	fatalIf(probe.NewError(e).Trace(args...), "Unable to get policy")

	// Create a new Minio Admin Client
	client, err := newAdminClient(aliasedURL)
	fatalIf(err, "Cannot get a configured admin connection.")

	fatalIf(probe.NewError(client.AddCannedPolicy(args.Get(1), string(policy))).Trace(args...), "Cannot add new policy")

	printMsg(userPolicyMessage{
		op:     "add",
		Policy: args.Get(1),
	})

	return nil
}
