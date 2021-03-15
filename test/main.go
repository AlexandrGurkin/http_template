package main

import (
	"os"
	"time"

	"github.com/AlexandrGurkin/common/consts"
	"github.com/AlexandrGurkin/common/middlewares"
	"github.com/AlexandrGurkin/common/xlog"
	"github.com/AlexandrGurkin/common/xlog/xlogrus"
	"github.com/AlexandrGurkin/http_template/cmd"
	"github.com/AlexandrGurkin/http_template/internal/handlers"
	"github.com/AlexandrGurkin/http_template/restapi"
	"github.com/AlexandrGurkin/http_template/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/spf13/cobra"

	_ "net/http/pprof"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run callback",
	Long:  `run callback`,
	Run:   runMain,
}

func init() {
	cmd.RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runMain(cmd *cobra.Command, args []string) {
	logger := xlogrus.NewXLogrus(xlog.LoggerCfg{
		Level: "trace",
		Out:   os.Stdout,
	})

	var err error
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logger.Fatal(err)
	}

	api := operations.NewTemplateForHTTPServerAPI(swaggerSpec)
	server := restapi.NewServer(api)

	//cfg := configuration.Config{}
	//err = viper.Unmarshal(&cfg)
	//if err != nil {
	//	log.WithFields(log.Fields{consts.FieldModule: runnerModule, consts.FieldAction: consts.ConfigParsing}).
	//		Fatalf("fail to unmarshal config [%s]", err.Error())
	//}

	server.Host = "0.0.0.0"
	server.Port = 8099
	restapi.SetMiddlewareConfig(middlewares.MiddlewareConfig{Logger: logger, Pprof: true})
	api.Logger = func(s string, i ...interface{}) {
		logger.WithXFields(xlog.Fields{consts.FieldModule: "swagger_api_loger"}).
			Infof(s, i...)
	}
	api.VersionGetVersionHandler = handlers.VersionHandler{}
	server.ConfigureAPI()
	server.KeepAlive = 10 * time.Second

	if err = server.Serve(); err != nil {
		logger.Fatal(err)
	}
}

func main() {
	cmd.Execute()
}
