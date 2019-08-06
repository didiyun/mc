/*
 * Minio Client (C) 2016 Minio, Inc.
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
	"os"
	"strings"
	"sync"
	"syscall"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/didiyun/mc/pkg/probe"
	"github.com/didiyun/mc/pkg/console"
)

var (
	watchFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "events",
			Value: "put,delete,get",
			Usage: "filter specific types of events; defaults to all events by default",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "filter events for a prefix",
		},
		cli.StringFlag{
			Name:  "suffix",
			Usage: "filter events for a suffix",
		},
		cli.BoolFlag{
			Name:  "recursive",
			Usage: "recursively watch for events",
		},
	}
)

var watchCmd = cli.Command{
	Name:   "watch",
	Usage:  "listen for object notification events",
	Action: mainWatch,
	Before: setGlobalsFromContext,
	Flags:  append(watchFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] PATH
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
EXAMPLES:
   1. Watch new S3 operations on a minio server
      $ {{.HelpName}} play/testbucket

   2. Watch new events for a specific prefix "output/"  on minio server.
      $ {{.HelpName}} --prefix "output/" play/testbucket

   3. Watch new events for a specific suffix ".jpg" on minio server.
      $ {{.HelpName}} --suffix ".jpg" play/testbucket

   4. Watch new events on a specific prefix and suffix on minio server.
      $ {{.HelpName}} --suffix ".jpg" --prefix "photos/" play/testbucket

   5. Watch for events on local directory.
      $ {{.HelpName}} /usr/share
`,
}

// checkWatchSyntax - validate all the passed arguments
func checkWatchSyntax(ctx *cli.Context) {
	if len(ctx.Args()) != 1 {
		cli.ShowCommandHelpAndExit(ctx, "watch", 1) // last argument is exit code
	}
}

// watchMessage container to hold one event notification
type watchMessage struct {
	Status string `json:"status"`
	Event  struct {
		Time string    `json:"time"`
		Size int64     `json:"size"`
		Path string    `json:"path"`
		Type EventType `json:"type"`
	} `json:"events"`
	Source struct {
		Host      string `json:"host,omitempty"`
		Port      string `json:"port,omitempty"`
		UserAgent string `json:"userAgent,omitempty"`
	} `json:"source,omitempty"`
}

func (u watchMessage) JSON() string {
	u.Status = "success"
	watchMessageJSONBytes, e := json.Marshal(u)
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(watchMessageJSONBytes)
}

func (u watchMessage) String() string {
	msg := console.Colorize("Time", fmt.Sprintf("[%s] ", u.Event.Time))
	if u.Event.Type == EventCreate {
		msg += console.Colorize("Size", fmt.Sprintf("%6s ", humanize.IBytes(uint64(u.Event.Size))))
	} else {
		msg += fmt.Sprintf("%6s ", "")
	}
	msg += console.Colorize("EventType", fmt.Sprintf("%s ", u.Event.Type))
	msg += console.Colorize("ObjectName", fmt.Sprintf("%s", u.Event.Path))
	return msg
}

func mainWatch(ctx *cli.Context) error {
	console.SetColor("Time", color.New(color.FgGreen))
	console.SetColor("Size", color.New(color.FgYellow))
	console.SetColor("EventType", color.New(color.FgCyan, color.Bold))
	console.SetColor("ObjectName", color.New(color.Bold))

	checkWatchSyntax(ctx)

	args := ctx.Args()
	path := args[0]

	prefix := ctx.String("prefix")
	suffix := ctx.String("suffix")
	events := strings.Split(ctx.String("events"), ",")
	recursive := ctx.Bool("recursive")

	s3Client, pErr := newClient(path)
	if pErr != nil {
		fatalIf(pErr.Trace(), "Cannot parse the provided url.")
	}

	params := watchParams{
		recursive: recursive,
		events:    events,
		prefix:    prefix,
		suffix:    suffix,
	}

	// Start watching on events
	wo, err := s3Client.Watch(params)
	fatalIf(err, "Cannot watch on the specified bucket.")

	trapCh := signalTrap(os.Interrupt, syscall.SIGTERM)

	// Initialize.. waitgroup to track the go-routine.
	wg := sync.WaitGroup{}

	// Increment wait group to wait subsequent routine.
	wg.Add(1)

	// Start routine to watching on events.
	go func() {
		defer wg.Done()

		// Wait for all events.
		for {
			select {
			case <-trapCh:
				// Signal received we are done.
				close(wo.doneChan)
				return
			case event, ok := <-wo.Events():
				if !ok {
					return
				}
				msg := watchMessage{}
				msg.Event.Path = event.Path
				msg.Event.Size = event.Size
				msg.Event.Time = event.Time
				msg.Event.Type = event.Type
				msg.Source.Host = event.Host
				msg.Source.Port = event.Port
				msg.Source.UserAgent = event.UserAgent
				printMsg(msg)
			case err, ok := <-wo.Errors():
				if !ok {
					return
				}
				errorIf(err, "Unable to watch for events.")
				return
			}
		}
	}()

	// Wait on the routine to be finished or exit.
	wg.Wait()

	return nil
}
