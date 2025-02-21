/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package workflow

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"vitess.io/vitess/go/cmd/vtctldclient/cli"
	"vitess.io/vitess/go/cmd/vtctldclient/command/vreplication/common"

	vtctldatapb "vitess.io/vitess/go/vt/proto/vtctldata"
)

var (
	workflowDeleteOptions = struct {
		Workflow         string
		KeepData         bool
		KeepRoutingRules bool
	}{}

	// WorkflowDelete makes a WorkflowDelete gRPC call to a vtctld.
	workflowDelete = &cobra.Command{
		Use:                   "delete",
		Short:                 "Delete a VReplication workflow.",
		Example:               `vtctldclient --server localhost:15999 workflow --keyspace customer delete --workflow commerce2customer`,
		DisableFlagsInUseLine: true,
		Aliases:               []string{"Delete"},
		Args:                  cobra.NoArgs,
		RunE:                  commandWorkflowDelete,
	}
)

func commandWorkflowDelete(cmd *cobra.Command, args []string) error {
	cli.FinishedParsing(cmd)

	req := &vtctldatapb.WorkflowDeleteRequest{
		Keyspace:         workflowOptions.Keyspace,
		Workflow:         workflowDeleteOptions.Workflow,
		KeepData:         workflowDeleteOptions.KeepData,
		KeepRoutingRules: workflowDeleteOptions.KeepRoutingRules,
	}
	resp, err := common.GetClient().WorkflowDelete(common.GetCommandCtx(), req)
	if err != nil {
		return err
	}

	// Sort the inner TabletInfo slice for deterministic output.
	sort.Slice(resp.Details, func(i, j int) bool {
		return resp.Details[i].Tablet.String() < resp.Details[j].Tablet.String()
	})

	data, err := cli.MarshalJSON(resp)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)

	return nil
}

func addWorkflowDeleteFlags(cmd *cobra.Command) {
	workflowDelete.Flags().StringVarP(&workflowDeleteOptions.Workflow, "workflow", "w", "", "The workflow you want to delete (required).")
	workflowDelete.MarkFlagRequired("workflow")
	workflowDelete.Flags().BoolVar(&workflowDeleteOptions.KeepData, "keep-data", false, "Keep the partially copied table data from the workflow in the target keyspace.")
	workflowDelete.Flags().BoolVar(&workflowDeleteOptions.KeepRoutingRules, "keep-routing-rules", false, "Keep the routing rules created for the workflow.")

}
