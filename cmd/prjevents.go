package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func NewPrjEventsCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "prjevents",
		Short: "run prjevents",
		Long:  `Run Project Events.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			/*
				// Create a CPU profile file
				f, err := os.Create("profile.prof")
				if err != nil {
					panic(err)
				}
				defer f.Close()

				// Start CPU profiling
				if err := pprof.StartCPUProfile(f); err != nil {
					panic(err)
				}
				defer pprof.StopCPUProfile()
			*/

			cmdCtx := cmd.Context()

			// start := time.Now()
			accessToken, err := loadGitlabAccessToken()
			if err != nil {
				utils.Redln(err)
				return
			}

			ec := NewEventClient(accessToken)
			filepath := "ignore/my_recent_events.yaml"
			route := "events/"
			myEvents, err := ec.updateRecentEvents(cmdCtx, app, cancel, filepath, route)
			if err != nil {
				utils.Redln(err)
				return
			}
			myProjectIDs := myEvents.ProjectIDs()

			numWorkers := 200
			pendingIDs := make(chan int, numWorkers)

			var wg sync.WaitGroup
			wg.Add(numWorkers)
			for range numWorkers {
				// if we capture worker, rember to alias
				go func() {
					defer wg.Done()

					for id := range pendingIDs {

						filepath := fmt.Sprintf("ignore/project_%d_events.yaml", id)
						route := fmt.Sprintf("projects/%d/events", id)
						_, err := ec.updateRecentEvents(cmdCtx, app, cancel, filepath, route)
						if err != nil {
							utils.Redln(err)
							return
						}
					}

				}()
			}
			for _, id := range myProjectIDs {
				pendingIDs <- id
			}
			close(pendingIDs)
			wg.Wait()
			fmt.Println("done updating project_x_events")
		},
	}
}
