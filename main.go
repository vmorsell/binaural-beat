package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bineural-beat",
		Usage: "play sounds that induce a binaural heat.",
		Commands: []*cli.Command{
			{
				Name:    "play",
				Aliases: []string{"p"},
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:    "frequency",
						Aliases: []string{"f"},
						Usage:   "Base frequency for the sound.",
						Value:   200,
					},
					&cli.StringFlag{
						Name:     "type",
						Aliases:  []string{"t"},
						Usage:    "Wave type (beta, alpha, theta or delta)",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					var beatFreq float64
					switch c.String("type") {
					case "beta":
						beatFreq = 14.5
					case "alpha":
						beatFreq = 10
					case "theta":
						beatFreq = 5.5
					case "delta":
						beatFreq = 1.75
					default:
						return fmt.Errorf("unknown wave type")
					}

					Play(c.Float64("frequency"), beatFreq)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
