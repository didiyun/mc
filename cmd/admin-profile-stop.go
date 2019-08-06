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
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/minio/cli"
	"github.com/didiyun/mc/pkg/probe"
	"github.com/didiyun/mc/pkg/console"
)

var adminProfileStopCmd = cli.Command{
	Name:            "stop",
	Usage:           "stop and download profile data",
	Action:          mainAdminProfileStop,
	Before:          setGlobalsFromContext,
	Flags:           globalFlags,
	HideHelpCommand: true,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
    2. Download latest profile data in the current directory
       $ {{.HelpName}} myminio/
`,
}

func checkAdminProfileStopSyntax(ctx *cli.Context) {
	if len(ctx.Args()) != 1 {
		cli.ShowCommandHelpAndExit(ctx, "stop", 1) // last argument is exit code
	}
}

// mainAdminProfileStop - the entry function of profile stop command
func mainAdminProfileStop(ctx *cli.Context) error {
	// Check for command syntax
	checkAdminProfileStopSyntax(ctx)

	// Get the alias parameter from cli
	args := ctx.Args()
	aliasedURL := args.Get(0)

	// Create a new Minio Admin Client
	client, err := newAdminClient(aliasedURL)
	if err != nil {
		fatalIf(err.Trace(aliasedURL), "Cannot initialize admin client.")
		return nil
	}

	// Create profile zip file
	tmpFile, e := ioutil.TempFile("", "mc-profile-")
	fatalIf(probe.NewError(e), "Unable to download profile data.")

	// Ask for profile data, which will come compressed with zip format
	zippedData, adminErr := client.DownloadProfilingData()
	fatalIf(probe.NewError(adminErr), "Unable to download profile data.")

	// Copy zip content to target download file
	_, e = io.Copy(tmpFile, zippedData)
	fatalIf(probe.NewError(e), "Unable to download profile data.")

	// Close everything
	zippedData.Close()
	tmpFile.Close()

	downloadPath := "profile.zip"

	fi, e := os.Stat(downloadPath)
	if e == nil && !fi.IsDir() {
		e = os.Rename(downloadPath, downloadPath+"."+time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
		fatalIf(probe.NewError(e), "Unable to create a backup of profile.zip")
	} else {
		if !os.IsNotExist(e) {
			fatal(probe.NewError(e), "Unable to download profile data.")
		}
	}

	fatalIf(probe.NewError(os.Rename(tmpFile.Name(), downloadPath)), "Unable to download profile data.")

	console.Infof("Profile data successfully downloaded as %s\n", downloadPath)
	return nil
}
