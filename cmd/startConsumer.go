/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"

	"github.com/spf13/cobra"
)

// startConsumerCmd represents the startConsumer command
var startConsumerCmd = &cobra.Command{
	Use:   "startConsumer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("startConsumer called")
		startNatsConsumer()
	},
}

func init() {
	natsCmd.AddCommand(startConsumerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startConsumerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startConsumerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startNatsConsumer(){
	nc, _ := nats.Connect(nats.DefaultURL)
	sub,err := nc.SubscribeSync("time")
	if err != nil {
		panic(err)
	}
	for {
		msg,err := sub.NextMsg(time.Minute)
		if err != nil {
			panic(err)
		}
		var message TimeMessage
		err = json.Unmarshal(msg.Data,&message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("consumed message: %+v\n",message)
	}
}