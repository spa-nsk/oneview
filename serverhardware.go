package oneview

import (
	"encoding/json"
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
