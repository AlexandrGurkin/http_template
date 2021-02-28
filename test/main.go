package main

import (
	"os"

	"github.com/AlexandrGurkin/common/consts"
	"github.com/AlexandrGurkin/common/middlewares"
	"github.com/AlexandrGurkin/common/xlog"
	"github.com/AlexandrGurkin/common/xlog/xlogrus"
	"github.com/AlexandrGurkin/http_template/restapi"
	"github.com/AlexandrGurkin/http_template/restapi/operations"
	"github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTemplateForHTTPServerAPI(swaggerSpec)
	server := restapi.NewServer(api)

	//cfg := configuration.Config{}
	//err = viper.Unmarshal(&cfg)
	//if err != nil {
	//	log.WithFields(log.Fields{consts.FieldModule: runnerModule, consts.FieldAction: consts.ConfigParsing}).
	//		Fatalf("fail to unmarshal config [%s]", err.Error())
	//}

	logger := xlogrus.NewXLogrus(xlog.LoggerCfg{
		Level: "trace",
		Out:   os.Stdout,
	})

	server.Host = "0.0.0.0"
	server.Port = 8022
	restapi.SetMiddlewareConfig(middlewares.MiddlewareConfig{Logger: logger})
	api.Logger = func(s string, i ...interface{}) {
		logger.WithXFields(xlog.Fields{consts.FieldModule: "swagger_api_loger"}).
			Debugf(s, i...)
	}
	server.ConfigureAPI()

	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
