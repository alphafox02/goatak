package mp

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
)

const (
	APP_PREF = "com.atakmap.app_preferences"
	CIV_PREF = "com.atakmap.app.civ_preferences"
	STREAMS  = "cot_streams"
)

type MissionPackage struct {
	version uint
	params  map[string]string
	files   []FileContent
}

func NewMissionPackage(uid, name string) *MissionPackage {
	return &MissionPackage{
		version: 2,
		params:  map[string]string{"uid": uid, "name": name},
		files:   nil,
	}
}

func (m *MissionPackage) Param(k, v string) {
	m.params[k] = v
}

func (m *MissionPackage) AddFiles(f ...FileContent) {
	m.files = append(m.files, f...)
}

func (m *MissionPackage) Manifest() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("<MissionPackageManifest version=\"%d\">\n", m.version))
	buf.WriteString("<Configuration>\n")

	for k, v := range m.params {
		buf.WriteString(fmt.Sprintf("  <Parameter name=\"%s\" value=\"%s\"/>\n", k, v))
	}

	buf.WriteString("</Configuration>\n")
	buf.WriteString("<Contents>\n")

	for _, v := range m.files {
		buf.WriteString(fmt.Sprintf("  <Content ignore=\"false\" zipEntry=\"%s\"/>\n", v.Name()))
	}

	buf.WriteString("</Contents>\n")
	buf.WriteString("</MissionPackageManifest>")

	return buf.Bytes()
}

func (m *MissionPackage) Create() ([]byte, error) {
	buff := new(bytes.Buffer)
	zipW := zip.NewWriter(buff)

	f, err := zipW.Create("MANIFEST/manifest.xml")
	if err != nil {
		return nil, err
	}

	_, err = f.Write(m.Manifest())

	if err != nil {
		return nil, err
	}

	for _, zf := range m.files {
		f1, err := zipW.Create(zf.Name())
		if err != nil {
			return nil, err
		}

		_, err = f1.Write(zf.Content())
		if err != nil {
			return nil, err
		}
	}

	err = zipW.Close()

	return buff.Bytes(), err
}

type FileContent interface {
	SetName(name string)
	Name() string
	Content() []byte
}

type FsFile struct {
	name string
	data []byte
}

func NewFsFile(name, path string) (*FsFile, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &FsFile{name: name, data: dat}, nil
}

func NewBlobFile(name string, data []byte) *FsFile {
	return &FsFile{name: name, data: data}
}

func (f *FsFile) Name() string {
	return f.name
}

func (f *FsFile) SetName(name string) {
	f.name = name
}

func (f *FsFile) Content() []byte {
	return f.data
}

type PrefFile struct {
	name string
	data map[string]map[string]string
}

func NewPrefFile(name string) *PrefFile {
	return &PrefFile{name: name, data: make(map[string]map[string]string)}
}

func (p *PrefFile) AddParam(pref, k, v string) {
	if _, ok := p.data[pref]; !ok {
		p.data[pref] = make(map[string]string)
	}

	p.data[pref][k] = v
}

func (p *PrefFile) Name() string {
	return p.name
}

func (p *PrefFile) SetName(name string) {
	p.name = name
}

func (p *PrefFile) Content() []byte {
	var sb bytes.Buffer

	sb.WriteString("<?xml version='1.0' standalone='yes'?>\n")
	sb.WriteString("<preferences>")

	for name, data := range p.data {
		sb.WriteString(fmt.Sprintf("<preference version=\"1\" name=\"%s\">\n", name))

		for k, v := range data {
			e := GetEntry(k, v)
			if e != "" {
				sb.WriteString(e + "\n")
			}
		}

		sb.WriteString("</preference>")
	}

	sb.WriteString("</preferences>")

	return sb.Bytes()
}
