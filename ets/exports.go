// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"path"
	"regexp"
)

// ProjectFile is a project file.
type ProjectFile struct {
	*zip.File

	ProjectID   string
	ProjectName string
}

// newProjectFile inspects a project meta file in order to find the real project file.
func newProjectFile(archive *zip.ReadCloser, metaFile *zip.File) (proj ProjectFile, err error) {
	r, err := metaFile.Open()
	if err != nil {
		return
	}

	var meta struct {
		Project struct {
			ID                 string `xml:"Id,attr"`
			ProjectInformation struct {
				Name   string `xml:"Name,attr"`
				Number uint   `xml:"ProjectId,attr"`
			}
		}
	}

	// Extract information from the meta file.
	err = xml.NewDecoder(r).Decode(&meta)
	r.Close()

	if err != nil {
		return
	}

	proj.ProjectID = meta.Project.ID
	proj.ProjectName = meta.Project.ProjectInformation.Name

	// Determine the path of the project file.
	projectFile := fmt.Sprintf(
		"%s/%d.xml",
		path.Dir(metaFile.Name),
		meta.Project.ProjectInformation.Number,
	)

	// Search for the actual project file.
	for _, file := range archive.File {
		if file.Name == projectFile {
			proj.File = file
			return
		}
	}

	return proj, fmt.Errorf("Could not find '%s' from project '%s'", projectFile, proj.ProjectID)
}

// ManufacturerFile is a manufacturer file.
type ManufacturerFile struct {
	*zip.File

	ManufacturerID string
	ContentID      string
}

// ExportArchive is a handle to an exported archive (.knxproj or .knxprod).
type ExportArchive struct {
	archive *zip.ReadCloser

	ProjectFiles      []ProjectFile
	ManufacturerFiles []ManufacturerFile
}

// OpenExportArchive opens the exported archive located at given path.
func OpenExportArchive(path string) (*ExportArchive, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	ex := &ExportArchive{archive: archive}

	if err = ex.findFiles(); err != nil {
		archive.Close()
		return nil, err
	}

	return ex, nil
}

var (
	projectMetaFileRe  = regexp.MustCompile("^(p|P)-([0-9a-zA-Z]+)/(p|P)roject.xml$")
	manufacturerFileRe = regexp.MustCompile("^(m|M)-([0-9a-zA-Z]+)/(m|M)-([0-9a-zA-Z]+)([^.]+).xml$")

	// TODO: Figure out if '/' is a universal path seperator in ZIP files.
)

// findFiles goes through the list of files inside the archive in order to find relevant files.
func (ex *ExportArchive) findFiles() error {
	for _, file := range ex.archive.File {
		if projectMetaFileRe.MatchString(file.Name) {
			// Process meta file in order to find the project file.
			projFile, err := newProjectFile(ex.archive, file)
			if err != nil {
				return err
			}

			ex.ProjectFiles = append(ex.ProjectFiles, projFile)
		} else if matches := manufacturerFileRe.FindStringSubmatch(file.Name); matches != nil {
			ex.ManufacturerFiles = append(ex.ManufacturerFiles, ManufacturerFile{
				File:           file,
				ManufacturerID: "M-" + matches[2],
				ContentID:      "M-" + matches[4] + matches[5],
			})
		}
	}

	return nil
}

// Close the archive handle.
func (ex *ExportArchive) Close() error {
	return ex.archive.Close()
}
