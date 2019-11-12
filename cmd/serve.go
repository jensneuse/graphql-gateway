package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/jensneuse/graphql-go-tools/pkg/execution"
	gqlHTTP "github.com/jensneuse/graphql-go-tools/pkg/http"
	"github.com/jensneuse/graphql-go-tools/pkg/playground"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
	schemaFile       string
	loggerConfigFile string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&schemaFile, "schema", "./schema.graphql", "schema is the configuration file")
	serveCmd.PersistentFlags().StringVar(&loggerConfigFile, "loggerConfig", "./logger.config.json", "configures the logger")
	_ = viper.BindPFlag("schema", rootCmd.PersistentFlags().Lookup("schema"))
	_ = viper.BindPFlag("loggerConfig", rootCmd.PersistentFlags().Lookup("loggerConfig"))
}

func logger() *zap.Logger {
	configData, err := ioutil.ReadFile(loggerConfigFile)
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
	handler, err := execution.NewHandler(schemaData, logger)
	if err != nil {
		log.Fatal(err)
	}

	upgrader := &ws.DefaultHTTPUpgrader
	upgrader.Header = http.Header{}
	upgrader.Header.Add("Sec-Websocket-Protocol", "graphql-ws")
	mux.HandleFunc("/time", func(writer http.ResponseWriter, request *http.Request) {
		_,err := writer.Write(fakeResponse())
		if err != nil {
			logger.Error("time_write_err",
				zap.Error(err),
			)
		}
	})
	mux.Handle(graphqlEndpoint, gqlHTTP.NewGraphqlHTTPHandlerFunc(handler,logger,upgrader))
	playgroundURLPrefix := "/playground"
	playgroundURL := ""
	err = playground.ConfigureHandlers(mux, playground.Config{
		URLPrefix:       playgroundURLPrefix,
		PlaygroundURL:   playgroundURL,
		GraphqlEndpoint: graphqlEndpoint,
		GraphQLSubscriptionEndpoint:graphqlEndpoint,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	addr := "0.0.0.0:9111"
	logger.Info("Listening",
		zap.String("add", addr),
	)
	fmt.Printf("Access Playground on: http://%s%s%s\n",prettyAddr(addr) , playgroundURLPrefix, playgroundURL)
	logger.Fatal("failed listening",
		zap.Error(http.ListenAndServe(addr, mux)),
	)
}

func fakeResponse () []byte {
	return []byte(`{"week_number":45,"utc_offset":"+01:00","utc_datetime":"2019-11-07T14:02:02.475928+00:00","unixtime":1573135322,"timezone":"Europe/Berlin","raw_offset":3600,"dst_until":null,"dst_offset":0,"dst_from":null,"dst":false,"day_of_year":311,"day_of_week":4,"datetime":"`+ time.Now().String() +`","client_ip":"92.216.144.100","abbreviation":"CET"}`)
}

func prettyAddr(addr string) string {
	return strings.Replace(addr,"0.0.0.0","localhost",-1)
}