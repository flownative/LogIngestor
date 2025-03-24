package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
	"github.com/papertrail/go-tail/follower"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
)

var version = "undefined"

func main() {
	var logfile, endpoint, streamName, ruleId string
	var tee bool

	app := &cli.App{
		Name:        "logingestor",
		HelpName:    "logingestor",
		Usage:       "Tail a log file and send lines to Azure Monitor Logs",
		Description: "A tool to tail a log file and send lines to Azure Monitor Logs",
		Version:     version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Usage:       "The `ENDPOINT` to send data to",
				EnvVars:     []string{"AZURE_MONITOR_ENDPOINT"},
				Destination: &endpoint,
			},
			&cli.StringFlag{
				Name:        "rule-id",
				Usage:       "The `RULE-ID` to send data to",
				EnvVars:     []string{"AZURE_MONITOR_RULE_ID"},
				Destination: &ruleId,
			},
			&cli.StringFlag{
				Name:        "stream-name",
				Usage:       "The `STREAM-NAME` to send data to",
				EnvVars:     []string{"AZURE_MONITOR_STREAM_NAME"},
				Destination: &streamName,
			},
			&cli.BoolFlag{
				Name:        "tee",
				Value:       false,
				Usage:       "If set, processed entries are output to stdout",
				Destination: &tee,
			},
		},
		Action: func(cCtx *cli.Context) error {
			cred, err := azidentity.NewDefaultAzureCredential(nil)
			if err != nil {
				// TODO: handle error better
				log.Fatalf("failed to obtain a credential: %v", err)
			}

			client, err := azlogs.NewClient(endpoint, cred, nil)
			if err != nil {
				// TODO: handle error better
				log.Fatalf("failed to obtain a client: %v", err)
			}

			logfile = cCtx.Args().Get(0)
			fmt.Println(fmt.Sprintf("tailing: %s", logfile))
			fmt.Println(fmt.Sprintf("to: %s/dataCollectionRules/%s/streams/%s", endpoint, ruleId, streamName))

			t, err := follower.New(logfile, follower.Config{
				Whence: io.SeekEnd,
				Offset: 0,
				Reopen: true,
			})

			if err != nil {
				// TODO: handle error better
				log.Fatalf("failed to obtain a follower: %v", err)
			}

			for line := range t.Lines() {
				// wrap line in [] so the ingestion accepts it
				wrappedLine := "[" + line.String() + "]"

				if tee {
					fmt.Println(wrappedLine)
				}

				// upload logs
				_, err = client.Upload(context.TODO(), ruleId, streamName, []byte(wrappedLine), nil)

				if err != nil {
					// TODO: handle error better
					log.Fatalf("failed to upload: %v", err)
				}
			}

			if t.Err() != nil {
				// TODO: handle error better
				log.Fatalf("failed to follow: %v", t.Err())
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		// TODO: handle error better
		log.Fatalf("boom: %v", err)
	}
}
