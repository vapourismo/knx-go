// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxproj11

import (
	"archive/zip"
	"encoding/xml"
	"io/ioutil"

	"github.com/vapourismo/knx-go/knx/cemi"
)

// GroupObject is a group object.
type GroupObject struct {
	ID      string         `xml:"Id,attr"`
	Name    string         `xml:"Name,attr"`
	Address cemi.GroupAddr `xml:"Address,attr"`
}

// GroupObjectRange is a range of group objects and ranges.
type GroupObjectRange struct {
	ID        string             `xml:"Id,attr"`
	Name      string             `xml:"Name,attr"`
	Objects   []GroupObject      `xml:"GroupAddress"`
	SubRanges []GroupObjectRange `xml:"GroupRange"`
}

// ComObjectRef is a reference to a communcation object.
type ComObjectRef struct {
	RefID      string
	SendIDs    []string
	ReceiveIDs []string
}

// UnmarshalXML extracts the ComObject information.
func (obj *ComObjectRef) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		RefID      string `xml:"RefId,attr"`
		Connectors struct {
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

	obj.RefID = doc.RefID
	obj.SendIDs = nil
	obj.ReceiveIDs = nil

	// Fill in Send/Receive connections to communication object refs.
	for _, con := range doc.Connectors.A {
		switch con.XMLName.Local {
		case "Send":
			obj.SendIDs = append(obj.SendIDs, con.RefID)

		case "Receive":
			obj.ReceiveIDs = append(obj.ReceiveIDs, con.RefID)
		}
	}

	return nil
}

// Device is a KNX device.
type Device struct {
	ID         string         `xml:"Id,attr"`
	Name       string         `xml:"Name,attr"`
	Address    uint           `xml:"Address,attr"`
	ComObjects []ComObjectRef `xml:"ComObjectInstanceRefs>ComObjectInstanceRef"`
}

// Line is a KNX line.
type Line struct {
	ID      string   `xml:"Id,attr"`
	Name    string   `xml:"Name,attr"`
	Address uint     `xml:"Address,attr"`
	Devices []Device `xml:"DeviceInstance"`
}

// Area is a KNX area.
type Area struct {
	ID      string `xml:"Id,attr"`
	Name    string `xml:"Name,attr"`
	Address uint   `xml:"Address,attr"`
	Lines   []Line `xml:"Line"`
}

// Installation is an installation.
type Installation struct {
	ID     string             `xml:"InstallationId,attr"`
	Name   string             `xml:"Name,attr"`
	Areas  []Area             `xml:"Topology>Area"`
	Groups []GroupObjectRange `xml:"GroupAddresses>GroupRanges>GroupRange"`
}

// Project is an ETS project.
type Project struct {
	ID            string         `xml:"Id,attr"`
	Installations []Installation `xml:"Installations>Installation"`
}

// processProject
func processProject(file *zip.File) (proj Project, err error) {
	// Acquire a handle to the file.
	r, err := file.Open()
	if err != nil {
		return
	}

	// Retrieve the contents of the file.
	contents, err := ioutil.ReadAll(r)
	r.Close()

	if err != nil {
		return
	}

	var doc struct {
		Project Project
	}

	// Extract the information.
	if err = xml.Unmarshal(contents, &doc); err != nil {
		return
	}

	return doc.Project, nil
}
