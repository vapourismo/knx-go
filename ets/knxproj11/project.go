// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxproj11

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"

	"fmt"

	"github.com/blang/semver"
)

// A ComObjectRefID is the ID of a communication object reference.
type ComObjectRefID string

// A ComObjectInstanceRef references a communication object via a reference to it. (Yes, you read
// that correctly.)
type ComObjectInstanceRef struct {
	RefID         ComObjectRefID
	DatapointType string
	Sends         []GroupAddressID
	Receives      []GroupAddressID
}

// UnmarshalXML extracts the ComObject information.
func (obj *ComObjectInstanceRef) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		RefID         string `xml:"RefID,attr"`
		DatapointType string `xml:"DatapointType,attr"`
		Connectors    struct {
			A []struct {
				XMLName xml.Name
				RefID   string `xml:"GroupAddressRefId,attr"`
			} `xml:",any"`
		}
	}

	// Decode element based on the layout above.
	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	obj.RefID = ComObjectRefID(doc.RefID)
	obj.DatapointType = doc.DatapointType
	obj.Sends = nil
	obj.Receives = nil

	// Fill in Send/Receive connections to communication object refs.
	for _, con := range doc.Connectors.A {
		switch con.XMLName.Local {
		case "Send":
			obj.Sends = append(obj.Sends, GroupAddressID(con.RefID))

		case "Receive":
			obj.Receives = append(obj.Receives, GroupAddressID(con.RefID))
		}
	}

	return nil
}

// A DeviceID is the ID of a device.
type DeviceID string

// Device is a KNX device.
type Device struct {
	ID         DeviceID               `xml:"Id,attr"`
	Name       string                 `xml:"Name,attr"`
	Address    uint                   `xml:"Address,attr"`
	ComObjects []ComObjectInstanceRef `xml:"ComObjectInstanceRefs>ComObjectInstanceRef"`
}

// A LineID is the ID of a line.
type LineID string

// Line is a KNX line.
type Line struct {
	ID      LineID   `xml:"Id,attr"`
	Name    string   `xml:"Name,attr"`
	Address uint     `xml:"Address,attr"`
	Devices []Device `xml:"DeviceInstance"`
}

// An AreaID is the ID of a line.
type AreaID string

// Area is a KNX area.
type Area struct {
	ID      AreaID `xml:"Id,attr"`
	Name    string `xml:"Name,attr"`
	Address uint   `xml:"Address,attr"`
	Lines   []Line `xml:"Line"`
}

// A GroupAddressID is the ID of a group address.
type GroupAddressID string

// GroupAddress is a group address.
type GroupAddress struct {
	ID      GroupAddressID `xml:"Id,attr"`
	Name    string         `xml:"Name,attr"`
	Address uint16         `xml:"Address,attr"`
}

// GroupAddressRange is a range of group addresses and sub-ranges.
type GroupAddressRange struct {
	Name       string              `xml:"Name,attr"`
	RangeStart uint16              `xml:"RangeStart,attr"`
	RangeEnd   uint16              `xml:"RangeEnd,attr"`
	Addresses  []GroupAddress      `xml:"GroupAddress"`
	SubRanges  []GroupAddressRange `xml:"GroupRange"`
}

// Installation is an installation.
type Installation struct {
	Name           string              `xml:"Name,attr"`
	Topology       []Area              `xml:"Topology>Area"`
	GroupAddresses []GroupAddressRange `xml:"GroupAddresses>GroupRanges>GroupRange"`
}

// Tool describes a tool.
type Tool struct {
	Name    string
	Version semver.Version
}

// Project is an ETS project.
type Project struct {
	ID            string         `xml:"Id,attr"`
	CreatedBy     Tool           `xml:"-"`
	Installations []Installation `xml:"Installations>Installation"`
}

// xmlns is the expected namespace for version 1.1 of the project spec.
const xmlns = "http://knx.org/xml/project/11"

// ErrInvalidNamespace indicates that given project file is scoped in an unknown namespace.
// The project that will be returned alongside this warning might still be usable.
var ErrInvalidNamespace = errors.New("Namespace does not match version 1.1 namespace")

func processProject(r io.Reader) (_ Project, err error) {
	var doc struct {
		Namespace   string `xml:"xmlns,attr"`
		ToolVersion string `xml:"ToolVersion,attr"`
		Project     Project
	}

	if err = xml.NewDecoder(r).Decode(&doc); err != nil {
		return
	}

	fmt.Sscanf(
		doc.ToolVersion,
		"%s %d.%d.%d",
		&doc.Project.CreatedBy.Name,
		&doc.Project.CreatedBy.Version.Major,
		&doc.Project.CreatedBy.Version.Minor,
		&doc.Project.CreatedBy.Version.Patch,
	)

	if doc.Namespace != xmlns {
		err = ErrInvalidNamespace
	}

	// We return the project even if namespace does not match. The project might be potentially
	// usable.
	return doc.Project, err
}

func processZippedProject(file *zip.File) (proj Project, err error) {
	r, err := file.Open()
	if err != nil {
		return
	}

	proj, err = processProject(r)
	r.Close()

	return
}
