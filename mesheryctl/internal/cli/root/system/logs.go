// Copyright 2019 The Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs",
	Long:  `Print history of Meshery's container logs and begin tailing them.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if ok := utils.IsMesheryRunning(); !ok {
			log.Error("No logs to show. Meshery is not running.")
			return
		}

		log.Info("Starting Meshery logging...")

		if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
			if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
				log.Fatal("start cmd: ", err)
			}
		}

		cmdlog := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "logs", "-f")
		cmdReader, err := cmdlog.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		if err := cmdlog.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmdlog.Wait(); err != nil {
			log.Fatal(err)
		}
	},
}