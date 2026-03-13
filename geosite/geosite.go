package geosite

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xtls/xray-core/app/router"
	"google.golang.org/protobuf/proto"
)

func LoadGeoSite(fn string) (*router.GeoSiteList, error) {
	if filepath.IsAbs(fn) || strings.Contains(fn, "..") {
		return nil, fmt.Errorf("invalid file path")
	}
	root, err := os.OpenRoot(".")
	if err != nil {
		return nil, err
	}
	file, err := root.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	geoSiteBytes, err1 := io.ReadAll(file)
	if err1 != nil {
		return nil, err1
	}
	var geoSiteList router.GeoSiteList
	if err2 := proto.Unmarshal(geoSiteBytes, &geoSiteList); err2 != nil {
		return nil, err2
	}
	return &geoSiteList, nil
}

func GetGeoSiteCodes(in *router.GeoSiteList) []string {
	result := make([]string, len(in.GetEntry()))
	for index, x := range in.GetEntry() {
		result[index] = x.CountryCode
	}
	return result
}

func CutGeoSiteCodes(in *router.GeoSiteList, codesToKeep []string) *router.GeoSiteList {
	out := &router.GeoSiteList{
		Entry: make([]*router.GeoSite, 0, len(codesToKeep)),
	}
	kept := make(map[string]bool, len(codesToKeep))
	for _, x := range in.GetEntry() {
		for _, y := range codesToKeep {
			u := strings.ToUpper(y)
			if x.CountryCode == u {
				if kept[u] {
					continue
				}
				out.Entry = append(out.Entry, x)
				kept[u] = true
			}
		}
	}

	return out
}

func SaveGeoSite(in *router.GeoSiteList, fn string) error {
	b, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	return os.WriteFile(fn, b, 0600)
}
