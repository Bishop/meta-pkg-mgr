package main

import (
	"fmt"
	"strconv"
)

type Hash map[string]string
type OutdatedRecords map[string]*PkgItem

type PkgItem struct {
	Manager   string
	Package   string
	Current   string
	Available string
	Family    string
	Major     string
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

	if value, ok := piece["family"]; ok {
		p.Family = value
	}

	if value, ok := piece["major"]; ok {
		p.Major = value
	}
}

func (p *PkgItem) HasMajor() bool {
	return p.Major != ""
}

func (p *PkgItem) IsYoungerThan(major string) bool {
	our, _ := strconv.ParseFloat(p.Major, 32)
	their, _ := strconv.ParseFloat(major, 32)

	return our > their
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
	installed := make(Hash)

	for _, item := range o {
		if item.HasMajor() && item.Installed() && item.IsYoungerThan(installed[item.Family]) {
			installed[item.Family] = item.Major
		}
	}

	for name, item := range o {
		major := item.HasMajor() && installed[item.Family] != "" && item.IsYoungerThan(installed[item.Family])

		if !(item.Installed() || major) || item.UpToDate() {
			delete(o, name)
		}
	}
}
