/*
 * Minio Client (C) 2015 Minio, Inc.
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
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/didiyun/mc/pkg/probe"
	"github.com/didiyun/mc/pkg/console"
)

// diff specific flags.
var (
	diffFlags = []cli.Flag{}
)

// Compute differences in object name, size, and date between two buckets.
var diffCmd = cli.Command{
	Name:   "diff",
	Usage:  "list differences in object name, size, and date between two buckets",
	Action: mainDiff,
	Before: setGlobalsFromContext,
	Flags:  append(diffFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] FIRST SECOND

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
DESCRIPTION:
  Diff only calculates differences in object name, size and time.
  It *DOES NOT* compare objects' contents.

LEGEND:
    > - object is only in source.
    < - object is only in destination.
    ! - newer object is in source.

EXAMPLES:
  1. Compare a local folder with a folder on Amazon S3 cloud storage.
     $ {{.HelpName}} ~/Photos s3/mybucket/Photos

  2. Compare two folders on a local filesystem.
     $ {{.HelpName}} ~/Photos /Media/Backup/Photos
`,
}

// diffMessage json container for diff messages
type diffMessage struct {
	Status        string       `json:"status"`
	FirstURL      string       `json:"first"`
	SecondURL     string       `json:"second"`
	Diff          differType   `json:"diff"`
	Error         *probe.Error `json:"error,omitempty"`
	firstContent  *clientContent
	secondContent *clientContent
}

// String colorized diff message
func (d diffMessage) String() string {
	msg := ""
	switch d.Diff {
	case differInFirst:
		msg = console.Colorize("DiffOnlyInFirst", "< "+d.FirstURL)
	case differInSecond:
		msg = console.Colorize("DiffOnlyInSecond", "> "+d.SecondURL)
	case differInType:
		msg = console.Colorize("DiffType", "! "+d.SecondURL)
	case differInSize:
		msg = console.Colorize("DiffSize", "! "+d.SecondURL)
	case differInTime:
		msg = console.Colorize("DiffTime", "! "+d.SecondURL)
	default:
		fatalIf(errDummy().Trace(d.FirstURL, d.SecondURL),
			"Unhandled difference between `"+d.FirstURL+"` and `"+d.SecondURL+"`.")
	}
	return msg

}

// JSON jsonified diff message
func (d diffMessage) JSON() string {
	d.Status = "success"
	diffJSONBytes, e := json.Marshal(d)
	fatalIf(probe.NewError(e),
		"Unable to marshal diff message `"+d.FirstURL+"`, `"+d.SecondURL+"` and `"+string(d.Diff)+"`.")
	return string(diffJSONBytes)
}

func checkDiffSyntax(ctx *cli.Context, encKeyDB map[string][]prefixSSEPair) {
	if len(ctx.Args()) != 2 {
		cli.ShowCommandHelpAndExit(ctx, "diff", 1) // last argument is exit code
	}
	for _, arg := range ctx.Args() {
		if strings.TrimSpace(arg) == "" {
			fatalIf(errInvalidArgument().Trace(ctx.Args()...), "Unable to validate empty argument.")
		}
	}
	URLs := ctx.Args()
	firstURL := URLs[0]
	secondURL := URLs[1]

	// Diff only works between two directories, verify them below.

	// Verify if firstURL is accessible.
	_, firstContent, err := url2Stat(firstURL, false, encKeyDB)
	if err != nil {
		fatalIf(err.Trace(firstURL), fmt.Sprintf("Unable to stat '%s'.", firstURL))
	}

	// Verify if its a directory.
	if !firstContent.Type.IsDir() {
		fatalIf(errInvalidArgument().Trace(firstURL), fmt.Sprintf("`%s` is not a folder.", firstURL))
	}

	// Verify if secondURL is accessible.
	_, secondContent, err := url2Stat(secondURL, false, encKeyDB)
	if err != nil {
		fatalIf(err.Trace(secondURL), fmt.Sprintf("Unable to stat '%s'.", secondURL))
	}

	// Verify if its a directory.
	if !secondContent.Type.IsDir() {
		fatalIf(errInvalidArgument().Trace(secondURL), fmt.Sprintf("`%s` is not a folder.", secondURL))
	}
}

// doDiffMain runs the diff.
func doDiffMain(firstURL, secondURL string) error {
	// Source and targets are always directories
	sourceSeparator := string(newClientURL(firstURL).Separator)
	if !strings.HasSuffix(firstURL, sourceSeparator) {
		firstURL = firstURL + sourceSeparator
	}
	targetSeparator := string(newClientURL(secondURL).Separator)
	if !strings.HasSuffix(secondURL, targetSeparator) {
		secondURL = secondURL + targetSeparator
	}

	// Expand aliased urls.
	firstAlias, firstURL, _ := mustExpandAlias(firstURL)
	secondAlias, secondURL, _ := mustExpandAlias(secondURL)

	firstClient, err := newClientFromAlias(firstAlias, firstURL)
	if err != nil {
		fatalIf(err.Trace(firstAlias, firstURL, secondAlias, secondURL),
			fmt.Sprintf("Failed to diff '%s' and '%s'", firstURL, secondURL))
	}

	secondClient, err := newClientFromAlias(secondAlias, secondURL)
	if err != nil {
		fatalIf(err.Trace(firstAlias, firstURL, secondAlias, secondURL),
			fmt.Sprintf("Failed to diff '%s' and '%s'", firstURL, secondURL))
	}

	// Diff first and second urls.
	for diffMsg := range objectDifference(firstClient, secondClient, firstURL, secondURL) {
		if diffMsg.Error != nil {
			errorIf(diffMsg.Error, "Unable to calculate objects difference.")
			// Ignore error and proceed to next object.
			continue
		}
		printMsg(diffMsg)
	}

	return nil
}

// mainDiff main for 'diff'.
func mainDiff(ctx *cli.Context) error {
	// Parse encryption keys per command.
	encKeyDB, err := getEncKeys(ctx)
	fatalIf(err, "Unable to parse encryption keys.")

	// check 'diff' cli arguments.
	checkDiffSyntax(ctx, encKeyDB)

	// Additional command specific theme customization.
	console.SetColor("DiffMessage", color.New(color.FgGreen, color.Bold))
	console.SetColor("DiffOnlyInFirst", color.New(color.FgRed))
	console.SetColor("DiffOnlyInSecond", color.New(color.FgGreen))
	console.SetColor("DiffType", color.New(color.FgMagenta))
	console.SetColor("DiffSize", color.New(color.FgYellow, color.Bold))
	console.SetColor("DiffTime", color.New(color.FgYellow, color.Bold))

	URLs := ctx.Args()
	firstURL := URLs.Get(0)
	secondURL := URLs.Get(1)

	return doDiffMain(firstURL, secondURL)
}
