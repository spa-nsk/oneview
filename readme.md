Пример использования HPE OneView

Загрузка демонстрационного пакета
```bash
$go get github.com/spa-nsk/oneview
```

ServerHardware  - структура описывающая сервер в OneView
```
type ServerHardware struct {
	Base    ov.ServerHardware		//базовые данные по серверу
	Memory  ServerHardwareMemory		//информация о модулях памяти
	Storage ServerHardwareLocalStorage	//информация о локальных хранилищах (дисках)
}
```
OVInfrastructure  -  структура описывающая объекты OneView
```
type OVInfrastructure struct {
	endpoints    []*ovEndpoint		//серверы с настроенным OneView
	Servers      []*ServerHardware		//данные по загруженным серверам
	ServersCount int			//количество серверов
}
```
получить данные сервера по серийному номеру
```
func (infra *OVInfrastructure) FindServerHardwareSN(sn string) (*ServerHardware, error)
```
Пример использования

```
package main
import (
	"fmt"
      	"github.com/spa-nsk/oneview"
)

func main() {
	infra := oneview.GlobalInitOVInfrastructure()
	infra.AddEndpoint("https://172.17.100.100", "mydomain", "mydomain\\user", "password")
	infra.LoadServerHardwareList()
	srv, err := infra.FindServerHardwareSN("CZ28510H7T")
	if err == nil{
		fmt.Println("Server ",srv.Base.Model, srv.Base.ProcessorCount,"x", srv.Base.ProcessorType, srv.Base.MemoryMb, "Mb RAM")
		for i, mem := range srv.Memory.Data {
			fmt.Println(i+1,mem.CapacityMiB,"Mb", mem.Manufacturer,"(", mem.PartNumber,")", mem.MemoryDeviceType,mem.OperatingSpeedMhz,"Mhz")
		}
	}
}

```
