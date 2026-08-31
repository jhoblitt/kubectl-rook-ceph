/*
Copyright 2026 The Rook Authors. All rights reserved.

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

package command

import (
	"fmt"
	"strings"

	"github.com/rook/kubectl-rook-ceph/pkg/exec"
	"github.com/rook/kubectl-rook-ceph/pkg/filesystem"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/validation"
)

var ExecCmd = &cobra.Command{
	Use:   "exec [component]",
	Short: "Open an interactive shell in a Rook Ceph component",
	Long: "Open an interactive shell in the Rook Ceph operator, toolbox, or a Ceph daemon pod.\n\n" +
		"The component defaults to toolbox and can also be operator, mon-<id>, mgr-<id>, osd-<id>, mds-<id>, or rgw-<id>.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		_, _, err := execTarget(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		labelSelector, container, err := execTarget(args)
		if err != nil {
			logging.Fatal(err)
		}
		namespace := cephClusterNamespace
		if container == "rook-ceph-operator" {
			namespace = operatorNamespace
		}
		if err := exec.RunShellInLabeledPod(cmd.Context(), clientSets, namespace, labelSelector, container); err != nil {
			logging.Fatal(err)
		}
	},
}

func execTarget(args []string) (string, string, error) {
	component := "toolbox"
	if len(args) == 1 {
		component = args[0]
	}

	if component == "toolbox" {
		return "app=rook-ceph-tools", "rook-ceph-tools", nil
	}
	if component == "operator" {
		return "app=rook-ceph-operator", "rook-ceph-operator", nil
	}

	daemonType, daemonID, found := strings.Cut(component, "-")
	if !found || daemonID == "" || len(validation.IsValidLabelValue(daemonID)) != 0 {
		return "", "", invalidExecComponent(component)
	}

	switch daemonType {
	case "mon":
		return fmt.Sprintf("app=rook-ceph-mon,mon=%s,mon_canary!=true", daemonID), daemonType, nil
	case "mgr", "osd", "mds", "rgw":
		return fmt.Sprintf("app=rook-ceph-%s,%s=%s", daemonType, daemonType, daemonID), daemonType, nil
	default:
		return "", "", invalidExecComponent(component)
	}
}

func invalidExecComponent(component string) error {
	return fmt.Errorf("invalid component %q: expected toolbox, operator, mon-<id>, mgr-<id>, osd-<id>, mds-<id>, or rgw-<id>", component)
}

func addCustomExecFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("pod-name", "", "Pod to execute commands in")
	cmd.PersistentFlags().String("pod-namespace", "", "Namespace of the target pod")
	cmd.PersistentFlags().String("pod-container", "", "Container in the target pod")
	cmd.PersistentFlags().String("mon-ip", "", "Ceph monitor IP (e.g. 10.0.0.1:6789)")
	cmd.PersistentFlags().String("user-id", "", "Ceph user ID for authentication")
	cmd.PersistentFlags().String("user-key", "", "Ceph user key for authentication")
}

func parseCustomExecConfig(cmd *cobra.Command) (*filesystem.CustomExecConfig, error) {
	pod, _ := cmd.Flags().GetString("pod-name")
	if pod == "" {
		return nil, nil
	}
	ns, _ := cmd.Flags().GetString("pod-namespace")
	container, _ := cmd.Flags().GetString("pod-container")
	monIP, _ := cmd.Flags().GetString("mon-ip")
	userID, _ := cmd.Flags().GetString("user-id")
	userKey, _ := cmd.Flags().GetString("user-key")

	if ns == "" || container == "" || monIP == "" || userID == "" || userKey == "" {
		return nil, fmt.Errorf(
			"--pod-namespace, --pod-container, --mon-ip, --user-id, and --user-key are all required when --pod-name is set")
	}
	return &filesystem.CustomExecConfig{
		PodName:      pod,
		PodNamespace: ns,
		Container:    container,
		MonIP:        monIP,
		UserID:       userID,
		UserKey:      userKey,
	}, nil
}
