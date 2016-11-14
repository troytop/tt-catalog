package sdltests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

type sdl struct {
	Name           string       `json:"name"`
	ProductVersion string       `json:"product_version"`
	SdlVersion     string       `json:"sdl_version"`
	Vendor         string       `json:"vendor"`
	Components     []component  `json:"components"`
	Parameters     []parameter  `json:"parameters"`
	PostFlight     []postFlight `json:"postflight"`
}

type component struct {
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Vendor       string               `json:"vendor"`
	Capabilities []string             `json:"capabilities"`
	Parameters   []componentParameter `json:"parameters"`
	ServicePorts []port               `json:"service_ports"`
}

type postFlight struct {
	Name       string               `json:"name"`
	RetryCount int                  `json:"retry_count"`
	Parameters []componentParameter `json:"parameters"`
}

type port struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Source   int    `json:"source"`
	Target   int    `json:"target"`
	Public   bool   `json:"publi"`
}

type componentParameter struct {
	Name string `json:"name"`
}

type parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

type serviceConfig struct {
	Description     string   `json:"description"`
	Categories      []string `json:"categories"`
	Labels          []string `json:"labels"`
	Type            string   `json:"type"`
	DeploymentGrade string   `json:"deployment_grade"`
	Supported       bool     `json:"supported"`
}

var portNameCharsetRegex = regexp.MustCompile("^[-a-z0-9]+$")
var portNameOneLetterRegexp = regexp.MustCompile("[a-z]")

func TestSDLsAreValid(t *testing.T) {
	servicesRoot := "../services"
	skipSDLValidation, ok := os.LookupEnv("SKIP_SDL_VALIDATION")
	if !ok {
		skipSDLValidation = ""
	}

	for _, vendor := range listDirs(t, servicesRoot) {
		vendorPath := servicesRoot + "/" + vendor.Name()

		if vendorPath == "../services/.DS_Store" {
			continue
		}

		//load all services
		for _, serviceName := range listDirs(t, vendorPath) {

			//every directory should be a service
			if !serviceName.IsDir() {
				continue
			}

			servicePath := filepath.Join(vendorPath, serviceName.Name())

			configFilePath := filepath.Join(vendorPath, serviceName.Name(), "config.json")
			configfile, configfileerr := ioutil.ReadFile(configFilePath)
			assert.Nil(t, configfileerr, "Error loading config file for service "+serviceName.Name())
			assert.True(t, len(configfile) > 0, "Found empty config file for service "+serviceName.Name())

			//check default version exists
			var c serviceConfig
			err := json.Unmarshal(configfile, &c)
			assert.Nil(t, err, "unmarshalling the service config failed for "+serviceName.Name())
			assert.NotEmpty(t, c.Type, "type must be defined")
			assert.NotEmpty(t, c.Description, "description must be defined")
			assert.NotEmpty(t, c.Categories, "categories must be defined")
			assert.NotEmpty(t, c.Labels, "labels must be defined")

			// DeploymentGrade should either be production or development or nil
			if c.DeploymentGrade != "" {
				deployGrade := strings.ToLower(c.DeploymentGrade)
				assert.True(t, deployGrade == "production" || deployGrade == "development", "deployment grade can be empty or development or production")
			}

			//subdirectories should be version of that service
			for _, productVersionDir := range listDirs(t, servicePath) {
				if !productVersionDir.IsDir() {
					continue
				}

				t.Log("Validating Service " + serviceName.Name() + " product version " + productVersionDir.Name())

				productVersionPath := filepath.Join(vendorPath, serviceName.Name(), productVersionDir.Name())

				for _, sdlVersionDir := range listDirs(t, productVersionPath) {
					if !sdlVersionDir.IsDir() {
						continue
					}

					t.Log("Validating Service " + serviceName.Name() + " product version " + productVersionDir.Name() + " sdl version " + sdlVersionDir.Name())

					sdlfile := filepath.Join(vendorPath, serviceName.Name(), productVersionDir.Name(), sdlVersionDir.Name(), "sdl/sdl.json")
					contents, fileerr := ioutil.ReadFile(sdlfile)
					t.Log("Reading Service" + sdlfile)

					//sdl should exist
					assert.Nil(t, fileerr, "Error loading sdl file for service "+serviceName.Name()+" version "+sdlVersionDir.Name())

					//sdl should not be empty
					assert.True(t, len(contents) > 0, "Found empty sdl file for service "+serviceName.Name()+" version "+sdlVersionDir.Name())

					//the sdl json object
					var s sdl
					err := json.Unmarshal(contents, &s)
					assert.Nil(t, err, "unmarshalling the sdl failed")
					//check contents match version and vendor directories
					assert.Equal(t, serviceName.Name(), s.Name, "service name does not match in sdl and directory")
					assert.Equal(t, vendor.Name(), s.Vendor, "vendor name does not match in sdl and directory")
					assert.Contains(t, productVersionDir.Name(), s.ProductVersion, "product_version version in sdl does not match dir its in")
					assert.Contains(t, sdlVersionDir.Name(), s.SdlVersion, "sdl_version in sdl does not match dir its in")
					assert.True(t, govalidator.IsSemver(s.SdlVersion), "SDL version specified in sdl is not of semver format")

					if strings.Contains(skipSDLValidation, serviceName.Name()) {
						t.Log("***** WARNING - " + serviceName.Name() + " - skipped SDL validation")
						continue
					}
					//convert sdl parameters to componentParameter type for check
					var sdlParams []componentParameter
					for _, param := range s.Parameters {

						sdlParams = append(sdlParams, componentParameter{Name: param.Name})

						//check that there is no default value set for parameters that contain 'PASSWORD' in the name
						if strings.Contains(strings.ToLower(param.Name), "password") {
							assert.Empty(t, param.Default, "SDL parameters that are passwords should not have a default value set")
						}
					}

					var compParams []componentParameter
					var compNames []string
					for _, comp := range s.Components {
						compNames = append(compNames, comp.Name)

						if comp.Name == "hsm-side-car" {
							p := componentParameter{Name: "SIDECAR_API_KEY"}
							assert.Contains(t, comp.Parameters, p, "SIDECAR_API_KEY is not present in hsm-side-car component")
							p = componentParameter{Name: "SIDECAR_CA_CERT"}
							assert.Contains(t, comp.Parameters, p, "SIDECAR_CA_CERT is not present in hsm-side-car component")
							p = componentParameter{Name: "tls-private-key-file"}
							assert.Contains(t, comp.Parameters, p, "tls-private-key-file is not present in hsm-side-car component")
							p = componentParameter{Name: "tls-cert-file"}
							assert.Contains(t, comp.Parameters, p, "tls-cert-file is not present in hsm-side-car component")
						}
						//log a warning if component has "ALL" capabilities set
						for _, cap := range comp.Capabilities {
							if cap == "ALL" {
								t.Log("***** WARNING - " + comp.Name + " - should verify that ALL capabilities is required")
							}
						}

						for _, port := range comp.ServicePorts {
							IsValidPortName(t, port.Name)
							assert.True(t, port.Source < 30000, "Port numbers should be less than 30000")
							assert.True(t, port.Target < 30000, "Port numbers should be less than 30000")
						}

						//check component parameters exist in parameter definition list
						for _, param := range comp.Parameters {
							assert.Contains(t, sdlParams, param, "component parameter not found in sdl parameter definition list")
							compParams = append(compParams, param)
						}
					}

					for _, postFlight := range s.PostFlight {
						for _, param := range postFlight.Parameters {
							assert.Contains(t, sdlParams, param, "component parameter not found in sdl parameter definition list")
							compParams = append(compParams, param)
						}
					}

					//check the parameter name and description are not equal and that all SDL parameters are used in components
					for _, param := range s.Parameters {
						assert.NotEqual(t, param.Name, param.Description, "Parameter name and description should not be the same")
						assert.Contains(t, compParams, componentParameter{Name: param.Name}, "SDL parameter defined but not used in a component")
					}

					//load the sdl and validate it against the schema
					sdlLoader := gojsonschema.NewReferenceLoader("file://" + sdlfile)
					sdlSchemaLoader := gojsonschema.NewReferenceLoader("file://./sdl_schema.json")
					result, err := gojsonschema.Validate(sdlSchemaLoader, sdlLoader)
					assert.Nil(t, err, "There was an unexpected error")
					for _, err := range result.Errors() {
						t.Log(err)
					}
					assert.Equal(t, len(result.Errors()), 0, "There was an error validating the SDL schema")
				}
			}
		}
	}
}

// https://github.com/kubernetes/kubernetes/blob/master/pkg/util/validation/validation.go
// IsValidPortName check that the argument is valid syntax. It must be non empty and no more than 15 characters long
// It must contain at least one letter [a-z] and it must contain only [a-z0-9-].
// Hyphens ('-') cannot be leading or trailing character of the string and cannot be adjacent to other hyphens.
// Although RFC 6335 allows upper and lower case characters but case is ignored for comparison purposes: (HTTP
// and http denote the same service).
func IsValidPortName(t *testing.T, port string) {
	assert.False(t, len(port) > 15, "Port name length shoudl not be more than 15")

	assert.True(t, portNameCharsetRegex.MatchString(port), "must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)")

	assert.True(t, portNameOneLetterRegexp.MatchString(port), "must contain at least one letter (a-z)")

	assert.False(t, strings.Contains(port, "--"), "must not contain consecutive hyphens")

	assert.False(t, port[0] == '-', "must not begin with a hyphen")

	assert.False(t, len(port) > 0 && port[len(port)-1] == '-', "must not end with a hyphen")
}

func listDirs(t *testing.T, path string) []os.FileInfo {
	dirList, err := ioutil.ReadDir(path)
	assert.Nil(t, err, "Error listing services directory")
	return dirList
}
