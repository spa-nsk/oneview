package oneview

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/rest"
	"github.com/HewlettPackard/oneview-golang/utils"
)

//MemoryLocation  - расположение можуля памяти Socket - номер CPU, Slot - номер слота для данного CPU
type MemoryLocation struct {
	Slot   int `json:"Slot"`
	Socket int `json:"Socket"`
}

//HPE  - структура показывающая состояние модуля памяти DIMMStatus в нормальном состоянии в значении "GoodInUse" MinimumVoltageVoltsX10 напряжение питания в вольтах
type HPE struct {
	DIMMStatus             string `json:"DIMMStatus"`
	MinimumVoltageVoltsX10 int    `json:"MinimumVoltageVoltsX10"`
}

//OemMemory  - данные по производителям, в данном случае HPE
type OemMemory struct {
	Hpe HPE `json:"Hpe"`
}

//MemoryStatus  - состояние модуля памяти Health в значении "OK",  State в значении "Enabled"
type MemoryStatus struct {
	Health string `json:"Health"`
	State  string `json:"State"`
}

//MemoryModule  - описание модуля памяти
type MemoryModule struct {
	BaseModuleType    string         `json:"BaseModuleType"`    //"RDIMM","UDIMM","SO_DIMM","LRDIMM","Mini_RDIMM","Mini_UDIMM","SO_RDIMM_72b","SO_UDIMM_72b","SO_DIMM_16b","SO_DIMM_32b"
	CapacityMiB       int            `json:"CapacityMiB"`       //емкость модуля в мегабайтах
	DeviceLocator     string         `json:"DeviceLocator"`     // "PROC1 DIMM 1" дублирует MemoryLocation в виде строки
	Manufacturer      string         `json:"Manufacturer"`      //производитель модуля памяти "HPE"
	MemoryDeviceType  string         `json:"MemoryDeviceType"`  //тип памяти "DDR","DDR2","DDR3","DDR4","DDR4_SDRAM","DDR4E_SDRAM","LPDDR4_SDRAM","DDR3_SDRAM","LPDDR3_SDRAM","DDR2_SDRAM","DDR2_SDRAM_FB_DIMM","DDR2_SDRAM_FB_DIMM_PROBE","DDR_SGRAM","DDR_SDRAM","ROM","SDRAM","EDO","FastPageMode","PipelinedNibble"
	MemoryLocation    MemoryLocation `json:"MemoryLocation"`    //расположение модуля памяти
	Name              string         `json:"Name"`              // "proc1dimm1" наименование модуля
	Oem               OemMemory      `json:"Oem"`               //дополнения от производителя
	OperatingSpeedMhz int            `json:"OperatingSpeedMhz"` //частота работы модуля памяти в мегагерцах
	PartNumber        string         `json:"PartNumber"`        //номер производителя
	RankCount         int            `json:"RankCount"`         //ранг модуля памяти
	Status            MemoryStatus   `json:"Status"`            //состояние модуля памяти
}

//ServerHardwareMemory  - структура возвращаемая при запросе, содержит количество и описание модулей памяти сервера.
type ServerHardwareMemory struct {
	URI             string         `json:"uri"`             //uri запроса
	CollectionState string         `json:"collectionState"` //состояние "Collected"
	Modified        time.Time      `json:"modified"`
	Data            []MemoryModule `json:"data"` //срез с данными модулей данные только по занятым слотам!!!
	Etag            string         `json:"etag"`
	Count           int            `json:"count"` //количество модулей
	Name            string         `json:"name"`  //"Memory"
}

// GetServerHardwareMemory gets a server hardware with uri
func GetServerHardwareMemory(c *ov.OVClient, uuid utils.Nstring) (ServerHardwareMemory, error) {

	//	var hardware ServerHardware
	var (
		serverHardwareMemory ServerHardwareMemory
		dataMem              []MemoryModule
	)
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, "/rest/server-hardware/"+uuid.String()+"/memory", nil)
	if err != nil {
		return serverHardwareMemory, err
	}

	if err := json.Unmarshal([]byte(data), &serverHardwareMemory); err != nil {
		return serverHardwareMemory, err
	}
	dataMem = make([]MemoryModule, 0)
	for _, rec := range serverHardwareMemory.Data {
		if rec.CapacityMiB > 0 {
			mem := MemoryModule{}
			mem = rec
			dataMem = append(dataMem, mem)
		}
	}
	serverHardwareMemory.Data = dataMem
	serverHardwareMemory.Count = len(dataMem)
	return serverHardwareMemory, nil
}

type EnvironmentalConfiguration struct {
	PowerHistorySupported        bool   `json:"powerHistorySupported"`
	ThermalHistorySupported      bool   `json:"thermalHistorySupported"`
	UtilizationHistorySupported  bool   `json:"utilizationHistorySupported"`
	CapHistorySupported          bool   `json:"capHistorySupported"`
	HistorySampleIntervalSeconds int    `json:"historySampleIntervalSeconds"`
	HistoryBufferSize            int    `json:"historyBufferSize"`
	PowerCapType                 string `json:"powerCapType"`
	IdleMaxPower                 int    `json:"idleMaxPower"`
	CalibratedMaxPower           int    `json:"calibratedMaxPower"`
	PsuList                      []struct {
		Capacity     int    `json:"capacity"`
		PsuID        int    `json:"psuId"`
		InputVoltage int    `json:"inputVoltage"`
		Side         string `json:"side"`
		PowerType    string `json:"powerType"`
	} `json:"psuList"`
	Height             int         `json:"height"`
	RackID             interface{} `json:"rackId"`
	USlot              int         `json:"uSlot"`
	RackName           interface{} `json:"rackName"`
	RelativeOrder      int         `json:"relativeOrder"`
	RackModel          string      `json:"rackModel"`
	LicenseRequirement string      `json:"licenseRequirement"`
	RackUHeight        int         `json:"rackUHeight"`
}

// GetServerHardwareMemory gets a server hardware with uri
func GetServerEnvConfig(c *ov.OVClient, uuid utils.Nstring) (EnvironmentalConfiguration, error) {
	var envConf EnvironmentalConfiguration

	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, "/rest/server-hardware/"+uuid.String()+"/environmentalConfiguration", nil)
	if err != nil {
		return envConf, err
	}

	if err := json.Unmarshal([]byte(data), &envConf); err != nil {
		return envConf, err
	}

	return envConf, nil
}

//CacheModuleStatus  - Health состояние модуля кэша "OK" при наличие и нормально состоянии, если отсутствует то ""
type CacheModuleStatus struct {
	Health string `json:"Health"`
}

//ControllerSatus состояние контроллера Health "OK", State "Enabled"
type ControllerSatus struct {
	Health string `json:"Health"`
	State  string `json:"State"`
}

//ControllerBoard  -
type ControllerBoard struct {
	Status ControllerSatus `json:"Status"`
}

//Version  - версия в виде строки VersionString
type Version struct {
	VersionString string `json:"VersionString"`
}

//FirmwareVersion  - текущая версия Firmware
type FirmwareVersion struct {
	Current Version `json:"Current"`
}

//Status  - состояние Health "OK", State "Enabled"
type Status struct {
	Health string `json:"Health"`
	State  string `json:"State"`
}

//LogicalDataDrive  - описание диска в логическом томе
type LogicalDataDrive struct {
	CapacityMiB            float64         `json:"CapacityMiB"`            //емкость диска в мегабайтах
	DiskDriveStatusReasons []string        `json:"DiskDriveStatusReasons"` //стаатусы "None"
	DiskDriveUse           string          `json:"DiskDriveUse"`           //использование диска в томе "Data"
	EncryptedDrive         bool            `json:"EncryptedDrive"`         //включено ли шифрование true, false
	FirmwareVersion        FirmwareVersion `json:"FirmwareVersion"`        //текущая версия прошивки
	Location               string          `json:"Location"`               //расположение "1I:1:7"
	LocationFormat         string          `json:"LocationFormat"`         //формат строкти расположения "ControllerPort:Box:Bay"
	MediaType              string          `json:"MediaType"`              //тип диска "HDD", ...
	Model                  string          `json:"Model"`                  //модель "MB4000JVYZQ"
	SerialNumber           string          `json:"SerialNumber"`           //серийный номер "ZC17LE2X"
	Status                 Status          `json:"Status"`                 //состояние
}

//LocalLogicalDrive  - описание тома
type LocalLogicalDrive struct {
	CapacityMiB            float64            `json:"CapacityMiB"`            //емкость в мегабайтах
	DataDrives             []LogicalDataDrive `json:"DataDrives"`             //диски
	LogicalDriveEncryption bool               `json:"LogicalDriveEncryption"` //шифрование
	LogicalDriveNumber     int                `json:"LogicalDriveNumber"`     //номер логического диска
	LogicalDriveType       string             `json:"LogicalDriveType"`       //тип логического диска
	Raid                   string             `json:"Raid"`                   //уровень Raid "5", ...
	Status                 Status             `json:"Status"`                 //состояние
}

//LocalPhysicalDrive  - описание диска
type LocalPhysicalDrive struct {
	CapacityMiB            float64         `json:"CapacityMiB"`            //емкость
	DiskDriveStatusReasons []string        `json:"DiskDriveStatusReasons"` //стаатусы "None"
	DiskDriveUse           string          `json:"DiskDriveUse"`           //использование диска в томе "Data"
	EncryptedDrive         bool            `json:"EncryptedDrive"`         //включено ли шифрование true, false
	FirmwareVersion        FirmwareVersion `json:"FirmwareVersion"`        //текущая версия прошивки
	Location               string          `json:"Location"`               //расположение "1I:1:7"
	LocationFormat         string          `json:"LocationFormat"`         //формат строкти расположения "ControllerPort:Box:Bay"
	MediaType              string          `json:"MediaType"`              //тип диска "HDD", ...
	Model                  string          `json:"Model"`                  //модель "MB4000JVYZQ"
	SerialNumber           string          `json:"SerialNumber"`           //серийный номер "ZC17LE2X"
	Status                 Status          `json:"Status"`                 //состояние
}

//LocalStorageEnclosure  - описание корзины для дисков
type LocalStorageEnclosure struct {
	DriveBayCount   int             `json:"DriveBayCount"`   //количество мест
	FirmwareVersion FirmwareVersion `json:"FirmwareVersion"` //версия прошивки
	Location        string          `json:"Location"`        //расположение "1I:1"
	LocationFormat  string          `json:"LocationFormat"`  //формат расположения "ControllerPort:Box"
	Status          Status          `json:"Status"`          //состояние
}

//LocalStorage  - описание локального хранилища
type LocalStorage struct {
	AdapterType              string                  `json:"AdapterType"`             //тип адаптера "SmartArray"
	CacheMemorySizeMiB       int                     `json:"CacheMemorySizeMiB"`      //размер кэша контроллера в мегабайтах
	CacheModuleSerialNumber  string                  `json:"CacheModuleSerialNumber"` //серийный номер контроллера
	CacheModuleStatus        CacheModuleStatus       `json:"CacheModuleStatus"`       //состояние кэша
	ControllerBoard          ControllerBoard         `json:"ControllerBoard"`
	EncryptionCspTestPassed  bool                    `json:"EncryptionCspTestPassed"`  //тезультат тестового прохода
	EncryptionEnabled        bool                    `json:"EncryptionEnabled"`        //шифрование включено?
	EncryptionSelfTestPassed bool                    `json:"EncryptionSelfTestPassed"` //результата теста шифрования
	FirmwareVersion          FirmwareVersion         `json:"FirmwareVersion"`          //прошивка
	Location                 string                  `json:"Location"`                 //расположение "Slot 3"
	LogicalDrives            []LocalLogicalDrive     `json:"LogicalDrives"`            //томы (логические диски)
	Model                    string                  `json:"Model"`                    //модель "Smart array p840 Controller"
	Name                     string                  `json:"Name"`                     //наименование "Controller in Slot 3"
	PhysicalDrives           []LocalPhysicalDrive    `json:"PhysicalDrives"`           //диски подключенные к хранилищу
	SerialNumber             string                  `json:"SerialNumber"`             //серийный номер
	Status                   Status                  `json:"Status"`                   //состоние
	StorageEnclosures        []LocalStorageEnclosure `json:"StorageEnclosures"`        //корзины для дисков
}

//ServerHardwareLocalStorage  - структура возвращаемая по запросу о хранилищах
type ServerHardwareLocalStorage struct {
	URI             string         `json:"uri"`             //ссылка на запрос
	CollectionState string         `json:"collectionState"` //состояние "Collected"
	Modified        time.Time      `json:"modified"`        //время модификации
	Data            []LocalStorage `json:"data"`            //срез хранилищ
	Etag            string         `json:"etag"`            //тэг
	Count           int            `json:"count"`           //количество
	Name            string         `json:"name"`            //"LocalStorage"
}

// GetServerHardwareLocalStorage  - запрос на получение информации по локальным хранилищам сервера по uuid
func GetServerHardwareLocalStorage(c *ov.OVClient, uuid utils.Nstring) (ServerHardwareLocalStorage, error) {

	//	var hardware ServerHardware
	var (
		serverHardwareLocalStorage ServerHardwareLocalStorage
	)
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, "/rest/server-hardware/"+uuid.String()+"/localStorage", nil)
	if err != nil {
		return serverHardwareLocalStorage, err
	}

	if err := json.Unmarshal([]byte(data), &serverHardwareLocalStorage); err != nil {
		return serverHardwareLocalStorage, err
	}
	return serverHardwareLocalStorage, nil
}

type Enclosure struct {
	Type                 string      `json:"type"`
	URI                  string      `json:"uri"`
	Category             string      `json:"category"`
	ETag                 time.Time   `json:"eTag"`
	Created              time.Time   `json:"created"`
	Modified             time.Time   `json:"modified"`
	RefreshState         string      `json:"refreshState"`
	StateReason          string      `json:"stateReason"`
	EnclosureType        string      `json:"enclosureType"`
	EnclosureTypeURI     string      `json:"enclosureTypeUri"`
	EnclosureModel       string      `json:"enclosureModel"`
	UUID                 string      `json:"uuid"`
	SerialNumber         string      `json:"serialNumber"`
	PartNumber           string      `json:"partNumber"`
	ReconfigurationState string      `json:"reconfigurationState"`
	UIDState             interface{} `json:"uidState"`
	LicensingIntent      string      `json:"licensingIntent"`
	DeviceBayCount       int         `json:"deviceBayCount"`
	DeviceBays           []struct {
		Type                                    string      `json:"type"`
		BayNumber                               int         `json:"bayNumber"`
		Model                                   interface{} `json:"model"`
		DevicePresence                          string      `json:"devicePresence"`
		ProfileURI                              interface{} `json:"profileUri"`
		DeviceURI                               string      `json:"deviceUri"`
		CoveredByProfile                        interface{} `json:"coveredByProfile"`
		CoveredByDevice                         string      `json:"coveredByDevice"`
		Ipv4Setting                             interface{} `json:"ipv4Setting"`
		Ipv6Setting                             interface{} `json:"ipv6Setting"`
		URI                                     string      `json:"uri"`
		Category                                string      `json:"category"`
		ETag                                    interface{} `json:"eTag"`
		Created                                 interface{} `json:"created"`
		Modified                                interface{} `json:"modified"`
		AvailableForHalfHeightProfile           bool        `json:"availableForHalfHeightProfile"`
		AvailableForFullHeightProfile           bool        `json:"availableForFullHeightProfile"`
		DeviceBayType                           string      `json:"deviceBayType"`
		DeviceFormFactor                        string      `json:"deviceFormFactor"`
		BayPowerState                           string      `json:"bayPowerState"`
		ChangeState                             string      `json:"changeState"`
		AvailableForHalfHeightDoubleWideProfile bool        `json:"availableForHalfHeightDoubleWideProfile"`
		AvailableForFullHeightDoubleWideProfile bool        `json:"availableForFullHeightDoubleWideProfile"`
		UUID                                    interface{} `json:"uuid"`
	} `json:"deviceBays"`
	InterconnectBayCount int `json:"interconnectBayCount"`
	InterconnectBays     []struct {
		BayNumber              int         `json:"bayNumber"`
		InterconnectURI        string      `json:"interconnectUri"`
		LogicalInterconnectURI interface{} `json:"logicalInterconnectUri"`
		InterconnectModel      string      `json:"interconnectModel"`
		Ipv4Setting            interface{} `json:"ipv4Setting"`
		Ipv6Setting            interface{} `json:"ipv6Setting"`
		SerialNumber           string      `json:"serialNumber"`
		InterconnectBayType    string      `json:"interconnectBayType"`
		ChangeState            string      `json:"changeState"`
		BayPowerState          string      `json:"bayPowerState"`
	} `json:"interconnectBays"`
	FanBayCount int `json:"fanBayCount"`
	FanBays     []struct {
		BayNumber       int    `json:"bayNumber"`
		DevicePresence  string `json:"devicePresence"`
		DeviceRequired  bool   `json:"deviceRequired"`
		Status          string `json:"status"`
		Model           string `json:"model"`
		PartNumber      string `json:"partNumber"`
		SparePartNumber string `json:"sparePartNumber"`
		FanBayType      string `json:"fanBayType"`
		ChangeState     string `json:"changeState"`
		State           string `json:"state"`
	} `json:"fanBays"`
	PowerSupplyBayCount int `json:"powerSupplyBayCount"`
	PowerSupplyBays     []struct {
		BayNumber          int    `json:"bayNumber"`
		DevicePresence     string `json:"devicePresence"`
		Status             string `json:"status"`
		Model              string `json:"model"`
		SerialNumber       string `json:"serialNumber"`
		PartNumber         string `json:"partNumber"`
		SparePartNumber    string `json:"sparePartNumber"`
		PowerSupplyBayType string `json:"powerSupplyBayType"`
		ChangeState        string `json:"changeState"`
	} `json:"powerSupplyBays"`
	EnclosureGroupURI    interface{} `json:"enclosureGroupUri"`
	FwBaselineURI        interface{} `json:"fwBaselineUri"`
	FwBaselineName       interface{} `json:"fwBaselineName"`
	IsFwManaged          bool        `json:"isFwManaged"`
	ForceInstallFirmware bool        `json:"forceInstallFirmware"`
	LogicalEnclosureURI  interface{} `json:"logicalEnclosureUri"`
	ManagerBays          []struct {
		BayNumber      int           `json:"bayNumber"`
		ManagerType    string        `json:"managerType"`
		UIDState       interface{}   `json:"uidState"`
		BayPowerState  string        `json:"bayPowerState"`
		FwVersion      string        `json:"fwVersion"`
		DevicePresence string        `json:"devicePresence"`
		Role           string        `json:"role"`
		IPAddress      string        `json:"ipAddress"`
		ChangeState    string        `json:"changeState"`
		FwBuildDate    string        `json:"fwBuildDate"`
		DhcpIpv6Enable bool          `json:"dhcpIpv6Enable"`
		Ipv6Addresses  []interface{} `json:"ipv6Addresses"`
		FqdnHostName   string        `json:"fqdnHostName"`
		DhcpEnable     bool          `json:"dhcpEnable"`
	} `json:"managerBays"`
	SupportState               string `json:"supportState"`
	SupportDataCollectionState string `json:"supportDataCollectionState"`
	SupportDataCollectionType  string `json:"supportDataCollectionType"`
	SupportDataCollectionsURI  string `json:"supportDataCollectionsUri"`
	RemoteSupportURI           string `json:"remoteSupportUri"`
	RemoteSupportSettings      struct {
		RemoteSupportCurrentState string `json:"remoteSupportCurrentState"`
		Destination               string `json:"destination"`
	} `json:"remoteSupportSettings"`
	CrossBars            []interface{} `json:"crossBars"`
	Partitions           []interface{} `json:"partitions"`
	ScopesURI            string        `json:"scopesUri"`
	Status               string        `json:"status"`
	Name                 string        `json:"name"`
	State                string        `json:"state"`
	Description          interface{}   `json:"description"`
	StandbyOaPreferredIP string        `json:"standbyOaPreferredIP"`
	ActiveOaPreferredIP  string        `json:"activeOaPreferredIP"`
	AssetTag             string        `json:"assetTag"`
	RackName             string        `json:"rackName"`
	VcmDomainID          string        `json:"vcmDomainId"`
	VcmDomainName        string        `json:"vcmDomainName"`
	VcmMode              bool          `json:"vcmMode"`
	VcmURL               string        `json:"vcmUrl"`
	OaBays               int           `json:"oaBays"`
	MigrationState       string        `json:"migrationState"`
}

func GetServerEnclosure(c *ov.OVClient, encuri utils.Nstring) (Enclosure, error) {

	var (
		encHardware Enclosure
	)

	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, encuri.String(), nil)
	if err != nil {
		return encHardware, err
	}

	if err := json.Unmarshal([]byte(data), &encHardware); err != nil {
		return encHardware, err
	}
	return encHardware, nil
}

type ServerSSOUrl struct {
	IloSsoURL string `json:"iloSsoUrl"`
}

func GetServerILOssoUrl(c *ov.OVClient, uuid utils.Nstring) (string, error) {

	var ssoILOUrl ServerSSOUrl
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, "/rest/server-hardware/"+uuid.String()+"/iloSsoUrl", nil)
	if err != nil {
		return ssoILOUrl.IloSsoURL, err
	}

	fmt.Println(string(data))

	if err := json.Unmarshal([]byte(data), &ssoILOUrl); err != nil {
		return ssoILOUrl.IloSsoURL, err
	}

	return ssoILOUrl.IloSsoURL, nil
}

type ServerjavaRemoteConsoleUrl struct {
	JavaRemoteConsoleUrl string `json:"javaRemoteConsoleUrl"`
}

func GetServerjavaRemoteConsoleUrl(c *ov.OVClient, uuid utils.Nstring) (string, error) {

	var javaILOUrl ServerjavaRemoteConsoleUrl
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	// rest call
	data, err := c.RestAPICall(rest.GET, "/rest/server-hardware/"+uuid.String()+"/javaRemoteConsoleUrl", nil)
	if err != nil {
		return javaILOUrl.JavaRemoteConsoleUrl, err
	}

	fmt.Println(string(data))

	if err := json.Unmarshal([]byte(data), &javaILOUrl); err != nil {
		return javaILOUrl.JavaRemoteConsoleUrl, err
	}

	return javaILOUrl.JavaRemoteConsoleUrl, nil
}

type DatacentersList struct {
}

// LoadDatcenterList
func (infra *OVInfrastructure) LoadDatcenterList(c *ov.OVClient, filters []string, sort string, start string, count string) (DatacentersList, error) {
	var (
		uri         = "/rest/datacenters"
		q           map[string]interface{}
		datacenters DatacentersList
	)
	q = make(map[string]interface{})

	if len(filters) > 0 {
		q["filter"] = filters
	}

	if sort != "" {
		q["sort"] = sort
	}

	if start != "" {
		q["start"] = start
	}

	if count != "" {
		q["count"] = count
	}
	// refresh login
	c.RefreshLogin()
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())

	data, err := c.RestAPICall(rest.GET, uri, nil, q)
	if err != nil {
		return datacenters, err
	}

	fmt.Println(string(data))
	/*
		if err := json.Unmarshal([]byte(data), &datacenters); err != nil {
			return datacenters, err
		}
	*/
	/*
		if count == "" {
			total := strconv.Itoa(serverlist.Total)
			return c.GetServerHardwareList(filters, "", "", total, "")
		}

		for i := 0; i < serverlist.Count; i++ {
			serverlist.Members[i].Client = c
		}
	*/
	return datacenters, nil

}
