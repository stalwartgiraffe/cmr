package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func NewPrjEventsCommand(cancel context.CancelFunc, cfg *CmdConfig) *cobra.Command {
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
			/*
				// Start tracing
				traceFile, err := os.Create("trace.out")
				if err != nil {
					panic(err)
				}
				defer traceFile.Close()

				if err := trace.Start(traceFile); err != nil {
					panic(err)
				}
				defer trace.Stop()

				traceCtx, cmdTask := trace.NewTask(cmdCtx, "prjEvents")
				defer cmdTask.End()
			*/

			// start := time.Now()
			accessToken, err := loadGitlabAccessToken()
			if err != nil {
				utils.Redln(err)
				return
			}

			ec := NewEventClient(accessToken)
			filepath := "ignore/my_recent_events.yaml"
			route := "events/"
			var logger AppLog
			logger = elog.New()
			myEvents, err := ec.updateRecentEvents(cmdCtx, logger, cancel, filepath, route)
			if err != nil {
				utils.Redln(err)
				return
			}
			myProjectIDs := myEvents.ProjectIDs()

			numWorkers := 200
			pendingIDs := make(chan int, numWorkers)

			logger = &elog.NoopLogger{}

			var wg sync.WaitGroup
			wg.Add(numWorkers)
			for worker := 0; worker < numWorkers; worker++ {
				// if we capture worker, rember to alias
				go func() {
					defer wg.Done()

					ids := []int{}
					for id := range pendingIDs {
						ids = append(ids, id)

						filepath := fmt.Sprintf("ignore/project_%d_events.yaml", id)
						route := fmt.Sprintf("projects/%d/events", id)
						// fmt.Println("worker", worker, "get ", id)
						_, err := ec.updateRecentEvents(cmdCtx, logger, cancel, filepath, route)
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
