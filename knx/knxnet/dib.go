// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"errors"
	"net"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

const (
	friendlyNameMaxLen = 30
)

// DescriptionType describes the type of a DeviceInformationBlock
type DescriptionType uint8

const (
	// DescriptionTypeDeviceInfo describes Device information e.g. KNX medium.
	DescriptionTypeDeviceInfo DescriptionType = 0x01

	// DescriptionTypeSupportedServiceFamilies describes Service families supported by the device.
	DescriptionTypeSupportedServiceFamilies DescriptionType = 0x02

	// DescriptionTypeIPConfig describes IP configuration
	DescriptionTypeIPConfig DescriptionType = 0x03

	// DescriptionTypeIPCurrentConfig describes current IP configuration
	DescriptionTypeIPCurrentConfig DescriptionType = 0x04

	// DescriptionTypeKNXAddresses describes KNX addresses
	DescriptionTypeKNXAddresses DescriptionType = 0x05

	// DescriptionTypeManufacturerData describes a DIB structure for further data defined by device manufacturer.
	DescriptionTypeManufacturerData DescriptionType = 0xfe
)

// KNXMedium describes the KNX medium type
type KNXMedium uint8

const (
	// KNXMediumTP1 is the TP1 medium
	KNXMediumTP1 KNXMedium = 0x02
	// KNXMediumPL110 is the PL110 medium
	KNXMediumPL110 KNXMedium = 0x04
	// KNXMediumRF is the RF medium
	KNXMediumRF KNXMedium = 0x10
	// KNXMediumIP is the IP medium
	KNXMediumIP KNXMedium = 0x20
)

// ProjectInstallationIdentifier describes a KNX project installation identifier
type ProjectInstallationIdentifier uint16

// DeviceStatus describes the device status
type DeviceStatus uint8

// DeviceSerialNumber desribes the serial number of a device
type DeviceSerialNumber [6]byte

// DeviceInformationBlock contains information about a device.
type DeviceInformationBlock struct {
	Type                    DescriptionType
	Medium                  KNXMedium
	Status                  DeviceStatus
	Source                  cemi.IndividualAddr
	ProjectIdentifier       ProjectInstallationIdentifier
	SerialNumber            DeviceSerialNumber
	RoutingMulticastAddress Address
	HardwareAddr            net.HardwareAddr
	FriendlyName            string
}

// Size returns the packed size.
func (DeviceInformationBlock) Size() uint {
	return 54
}

// Pack assembles the device information structure in the given buffer.
func (info *DeviceInformationBlock) Pack(buffer []byte) {
	buf := make([]byte, friendlyNameMaxLen)
	util.PackString(buf, friendlyNameMaxLen, info.FriendlyName)

	util.PackSome(
		buffer,
		uint8(info.Size()), uint8(info.Type),
		uint8(info.Medium), uint8(info.Status),
		uint16(info.Source),
		uint16(info.ProjectIdentifier),
		info.SerialNumber[:],
		info.RoutingMulticastAddress[:],
		[]byte(info.HardwareAddr),
		buf,
	)
}

// Unpack parses the given data in order to initialize the structure.
func (info *DeviceInformationBlock) Unpack(data []byte) (n uint, err error) {
	var length uint8

	info.HardwareAddr = make([]byte, 6)
	if n, err = util.UnpackSome(
		data,
		&length, (*uint8)(&info.Type),
		(*uint8)(&info.Medium), (*uint8)(&info.Status),
		(*uint16)(&info.Source),
		(*uint16)(&info.ProjectIdentifier),
		info.SerialNumber[:],
		info.RoutingMulticastAddress[:],
		[]byte(info.HardwareAddr),
	); err != nil {
		return
	}

	nn, err := util.UnpackString(data[n:], friendlyNameMaxLen, &info.FriendlyName)
	if err != nil {
		return n, err
	}
	n += nn

	if length != uint8(info.Size()) {
		return n, errors.New("Device info structure length is invalid")
	}

	return
}

// SupportedServicesDIB contains information about the supported services of a device
type SupportedServicesDIB struct {
	Type     DescriptionType
	Families []ServiceFamily
}

// Size returns the packed size.
func (sdib SupportedServicesDIB) Size() uint {
	size := uint(2)
	for _, f := range sdib.Families {
		size += f.Size()
	}

	return size
}

// Pack assembles the supported services structure in the given buffer.
func (sdib *SupportedServicesDIB) Pack(buffer []byte) {
	util.PackSome(
		buffer,
		uint8(sdib.Size()), uint8(sdib.Type),
	)

	offset := uint(2)
	for _, f := range sdib.Families {
		f.Pack(buffer[offset:])
		offset += f.Size()
	}
}

// Unpack parses the given data in order to initialize the structure.
func (sdib *SupportedServicesDIB) Unpack(data []byte) (n uint, err error) {
	var length uint8
	if n, err = util.UnpackSome(
		data,
		&length, (*uint8)(&sdib.Type),
	); err != nil {
		return
	}

	for n < uint(length) {
		f := ServiceFamily{}
		nn, err := f.Unpack(data[n:])
		if err != nil {
			return n, errors.New("Unable to unpack service family")
		}

		n += nn
		sdib.Families = append(sdib.Families, f)
	}

	if length != uint8(sdib.Size()) {
		return n, errors.New("Supported Services structure length is invalid")
	}

	return
}

// ServiceFamilyType describes a KNXnet service family type
type ServiceFamilyType uint8

const (
	// ServiceFamilyTypeIPCore is the KNXnet/IP Core family type
	ServiceFamilyTypeIPCore = 0x02
	// ServiceFamilyTypeIPDeviceManagement is the KNXnet/IP Device Management family type
	ServiceFamilyTypeIPDeviceManagement = 0x03
	// ServiceFamilyTypeIPTunnelling is the KNXnet/IP Tunnelling family type
	ServiceFamilyTypeIPTunnelling = 0x04
	// ServiceFamilyTypeIPRouting is the KNXnet/IP Routing family type
	ServiceFamilyTypeIPRouting = 0x05
	// ServiceFamilyTypeIPRemoteLogging is the KNXnet/IP Remote Logging family type
	ServiceFamilyTypeIPRemoteLogging = 0x06
	// ServiceFamilyTypeIPRemoteConfigurationAndDiagnosis is the KNXnet/IP Remote Configuration and Diagnosis family type
	ServiceFamilyTypeIPRemoteConfigurationAndDiagnosis = 0x07
	// ServiceFamilyTypeIPObjectServer is the KNXnet/IP Object Server family type
	ServiceFamilyTypeIPObjectServer = 0x08
)

// ServiceFamily describes a KNXnet service supported by a device
type ServiceFamily struct {
	Type    ServiceFamilyType
	Version uint8
}

// Size returns the packed size.
func (ServiceFamily) Size() uint {
	return 2
}

// Pack assembles the service family structure in the given buffer.
func (f *ServiceFamily) Pack(buffer []byte) {
	util.PackSome(
		buffer,
		uint8(f.Type), f.Version,
	)
}

// Unpack parses the given data in order to initialize the structure.
func (f *ServiceFamily) Unpack(data []byte) (n uint, err error) {
	return util.UnpackSome(data, (*uint8)(&f.Type), &f.Version)
}
