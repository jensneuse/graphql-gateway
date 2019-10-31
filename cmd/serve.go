
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jensneuse/graphql-go-tools/pkg/execution"
	"github.com/jensneuse/graphql-go-tools/pkg/playground"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the gateway",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var (
	schemaFile string
	loggerConfigFile string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&schemaFile, "schema", "./schema.graphql", "schema is the configuration file")
	serveCmd.PersistentFlags().StringVar(&loggerConfigFile, "loggerConfig", "./logger.config.json", "configures the logger")
	_ = viper.BindPFlag("schema", rootCmd.PersistentFlags().Lookup("schema"))
	_ = viper.BindPFlag("loggerConfig", rootCmd.PersistentFlags().Lookup("loggerConfig"))
}

func logger () *zap.Logger {
	configData,err := ioutil.ReadFile(loggerConfigFile)
	if err != nil {
		panic(err)
	}
	var cfg zap.Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func startServer() {

	logger := logger()
	logger.Info("logger initialized")

	mux := http.NewServeMux()
	graphqlEndpoint := "/graphql"
	schemaData, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	handler, err := execution.NewHandler(schemaData,logger)
	if err != nil {
		log.Fatal(err)
	}
	mux.HandleFunc(graphqlEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		buf := bytes.NewBuffer(make([]byte, 0, 4096))
		err := handler.Handle(request.Body, buf)
		if err != nil {
			err := json.NewEncoder(writer).Encode(struct {
				Errors []struct {
					Message string `json:"message"`
				} `json:"errors"`
			}{
				Errors: []struct {
					Message string `json:"message"`
				}{
					{
						Message: err.Error(),
					},
				},
			})
			if err != nil {
				logger.Fatal("error encoding",zap.Error(err))
			}
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		_, _ = buf.WriteTo(writer)
	})
	playgroundURLPrefix := "/playground"
	playgroundURL := ""
	err = playground.ConfigureHandlers(mux, playground.Config{
		URLPrefix:       playgroundURLPrefix,
		PlaygroundURL:   playgroundURL,
		GraphqlEndpoint: graphqlEndpoint,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	addr := "localhost:9111"
	logger.Info("Listening",
		zap.String("add", addr),
	)
	fmt.Printf("Access Playground on: http://%s%s%s",addr,playgroundURLPrefix,playgroundURL)
	logger.Fatal("failed listening",
		zap.Error(http.ListenAndServe(addr, mux)),
	)
}
