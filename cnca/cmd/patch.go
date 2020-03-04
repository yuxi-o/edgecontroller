// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use: "patch",
	Short: "Patch an active LTE CUPS userplane or NGC AF TI subscription or " +
		"NGC AF PFD Transaction or NGC AF PFD Application " +
		"using YAML configuration file",
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		ymlFile, _ := cmd.Flags().GetString("filename")
		if ymlFile == "" {
			fmt.Println(errors.New("YAML file missing"))
			return
		}

		data, err := ioutil.ReadFile(ymlFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		switch c.Kind {
		case "ngc":
			if pfdCommandCalled == true {
				fmt.Println(errors.New("Incorrect `kind` in YAML file"))
				return
			} else if len(args) > 1 {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}
			var s AFTrafficInfluSub
			if err = yaml.Unmarshal(data, &s); err != nil {
				fmt.Println(err)
				return
			}

			sub, err := yaml.Marshal(s.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			sub, err = y2j.YAMLToJSON(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			// patch subscription
			err = AFPatchSubscription(args[0], sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])

		case "lte":
			if pfdCommandCalled == true {
				fmt.Println(errors.New("Incorrect `kind` in YAML file"))
				return
			} else if len(args) > 1 {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}
			var u LTEUserplane
			if err = yaml.Unmarshal(data, &u); err != nil {
				fmt.Println(err)
				return
			}

			up, err := yaml.Marshal(u.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			up, err = y2j.YAMLToJSON(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			// patch userplane
			err = LtePatchUserplane(args[0], up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])
		case "ngc_pfd":

			if pfdCommandCalled == false {
				fmt.Println(errors.New("Incorrect `kind` in YAML file"))
				return
			}
			if len(args) < 2 {
				fmt.Println(errors.New("Missing input(s)"))
				return
			}

			if args[0] == "transaction" && args[1] != "" {
				if len(args) > 2 {
					if args[2] == "application" && len(args) > 3 {
						// patch PFD Application

						var s AFPfdData
						if err = yaml.Unmarshal(data, &s); err != nil {
							fmt.Println(err)
							return
						}

						pfdAppData := getPfdAppData(s)

						app, err := json.Marshal(pfdAppData)
						if err != nil {
							fmt.Println(err)
							return
						}

						err = AFPatchPfdApplication(args[1], args[3], app)
						if err != nil {
							klog.Info(err)
							return
						}
						fmt.Printf("Application %s patched\n", args[3])
						return
					}
				} else {
					// patch PFD Transaction

					var s AFPfdManagement
					if err = yaml.Unmarshal(data, &s); err != nil {
						fmt.Println(err)
						return
					}

					pfdTransData := getPfdTransData(s)

					trans, err := json.Marshal(pfdTransData)
					if err != nil {
						fmt.Println(err)
						return
					}

					err = AFPatchPfdTransaction(args[1], trans)
					if err != nil {
						klog.Info(err)
						return
					}
					fmt.Printf("Transaction %s patched\n", args[1])
					return
				}
			}
			fmt.Println(errors.New("Invalid input(s)"))
		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

func init() {

	const help = `Patch an active LTE CUPS userplane or NGC AF TI subscription or NGC AF PFD Transaction or NGC AF PFD Application using YAML configuration file

Usage:
  cnca {pfd | <none>} patch { <userplane-id> | <subscription-id> | transaction <transaction-id> {application <application-id> | <none>}  } -f <config.yml>

Example:
  cnca patch <userplane-id> -f <config.yml>
  cnca patch <subscription-id> -f <config.yml>
  cnca pfd patch transactions
  cnca pfd patch transaction <transaction-id> -f <config.yml>
  cnca pfd patch transaction <transaction-id> application <application-id> -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`
	// add `patch` command
	cncaCmd.AddCommand(patchCmd)
	pfdCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	applyCmd.MarkFlagRequired("filename")
	patchCmd.SetHelpTemplate(help)
}

func getPfdAppData(inputPfdAppData AFPfdData) PfdData {

	var pfdAppData PfdData

	pfdAppData.ExternalAppID = inputPfdAppData.Policy.ExternalAppID
	pfdAppData.Self = Link(inputPfdAppData.Policy.Self)

	if inputPfdAppData.Policy.AllowedDelay != nil {
		allowedDelay := DurationSecRm(*inputPfdAppData.Policy.AllowedDelay)
		pfdAppData.AllowedDelay = &allowedDelay
	}
	if inputPfdAppData.Policy.CachingTime != nil {
		cachingTime := DurationSecRo(*inputPfdAppData.Policy.CachingTime)
		pfdAppData.CachingTime = &cachingTime
	}
	if inputPfdAppData.Policy.Pfds != nil {
		pfdAppData.Pfds = make(map[string]Pfd)
	}

	for _, inputPfdData := range inputPfdAppData.Policy.Pfds {
		pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
	}

	return pfdAppData
}
