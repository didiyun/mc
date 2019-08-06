package cmd

import (
	"strings"

	"github.com/didiyun/mc/pkg/console"
	"github.com/fatih/color"
	"github.com/minio/cli"
)

// ck specific flags.
var (
	checkkFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "repair",
			Usage: "repair objects",
		},
	}
)

// show object metadata
var checkCmd = cli.Command{
	Name:   "check",
	Usage:  "check big object content",
	Action: mainCheck,
	Before: setGlobalsFromContext,
	Flags:  append(checkkFlags, globalFlags...),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] TARGET [TARGET ...]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

ENVIRONMENT VARIABLES:
   MC_ENCRYPT_KEY: List of comma delimited prefix=secret values

EXAMPLES:
   1. Check object content md5 of mybucket on DIDIYUN S3 cloud storage.
      $ {{.HelpName}} s3/mybucket/object ./localfile
`,
}

// checkStatSyntax - validate all the passed arguments
func checkSyntax(ctx *cli.Context, encKeyDB map[string][]prefixSSEPair) {
	if !ctx.Args().Present() {
		cli.ShowCommandHelpAndExit(ctx, "check", 1) // last argument is exit code
	}

	args := ctx.Args()
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			fatalIf(errInvalidArgument().Trace(args...), "Unable to validate empty argument.")
		}
	}
}

// mainCheck - is a handler for mc stat command
func mainCheck(ctx *cli.Context) error {
	// Additional command specific theme customization.
	console.SetColor("Name", color.New(color.Bold, color.FgCyan))
	console.SetColor("Date", color.New(color.FgWhite))
	console.SetColor("Size", color.New(color.FgWhite))
	console.SetColor("ETag", color.New(color.FgWhite))
	console.SetColor("EncryptionHeaders", color.New(color.FgWhite))
	console.SetColor("Metadata", color.New(color.FgWhite))
	// check 'check' cli arguments.
	encKeyDB, cErr := getEncKeys(ctx)
	fatalIf(cErr, "Unable to parse encryption keys.")
	// 检查参数
	checkSyntax(ctx, encKeyDB)
	args := ctx.Args()
	// mimic operating system tool behavior.
	if !ctx.Args().Present() {
		args = []string{"."}
	}
	// Set command flags from context.
	isRepair := ctx.Bool("repair")
	sourceURL := args[0]
	localFile := args[1]

	// 获取状态
	objectParts, err := getObjectPartList(sourceURL, false, encKeyDB)
	if cErr != nil {
		return err
	}
	checkResultList, err := doCheck(objectParts, localFile)
	if err != nil {
		return err
	}

	if len(checkResultList) > 0 {
		color.Blue("\nAll Wrong Info\n")
		for _, errlist := range checkResultList {
			color.Red("Wrong: PartNumber[%d] LocalMD5[%s] RemoteMd5[%s] Size[%d] LastModified[%s]\n",
				errlist.PartNumber, errlist.SrcMD5, errlist.ETag, errlist.Size, errlist.LastModified)
		}
	}

	color.Blue("\nALL Total: Right[%d] Wrong[%d]\n", len(objectParts)-len(checkResultList), len(checkResultList))
	// 修复文件
	if isRepair == false {
		return nil
	}

	objectSize := getObjectSize(objectParts)
	if err = doRepair(sourceURL, checkResultList, localFile, objectSize); err != nil {
		color.Red("\n%s\n", err.Error())
		return err
	}
	color.Blue("\nRepair succedd!\n")

	return nil
}
