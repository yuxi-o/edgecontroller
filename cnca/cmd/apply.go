// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use: "apply",
	Short: "Apply LTE CUPS userplane or NGC AF TI subscription or NGC AF PFD" +
		"Transaction or NGC AF PFD Application using YAML configuration file",
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

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

			// create new subscription
			subLoc, err := AFCreateSubscription(sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Subscription URI:", subLoc)

		case "lte":
			if pfdCommandCalled == true {
				fmt.Println(errors.New("Incorrect `kind` in YAML file"))
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

			// create new LTE userplane
			upID, err := LteCreateUserplane(up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Userplane:", upID)

		case "ngc_pfd":
			if pfdCommandCalled == false {
				fmt.Println(errors.New("Incorrect `kind` in YAML file"))
				return
			}
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

			// create new AF PFD Transaction
			pfdData, self, err := AFCreatePfdTransaction(trans)
			if err != nil {
				klog.Info(err)
				return
			}

			if pfdData != nil {
				//Convert the json PFD Transaction data into struct
				pfdTrans := PfdManagement{}
				err = json.Unmarshal(pfdData, &pfdTrans)
				if err != nil {
					klog.Info(err)
					return
				}
				if self == "" {
					self = string(pfdTrans.Self)
				}

				fmt.Printf("PFD Transaction URI: %s\n", self)
				fmt.Printf("PFD Transaction ID: %s\n",
					getTransIdFromUrl(self))
				fmt.Println("    Application IDs:")

				var appStatus map[string]string

				appStatus = make(map[string]string)
				for k, _ := range pfdTrans.PfdDatas {
					appStatus[k] = "Created"
				}
				for _, v := range pfdTrans.PfdReports {
					for _, str := range PfdReport(v).ExternalAppIds {
						appStatus[str] = string(PfdReport(v).FailureCode)
					}
				}
				for k, v := range appStatus {
					if v != "Created" {
						fmt.Printf("      - %s : Failed (Reason: %s)\n", k, v)
					} else {
						fmt.Printf("      - %s : %s\n", k, v)
					}
				}
			} else {
				fmt.Printf("PFD Transaction URI: %s\n", self)
				fmt.Printf("PFD Transaction ID: %s\n",
					getTransIdFromUrl(self))
			}
		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

func init() {

	const help = `Apply LTE CUPS userplane or NGC AF TI subscription or NGC AF PFD Transaction or NGC AF PFD Application using YAML configuration file

Usage:
  cnca {pfd | <none>} apply -f <config.yml>

Example:
  cnca apply -f <config.yml>
  cnca pfd apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`
	// add `apply` command
	cncaCmd.AddCommand(applyCmd)
	pfdCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	applyCmd.MarkFlagRequired("filename")
	applyCmd.SetHelpTemplate(help)
}

func getTransIdFromUrl(url string) string {
	urlElements := strings.Split(url, "/")
	for index, str := range urlElements {
		if str == "transactions" {
			return urlElements[index+1]
		}
	}
	return ""
}

func getPfdTransData(inputPfdTransData AFPfdManagement) PfdManagement {

	var pfdTransData PfdManagement

	if inputPfdTransData.Policy.SuppFeat != nil {
		*pfdTransData.SuppFeat =
			SupportedFeatures(*inputPfdTransData.Policy.SuppFeat)
	}

	if inputPfdTransData.Policy.PfdDatas != nil {
		pfdTransData.PfdDatas = make(map[string]PfdData)
	}

	for _, inputPfdAppData := range inputPfdTransData.Policy.PfdDatas {
		var pfdAppData PfdData

		pfdAppData.ExternalAppID = inputPfdAppData.ExternalAppID
		pfdAppData.Self = Link(inputPfdAppData.Self)

		if inputPfdAppData.AllowedDelay != nil {
			allowedDelay := DurationSecRm(*inputPfdAppData.AllowedDelay)
			pfdAppData.AllowedDelay = &allowedDelay
		}
		if inputPfdAppData.CachingTime != nil {
			cachingTime := DurationSecRo(*inputPfdAppData.CachingTime)
			pfdAppData.CachingTime = &cachingTime
		}
		if inputPfdAppData.Pfds != nil {
			pfdAppData.Pfds = make(map[string]Pfd)
		}

		for _, inputPfdData := range inputPfdAppData.Pfds {
			pfdAppData.Pfds[inputPfdData.PfdID] = Pfd(inputPfdData)
		}
		pfdTransData.PfdDatas[pfdAppData.ExternalAppID] = pfdAppData
	}
	return pfdTransData
}
