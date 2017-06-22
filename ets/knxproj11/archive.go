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

// A ProjectFile contains project information.
type ProjectFile struct {
	file *zip.File
	ID   string
	Name string
}

// Process the project file in order to retrieve more information.
func (pm *ProjectFile) Process() (Project, error) {
	return processZippedProject(pm.file)
}

// findProjectFile analyzes the contents of the given project meta file in order to find the file
// that contains the actual project.
func findProjectFile(archive *zip.ReadCloser, file *zip.File) (_ ProjectFile, err error) {
	meta, err := readProjectMeta(file)
	if err != nil {
		return
	}

	// Determine the file name which contains the actual project.
	projectFile := path.Join(
		path.Dir(file.Name),
		fmt.Sprintf("%d.xml", meta.Project.ProjectInformation.Number),
	)

	// Find the project file.
	for _, file := range archive.File {
		if file.Name == projectFile {
			return ProjectFile{
				file: file,
				ID:   meta.Project.ID,
				Name: meta.Project.ProjectInformation.Name,
			}, nil
		}
	}

	err = fmt.Errorf(
		"Could not find project file '%s' from project '%s'",
		projectFile, meta.Project.ID,
	)
	return
}

// A ManufacturerFile contains manufacturerer data.
type ManufacturerFile struct {
	file *zip.File
}

// A ProjectArchive is a handle to a .knxproj file.
type ProjectArchive struct {
	archive       *zip.ReadCloser
	Projects      []ProjectFile
	Manufacturers []ManufacturerFile
}

var (
	projectFileRe      = regexp.MustCompile("^(p|P)-([0-9a-zA-Z]+)/(p|P)roject.xml$")
	manufacturerFileRe = regexp.MustCompile("^(m|M)-([0-9a-zA-Z]+)/(m|M)-([0-9a-zA-Z]+)([^.]+).xml$")

	// TODO: Figure out if '/' is a universal path seperator in ZIP files.
)

// Open a .knxproj file.
func Open(file string) (*ProjectArchive, error) {
	archive, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}

	var projects []ProjectFile
	var manufacturers []ManufacturerFile

	for _, file := range archive.File {
		if projectFileRe.MatchString(file.Name) {
			proj, err := findProjectFile(archive, file)
			if err != nil {
				archive.Close()
				return nil, err
			}

			projects = append(projects, proj)
		} else if manufacturerFileRe.MatchString(file.Name) {
			manufacturers = append(manufacturers, ManufacturerFile{file: file})
		}
	}

	return &ProjectArchive{
		archive:       archive,
		Projects:      projects,
		Manufacturers: manufacturers,
	}, nil
}

// Close the handle.
func (pz *ProjectArchive) Close() error {
	return pz.archive.Close()
}
