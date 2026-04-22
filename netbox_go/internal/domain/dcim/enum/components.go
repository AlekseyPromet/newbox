// Package enum содержит перечисления (enums) для компонентов DCIM
package enum

import "netbox_go/pkg/types"

// CableType представляет типы кабелей
type CableType string

const (
	CableTypeCat3        CableType = "cat3"
	CableTypeCat5        CableType = "cat5"
	CableTypeCat5e       CableType = "cat5e"
	CableTypeCat6        CableType = "cat6"
	CableTypeCat6a       CableType = "cat6a"
	CableTypeCat7        CableType = "cat7"
	CableTypeCat7a       CableType = "cat7a"
	CableTypeCat8        CableType = "cat8"
	CableTypeDACDirect   CableType = "dac-direct-attach"
	CableTypeFiberSingle CableType = "fiber-single-mode"
	CableTypeFiberMulti  CableType = "fiber-multi-mode"
	CableTypePower       CableType = "power"
	CableTypeCoaxial     CableType = "coaxial"
)

// GetAllCableTypes возвращает все возможные типы кабелей
func GetAllCableTypes() []CableType {
	return []CableType{
		CableTypeCat3,
		CableTypeCat5,
		CableTypeCat5e,
		CableTypeCat6,
		CableTypeCat6a,
		CableTypeCat7,
		CableTypeCat7a,
		CableTypeCat8,
		CableTypeDACDirect,
		CableTypeFiberSingle,
		CableTypeFiberMulti,
		CableTypePower,
		CableTypeCoaxial,
	}
}

// Validate проверяет корректность типа кабеля
func (c CableType) Validate() error {
	switch c {
	case CableTypeCat3, CableTypeCat5, CableTypeCat5e, CableTypeCat6,
		CableTypeCat6a, CableTypeCat7, CableTypeCat7a, CableTypeCat8,
		CableTypeDACDirect, CableTypeFiberSingle, CableTypeFiberMulti,
		CableTypePower, CableTypeCoaxial:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление типа кабеля
func (c CableType) String() string {
	return string(c)
}

// CableStatus представляет статусы кабеля
type CableStatus string

const (
	CableStatusConnected   CableStatus = "connected"
	CableStatusPlanned     CableStatus = "planned"
	CableStatusDecommissioning CableStatus = "decommissioning"
)

// GetAllCableStatuses возвращает все возможные статусы кабеля
func GetAllCableStatuses() []CableStatus {
	return []CableStatus{
		CableStatusConnected,
		CableStatusPlanned,
		CableStatusDecommissioning,
	}
}

// Validate проверяет корректность статуса кабеля
func (c CableStatus) Validate() error {
	switch c {
	case CableStatusConnected, CableStatusPlanned, CableStatusDecommissioning:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса кабеля
func (c CableStatus) String() string {
	return string(c)
}

// Color возвращает цвет статуса кабеля для UI
func (c CableStatus) Color() string {
	switch c {
	case CableStatusConnected:
		return "green"
	case CableStatusPlanned:
		return "blue"
	case CableStatusDecommissioning:
		return "yellow"
	default:
		return "gray"
	}
}

// ConsolePortType представляет типы консольных портов
type ConsolePortType string

const (
	ConsolePortTypeSerial     ConsolePortType = "serial"
	ConsolePortTypeUSB        ConsolePortType = "usb"
	ConsolePortTypeUSBMini    ConsolePortType = "usb-mini"
	ConsolePortTypeUSBMicro   ConsolePortType = "usb-micro"
	ConsolePortTypeUSBC       ConsolePortType = "usb-c"
)

// GetAllConsolePortTypes возвращает все возможные типы консольных портов
func GetAllConsolePortTypes() []ConsolePortType {
	return []ConsolePortType{
		ConsolePortTypeSerial,
		ConsolePortTypeUSB,
		ConsolePortTypeUSBMini,
		ConsolePortTypeUSBMicro,
		ConsolePortTypeUSBC,
	}
}

// Validate проверяет корректность типа консольного порта
func (c ConsolePortType) Validate() error {
	switch c {
	case ConsolePortTypeSerial, ConsolePortTypeUSB, ConsolePortTypeUSBMini,
		ConsolePortTypeUSBMicro, ConsolePortTypeUSBC:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление типа консольного порта
func (c ConsolePortType) String() string {
	return string(c)
}

// PowerPortType представляет типы портов питания
type PowerPortType string

const (
	PowerPortTypeIEC60320C13  PowerPortType = "iec-60320-c13"
	PowerPortTypeIEC60320C15  PowerPortType = "iec-60320-c15"
	PowerPortTypeIEC60320C19  PowerPortType = "iec-60320-c19"
	PowerPortTypeIEC60320C21  PowerPortType = "iec-60320-c21"
	PowerPortTypeNEMA515R     PowerPortType = "nema-5-15r"
	PowerPortTypeNEMA520R     PowerPortType = "nema-5-20r"
	PowerPortTypeNEMAL520R    PowerPortType = "nema-l5-20r"
	PowerPortTypeNEMAL530R    PowerPortType = "nema-l5-30r"
	PowerPortTypeNEMAL620R    PowerPortType = "nema-l6-20r"
	PowerPortTypeNEMAL630R    PowerPortType = "nema-l6-30r"
	PowerPortTypeCS6360C      PowerPortType = "cs6360c"
	PowerPortTypeCS8165C      PowerPortType = "cs8165c"
	PowerPortTypeCS8265C      PowerPortType = "cs8265c"
	PowerPortTypeCS8365C      PowerPortType = "cs8365c"
	PowerPortTypeCS8465C      PowerPortType = "cs8465c"
	PowerPortTypeITAE         PowerPortType = "ita-e"
	PowerPortTypeITAF         PowerPortType = "ita-f"
	PowerPortTypeITAG         PowerPortType = "ita-g"
	PowerPortTypeITAH         PowerPortType = "ita-h"
	PowerPortTypeITA          PowerPortType = "ita-i"
	PowerPortTypeITAJ         PowerPortType = "ita-j"
	PowerPortTypeITAK         PowerPortType = "ita-k"
	PowerPortTypeITAL         PowerPortType = "ita-l"
	PowerPortTypeITAM         PowerPortType = "ita-m"
	PowerPortTypeITAN         PowerPortType = "ita-n"
	PowerPortTypeITAO         PowerPortType = "ita-o"
	PowerPortTypeDC           PowerPortType = "dc-terminal"
	PowerPortTypeMolexMicrofit PowerPortType = "molex-micro-fit-1x2"
	PowerPortTypeMolexMicrofit2 PowerPortType = "molex-micro-fit-2x2"
	PowerPortTypeMolexMicrofit4 PowerPortType = "molex-micro-fit-2x4"
)

// GetAllPowerPortTypes возвращает все возможные типы портов питания
func GetAllPowerPortTypes() []PowerPortType {
	return []PowerPortType{
		PowerPortTypeIEC60320C13,
		PowerPortTypeIEC60320C15,
		PowerPortTypeIEC60320C19,
		PowerPortTypeIEC60320C21,
		PowerPortTypeNEMA515R,
		PowerPortTypeNEMA520R,
		PowerPortTypeNEMAL520R,
		PowerPortTypeNEMAL530R,
		PowerPortTypeNEMAL620R,
		PowerPortTypeNEMAL630R,
		PowerPortTypeCS6360C,
		PowerPortTypeCS8165C,
		PowerPortTypeCS8265C,
		PowerPortTypeCS8365C,
		PowerPortTypeCS8465C,
		PowerPortTypeITAE,
		PowerPortTypeITAF,
		PowerPortTypeITAG,
		PowerPortTypeITAH,
		PowerPortTypeITA,
		PowerPortTypeITAJ,
		PowerPortTypeITAK,
		PowerPortTypeITAL,
		PowerPortTypeITAM,
		PowerPortTypeITAN,
		PowerPortTypeITAO,
		PowerPortTypeDC,
		PowerPortTypeMolexMicrofit,
		PowerPortTypeMolexMicrofit2,
		PowerPortTypeMolexMicrofit4,
	}
}

// Validate проверяет корректность типа порта питания
func (p PowerPortType) Validate() error {
	switch p {
	case PowerPortTypeIEC60320C13, PowerPortTypeIEC60320C15, PowerPortTypeIEC60320C19,
		PowerPortTypeIEC60320C21, PowerPortTypeNEMA515R, PowerPortTypeNEMA520R,
		PowerPortTypeNEMAL520R, PowerPortTypeNEMAL530R, PowerPortTypeNEMAL620R,
		PowerPortTypeNEMAL630R, PowerPortTypeCS6360C, PowerPortTypeCS8165C,
		PowerPortTypeCS8265C, PowerPortTypeCS8365C, PowerPortTypeCS8465C,
		PowerPortTypeITAE, PowerPortTypeITAF, PowerPortTypeITAG, PowerPortTypeITAH,
		PowerPortTypeITA, PowerPortTypeITAJ, PowerPortTypeITAK, PowerPortTypeITAL,
		PowerPortTypeITAM, PowerPortTypeITAN, PowerPortTypeITAO, PowerPortTypeDC,
		PowerPortTypeMolexMicrofit, PowerPortTypeMolexMicrofit2, PowerPortTypeMolexMicrofit4:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление типа порта питания
func (p PowerPortType) String() string {
	return string(p)
}

// InterfaceType представляет типы сетевых интерфейсов
type InterfaceType string

const (
	InterfaceTypeVirtual      InterfaceType = "virtual"
	InterfaceTypeBridge       InterfaceType = "bridge"
	InterfaceTypeLAG          InterfaceType = "lag"
	InterfaceType100MEFixed   InterfaceType = "100base-tx"
	InterfaceType100MEPPFIXED InterfaceType = "100base-tx-pp"
	InterfaceType1GEFixed     InterfaceType = "1000base-t"
	InterfaceType1GEPPFIXED   InterfaceType = "1000base-t-pp"
	InterfaceType1GESFP       InterfaceType = "1000base-x-sfp"
	InterfaceType2GEXFP       InterfaceType = "2.5gbase-t"
	InterfaceType5GEXFP       InterfaceType = "5gbase-t"
	InterfaceType10GEXFP      InterfaceType = "10gbase-t"
	InterfaceType10GESFPPlus  InterfaceType = "10gbase-x-sfpp"
	InterfaceType10GESFP28    InterfaceType = "10gbase-x-sfp28"
	InterfaceType25GESFP28    InterfaceType = "25gbase-x-sfp28"
	InterfaceType40GEQSFPPlus InterfaceType = "40gbase-x-qsfpp"
	InterfaceType50GEPAM4     InterfaceType = "50gbase-x-pam4"
	InterfaceType100GESFP28   InterfaceType = "100gbase-x-sfp28"
	InterfaceType100GEQSFP28  InterfaceType = "100gbase-x-qsfp28"
	InterfaceType200GECFP     InterfaceType = "200gbase-x-cfp2"
	InterfaceType200GEQSFP56  InterfaceType = "200gbase-x-qsfp56"
	InterfaceType200GESFPDD   InterfaceType = "200gbase-x-sfpdd"
	InterfaceType400GEQSFPDD  InterfaceType = "400gbase-x-qsfpdd"
	InterfaceType400GEOSFP    InterfaceType = "400gbase-x-osfp"
	InterfaceType800GEQSFPDD  InterfaceType = "800gbase-x-qsfpdd"
	InterfaceType800GEOSFP    InterfaceType = "800gbase-x-osfp"
	InterfaceType16GEFC       InterfaceType = "16gfc-fibre"
	InterfaceType32GEFC       InterfaceType = "32gfc-fibre"
	InterfaceType64GEFC       InterfaceType = "64gfc-fibre"
	InterfaceType128GEFC      InterfaceType = "128gfc-fibre"
	InterfaceType256GEFC      InterfaceType = "256gfc-fibre"
	InterfaceType512GEFC      InterfaceType = "512gfc-fibre"
	InterfaceTypeIEEE80211A   InterfaceType = "ieee802.11a"
	InterfaceTypeIEEE80211G   InterfaceType = "ieee802.11g"
	InterfaceTypeIEEE80211N   InterfaceType = "ieee802.11n"
	InterfaceTypeIEEE80211AC  InterfaceType = "ieee802.11ac"
	InterfaceTypeIEEE80211AX  InterfaceType = "ieee802.11ax"
	InterfaceTypeIEEE80211AY  InterfaceType = "ieee802.11ay"
	InterfaceTypeIEEE80211BE  InterfaceType = "ieee802.11be"
	InterfaceTypeGSM          InterfaceType = "gsm"
	InterfaceTypeCDMA         InterfaceType = "cdma"
	InterfaceTypeLTE          InterfaceType = "lte"
	InterfaceTypeSONET        InterfaceType = "sonet"
	InterfaceTypeXDSL         InterfaceType = "xdsl"
	InterfaceTypeDOCSIS       InterfaceType = "docsis"
	InterfaceTypePON          InterfaceType = "pon"
	InterfaceTypeGPON         InterfaceType = "gpon"
	InterfaceTypeXGS          InterfaceType = "xgs-pon"
	InterfaceTypeNGPON2       InterfaceType = "ng-pon2"
	InterfaceTypeEPON         InterfaceType = "epon"
	InterfaceType10GEPON      InterfaceType = "10g-epon"
)

// GetAllInterfaceTypes возвращает все возможные типы интерфейсов
func GetAllInterfaceTypes() []InterfaceType {
	types := []InterfaceType{
		InterfaceTypeVirtual,
		InterfaceTypeBridge,
		InterfaceTypeLAG,
		InterfaceType100MEFixed,
		InterfaceType100MEPPFIXED,
		InterfaceType1GEFixed,
		InterfaceType1GEPPFIXED,
		InterfaceType1GESFP,
		InterfaceType2GEXFP,
		InterfaceType5GEXFP,
		InterfaceType10GEXFP,
		InterfaceType10GESFPPlus,
		InterfaceType10GESFP28,
		InterfaceType25GESFP28,
		InterfaceType40GEQSFPPlus,
		InterfaceType50GEPAM4,
		InterfaceType100GESFP28,
		InterfaceType100GEQSFP28,
		InterfaceType200GECFP,
		InterfaceType200GEQSFP56,
		InterfaceType200GESFPDD,
		InterfaceType400GEQSFPDD,
		InterfaceType400GEOSFP,
		InterfaceType800GEQSFPDD,
		InterfaceType800GEOSFP,
		InterfaceType16GEFC,
		InterfaceType32GEFC,
		InterfaceType64GEFC,
		InterfaceType128GEFC,
		InterfaceType256GEFC,
		InterfaceType512GEFC,
		InterfaceTypeIEEE80211A,
		InterfaceTypeIEEE80211G,
		InterfaceTypeIEEE80211N,
		InterfaceTypeIEEE80211AC,
		InterfaceTypeIEEE80211AX,
		InterfaceTypeIEEE80211AY,
		InterfaceTypeIEEE80211BE,
		InterfaceTypeGSM,
		InterfaceTypeCDMA,
		InterfaceTypeLTE,
		InterfaceTypeSONET,
		InterfaceTypeXDSL,
		InterfaceTypeDOCSIS,
		InterfaceTypePON,
		InterfaceTypeGPON,
		InterfaceTypeXGS,
		InterfaceTypeNGPON2,
		InterfaceTypeEPON,
		InterfaceType10GEPON,
	}
	return types
}

// Validate проверяет корректность типа интерфейса
func (i InterfaceType) Validate() error {
	validTypes := GetAllInterfaceTypes()
	for _, t := range validTypes {
		if i == t {
			return nil
		}
	}
	return types.ErrInvalidStatus
}

// String возвращает строковое представление типа интерфейса
func (i InterfaceType) String() string {
	return string(i)
}

// InterfaceMode представляет режимы работы VLAN на интерфейсе
type InterfaceMode string

const (
	InterfaceModeAccess  InterfaceMode = "access"
	InterfaceModeTagged  InterfaceMode = "tagged"
	InterfaceModeTaggedAll InterfaceMode = "tagged-all"
)

// GetAllInterfaceModes возвращает все возможные режимы интерфейса
func GetAllInterfaceModes() []InterfaceMode {
	return []InterfaceMode{
		InterfaceModeAccess,
		InterfaceModeTagged,
		InterfaceModeTaggedAll,
	}
}

// Validate проверяет корректность режима интерфейса
func (m InterfaceMode) Validate() error {
	switch m {
	case InterfaceModeAccess, InterfaceModeTagged, InterfaceModeTaggedAll:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление режима интерфейса
func (m InterfaceMode) String() string {
	return string(m)
}

// LinkStatus представляет статус линка интерфейса
type LinkStatus string

const (
	LinkStatusUp          LinkStatus = "up"
	LinkStatusDown        LinkStatus = "down"
	LinkStatusTesting     LinkStatus = "testing"
	LinkStatusUnknown     LinkStatus = "unknown"
	LinkStatusDormant     LinkStatus = "dormant"
	LinkStatusNotPresent  LinkStatus = "not-present"
	LinkStatusLowerLayerDown LinkStatus = "lower-layer-down"
)

// GetAllLinkStatuses возвращает все возможные статусы линка
func GetAllLinkStatuses() []LinkStatus {
	return []LinkStatus{
		LinkStatusUp,
		LinkStatusDown,
		LinkStatusTesting,
		LinkStatusUnknown,
		LinkStatusDormant,
		LinkStatusNotPresent,
		LinkStatusLowerLayerDown,
	}
}

// Validate проверяет корректность статуса линка
func (l LinkStatus) Validate() error {
	switch l {
	case LinkStatusUp, LinkStatusDown, LinkStatusTesting, LinkStatusUnknown,
		LinkStatusDormant, LinkStatusNotPresent, LinkStatusLowerLayerDown:
		return nil
	default:
		return types.ErrInvalidStatus
	}
}

// String возвращает строковое представление статуса линка
func (l LinkStatus) String() string {
	return string(l)
}

// Color возвращает цвет статуса линка для UI
func (l LinkStatus) Color() string {
	switch l {
	case LinkStatusUp:
		return "green"
	case LinkStatusDown:
		return "red"
	case LinkStatusTesting:
		return "yellow"
	case LinkStatusUnknown:
		return "gray"
	case LinkStatusDormant:
		return "blue"
	default:
		return "gray"
	}
}
