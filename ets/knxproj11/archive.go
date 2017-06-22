// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxproj11

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"path"
	"regexp"
)

// projectMetaDocument is used in conjunction with xml.Unmarshal to extract the meta information
// from a XML document.
type projectMetaDocument struct {
	Project struct {
		ID                 string `xml:"Id,attr"`
		ProjectInformation struct {
			Name   string `xml:"Name,attr"`
			Number uint   `xml:"ProjectId,attr"`
		}
	}
}

// readProjectMeta extracts project meta information from a file.
func readProjectMeta(file *zip.File) (meta projectMetaDocument, err error) {
	r, err := file.Open()
	if err != nil {
		return
	}

	err = xml.NewDecoder(r).Decode(&meta)
	r.Close()

	return
}

// ProjectMeta contains meta information for a project.
type ProjectMeta struct {
	file *zip.File
	ID   string
	Name string
}

// Process the project file in order to retrieve more information.
func (pm *ProjectMeta) Process() (Project, error) {
	return processZippedProject(pm.file)
}

// A ProjectArchive is a handle to a .knxproj file.
type ProjectArchive struct {
	archive  *zip.ReadCloser
	Projects []ProjectMeta
}

// projectFileRe is a regular expression that matches against project meta files.
var projectFileRe = regexp.MustCompile("^(p|P)-([0-9a-zA-Z]+)/(p|P)roject.xml$")

// Open a .knxproj file.
func Open(file string) (*ProjectArchive, error) {
	archive, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}

	var projects []ProjectMeta

outerLoop:
	for _, file := range archive.File {
		// Is this a project meta file?
		if projectFileRe.MatchString(file.Name) {
			meta, err := readProjectMeta(file)
			if err != nil {
				archive.Close()
				return nil, err
			}

			// Determine the file name which contains the actual project.
			projectFile := path.Join(
				path.Dir(file.Name),
				fmt.Sprintf("%d.xml", meta.Project.ProjectInformation.Number),
			)

			for _, file := range archive.File {
				if file.Name == projectFile {
					projects = append(projects, ProjectMeta{
						file: file,
						ID:   meta.Project.ID,
						Name: meta.Project.ProjectInformation.Name,
					})

					continue outerLoop
				}
			}

			archive.Close()
			return nil, fmt.Errorf("Could not find file: %s", projectFile)
		}
	}

	return &ProjectArchive{
		archive:  archive,
		Projects: projects,
	}, nil
}

// Close the handle.
func (pz *ProjectArchive) Close() error {
	return pz.archive.Close()
}
