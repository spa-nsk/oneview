//demo for use information from HPE OneView
package main
import (
      github.com/spa-nsk/oneview
)

func main() {
	infra := GlobalInitOVInfrastructure()
	infra.AddEndpoint("https://172.17.100.100", "mydomain", "mydomain\\user", "password")
	infra.LoadServerHardwareList()

}
