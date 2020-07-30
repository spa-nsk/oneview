package oneview

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/HewlettPackard/oneview-golang/ov"
)

type ovEndpoint struct {
	domain   string
	login    string
	password string
	endpoint string
}

//ServerHardware  - структура описывающая сервер в OneView
type ServerHardware struct {
	Base    ov.ServerHardware
	Memory  ServerHardwareMemory
	Storage ServerHardwareLocalStorage
}

//OVInfrastructure  -  структура описывающая объекты OneView
type OVInfrastructure struct {
	endpoints    []*ovEndpoint
	Servers      []*ServerHardware
	ServersCount int
}

//AddEndpoint  - функция добавления точки подключения к списку подключений
func (infra *OVInfrastructure) AddEndpoint(endpoint string, domain string, login string, password string) {
	infra.endpoints = append(infra.endpoints, &ovEndpoint{endpoint: endpoint, domain: domain, login: login, password: password})
}

var infrastructureGlobal *OVInfrastructure

//Init  - функция инициализации глобальной структуры
func (infra *OVInfrastructure) Init() {
	if infra == nil {
		infrastructureGlobal = &OVInfrastructure{}
		infra = infrastructureGlobal //проверить корректность обработки
	}
	infra.Servers = make([]*ServerHardware, 0)
	infra.ServersCount = 0
	infra.endpoints = make([]*ovEndpoint, 0)
}

//Destroy  - функция разрушения глобальной структуры
func (infra *OVInfrastructure) Destroy() {
	if infra == nil {
		infrastructureGlobal = &OVInfrastructure{}
		infra = infrastructureGlobal //проверить корректность обработки
	}
	infra.Servers = make([]*ServerHardware, 0)
	infra.ServersCount = 0
	infra.endpoints = make([]*ovEndpoint, 0)
}

//GlobalInitOVInfrastructure  - инициализация пакета
func GlobalInitOVInfrastructure() *OVInfrastructure {

	infrastructureGlobal = &OVInfrastructure{}
	infrastructureGlobal.Init()

	return infrastructureGlobal
}

//GetGlobalOVInfrastructure  - получение ссылки глобальную структуру
func GetGlobalOVInfrastructure() *OVInfrastructure {
	if infrastructureGlobal == nil {
		infrastructureGlobal = &OVInfrastructure{}
		infrastructureGlobal.Init()
	}
	return infrastructureGlobal
}

//DestroyGlobalOVInfrastructureObj  - уничтожение глобальгной структуры
func DestroyGlobalOVInfrastructureObj() {
	if infrastructureGlobal != nil {
		infrastructureGlobal.Destroy()
	}
}

//LoadServerHardwareList  - загрузка информации со всех точек подключения по всем серверам
func (infra *OVInfrastructure) LoadServerHardwareList() ([]*ServerHardware, error) {
	for _, endpoint := range infra.endpoints {
		var ClientOV *ov.OVClient
		ovc := ClientOV.NewOVClient(
			endpoint.login,
			endpoint.password,
			endpoint.domain,
			endpoint.endpoint,
			false,
			1600,
			"*")
		filters := []string{""}
		sort := ""
		start := ""
		count := ""
		expand := "all"
		ServerList, err := ovc.GetServerHardwareList(filters, sort, start, count, expand)
		if err == nil {
			for _, rec := range ServerList.Members {
				s := ServerHardware{}
				s.Base = rec
				s.Memory, _ = GetServerHardwareMemory(ovc, rec.UUID)        //запрос по памяти в сервере
				s.Storage, _ = GetServerHardwareLocalStorage(ovc, rec.UUID) //запрос по локальным хранилищам
				infra.Servers = append(infra.Servers, &s)
				infra.ServersCount++

			}
		} else {
			fmt.Println("Failed to fetch server List : ", err)
		}
		for i := ServerList.Count; i < ServerList.Total; i = i + ServerList.Count {
			if ServerList.Total-i > ServerList.Count {
				count = strconv.Itoa(ServerList.Count)
			} else {
				count = strconv.Itoa(ServerList.Total - i)
			}
			start = strconv.Itoa(i)
			ServerList, err = ovc.GetServerHardwareList(filters, sort, start, count, expand)
			if err == nil {
				for _, rec := range ServerList.Members {
					s := ServerHardware{}
					s.Base = rec
					s.Memory, _ = GetServerHardwareMemory(ovc, rec.UUID)        //запрос по памяти в сервере
					s.Storage, _ = GetServerHardwareLocalStorage(ovc, rec.UUID) //запрос по локальным хранилищам
					infra.Servers = append(infra.Servers, &s)
					infra.ServersCount++
				}
			} else {
				fmt.Println("Failed to fetch server List : ", err)
			}
		}

	}
	return infra.Servers, nil
}

//FindServerHardwareSN  - поиск информации со всех точек подключения по серийному номеру
func (infra *OVInfrastructure) FindServerHardwareSN(sn string) (*ServerHardware, error) {
	for _, srvHW := range infra.Servers {
		if sn == string(srvHW.Base.SerialNumber) {
			return srvHW, nil
		}
	}
	return nil, errors.New("Serial Number not found")
}
