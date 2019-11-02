//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"quebic-faas/common"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	setupIngressCmds()
	setupIngressFlags()
}

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Ingress commonds",
	Long:  `Ingress commonds`,
}

func setupIngressCmds() {
	ingressCmd.AddCommand(ingressDescribeCmd)
}

func setupIngressFlags() {
}

var ingressDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "ingress : describe",
	Long:  `ingress : describe`,
	Run: func(cmd *cobra.Command, args []string) {
		ingressInfo(cmd, args)
	},
}

func ingressInfo(cmd *cobra.Command, args []string) {

	deployment, err := getDeployment()
	if err != nil {
		prepareError(cmd, err)
	}

	details, err := deployment.ListByName(quebicManagerComponentID)
	if err != nil {
		prepareError(cmd, err)
	}

	if details.Status == common.KubeStatusFalse {
		prepareError(cmd, fmt.Errorf("manager not ready yet"))
	}

	ingressDetails, err := deployment.IngressDescribe(waitForAvailable)
	if err != nil {
		prepareError(cmd, fmt.Errorf("manager not ready yet"))
	}

	ingressIP := ingressDetails.IP

	var comps [][]string
	comps = append(comps, []string{quebicManagerComponentID, ingressIP, common.IngressHostManager})
	comps = append(comps, []string{common.ComponentAPIGateway, ingressIP, common.IngressHostAPIGateway})

	//prepare table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Component",
		"Address",
		"Host",
	})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(comps)
	table.Render()

}
