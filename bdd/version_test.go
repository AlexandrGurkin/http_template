package bdd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlexandrGurkin/http_template/client/api/version"
	"github.com/cucumber/godog"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func NewVersionDriver() VersionDriver {
	url := os.Getenv("VERSION_URL")
	clientTransport := httptransport.New(url, "/api/v1", []string{"http"})
	api := version.New(clientTransport, strfmt.Default)
	return VersionDriver{url: url, api: api}
}

type VersionDriver struct {
	url            string
	api            version.ClientService
	lastResVersion *version.GetVersionOK
	lastErrVersion error
}

func (d *VersionDriver) iSendRequestToVersionHandler() error {
	d.lastResVersion, d.lastErrVersion = d.api.GetVersion(version.NewGetVersionParams())
	if d.lastErrVersion != nil {
		if _, ok := d.lastErrVersion.(*runtime.APIError); ok {
			return nil
		}
	}
	return d.lastErrVersion
}

func (d *VersionDriver) theResponseCodeShouldBe(arg1 int) error {
	tmpCode := 0
	if d.lastResVersion != nil {
		tmpCode = 200
	} else {
		if e, ok := d.lastErrVersion.(*runtime.APIError); ok {
			tmpCode = e.Code
		}
	}

	if tmpCode == arg1 {
		return nil
	} else {
		return fmt.Errorf("expected response code to be: [%d], but actual is: [%d]", arg1, tmpCode)
	}
}

func (d *VersionDriver) theResponseShouldMatchJson(arg1 string) error {
	if d.lastResVersion != nil {
		if len(d.lastResVersion.Payload.Version) > 0 {
			return nil
		} else {
			str, err := json.Marshal(d.lastResVersion.Payload)
			if err != nil {
				return err
			}
			return fmt.Errorf("expected resp with \"version\" field, bun actual is [%s]", string(str))
		}
	}
	return fmt.Errorf("response is nil")
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	driver := NewVersionDriver()

	ctx.Step(`^I send request to version handler$`, driver.iSendRequestToVersionHandler)
	ctx.Step(`^the response code should be (\d+)$`, driver.theResponseCodeShouldBe)
	ctx.Step(`^the response should contain field: "([^"]*)"$`, driver.theResponseShouldMatchJson)
}
