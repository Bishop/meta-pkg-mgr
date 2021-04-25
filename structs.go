package main

import "fmt"

type Hash map[string]string
type OutdatedRecords map[string]*PkgItem

type Config struct {
	PkgConfigs []PkgConfig `json:"pkg"`
}

type PkgConfig struct {
	Name  string       `json:"name"`
	Shell string       `json:"shell"`
	Flow  []PkgCommand `json:"flow"`
}

type PkgCommand struct {
	Command string `json:"cmd"`
	RegExp  string `json:"re"`
}

type PkgItem struct {
	Manager   string
	Package   string
	Current   string
	Available string
}

func (p *PkgItem) Update(piece Hash) {
	if value, ok := piece["pkg"]; ok {
		p.Package = value
	}

	if value, ok := piece["current"]; ok {
		p.Current = value
	}

	if value, ok := piece["available"]; ok {
		p.Available = value
	}
}

func (p *PkgItem) Installed() bool {
	return p.Current != ""
}

func (p *PkgItem) UpToDate() bool {
	return p.Current == p.Available || p.Available == ""
}

func (p *PkgItem) String() string {
	return fmt.Sprintf("%-8s %-20s %-10s %s", p.Manager, p.Package, p.Current, p.Available)
}

func (o OutdatedRecords) Update(mgr string, item Hash) {
	if _, ok := o[item["pkg"]]; !ok {
		o[item["pkg"]] = &PkgItem{Manager: mgr}
	}

	o[item["pkg"]].Update(item)
}

func (o OutdatedRecords) Filter() {
	for name, item := range o {
		if !item.Installed() || item.UpToDate() {
			delete(o, name)
		}
	}
}
