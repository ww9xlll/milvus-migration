package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zilliztech/milvus-migration/core/gstore"
	"github.com/zilliztech/milvus-migration/core/util"
	"github.com/zilliztech/milvus-migration/internal/log"
	"github.com/zilliztech/milvus-migration/starter"
	"go.uber.org/zap"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start to migrate data",

	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()
		ctx := context.Background()

		jobId := util.GenerateUUID("start")
		fmt.Println("jodId is ", jobId)

		defer func() {
			if _any := recover(); _any != nil {
				handlePanic(_any, jobId)
				return
			}
		}()
		err := starter.Start(ctx, configFile, jobId)
		if err != nil {
			log.Error("[start migration error]", zap.Error(err))
			return
		}
		fmt.Printf("Migration Success! Job %s cost=[%f]\n", jobId, time.Since(start).Seconds())
		printJobMessage(jobId)
	},
}

func printJobMessage(jobId string) {
	jobInfo, _ := gstore.GetJobInfo(jobId)
	val, _ := json.Marshal(&jobInfo)
	fmt.Printf("Migration JobInfo: %s\n", string(val))

	procInfo := gstore.GetProcessHandler(jobId)
	val, _ = json.Marshal(&procInfo)
	fmt.Printf("Migration ProcessInfo: %s, Process:%d\n", string(val), procInfo.CalcProcess())

	fileTaskInfo := gstore.GetFileTask(jobId)
	val, _ = json.Marshal(&fileTaskInfo)
	fmt.Printf("Migration FileTaskInfo:  %s\n", string(val))
}

func handlePanic(_any any, jobId string) {
	var errMsg string
	err, ok := _any.(error)
	if ok {
		errMsg = err.Error()
	} else {
		errMsg, _ = _any.(string)
	}
	if err == nil {
		err = errors.New(errMsg)
	}
	fmt.Printf("Migration panic error! Job: %s , err: %s\n", jobId, errMsg)
	gstore.RecordJobError(jobId, err)
}

func init() {
	// ./milvus-migration start --config=/Users/zilliz/gitCode/cloud_team/milvus-migration/configs/migration_targetMinio.yaml
	rootCmd.AddCommand(startCmd)
}
