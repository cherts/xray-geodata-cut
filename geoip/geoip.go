package geoip

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xtls/xray-core/app/router"
	"github.com/xtls/xray-core/common/net"
	"google.golang.org/protobuf/proto"
)

func LoadGeoIP(fn string) (*router.GeoIPList, error) {
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
	geoIPBytes, err1 := io.ReadAll(file)
	if err1 != nil {
		return nil, err1
	}
	var geoIPList router.GeoIPList
	if err2 := proto.Unmarshal(geoIPBytes, &geoIPList); err2 != nil {
		return nil, err2
	}
	return &geoIPList, nil
}

func GetGeoIPCodes(in *router.GeoIPList) []string {
	result := make([]string, len(in.GetEntry()))
	for index, x := range in.GetEntry() {
		result[index] = x.CountryCode
	}
	return result
}

func CutGeoIPCodes(in *router.GeoIPList, codesToKeep []string, trimIPv6 bool) *router.GeoIPList {
	out := &router.GeoIPList{
		Entry: make([]*router.GeoIP, 0, len(codesToKeep)+1),
	}
	kept := make(map[string]bool, len(codesToKeep)+1)
	for _, x := range in.GetEntry() {
		for _, y := range codesToKeep {
			u := strings.ToUpper(x.CountryCode)
			switch u {
			case strings.ToUpper(strings.TrimSpace(y)), "PRIVATE":
				{
					if kept[u] {
						continue
					}
					if trimIPv6 {
						newEntry := &router.GeoIP{
							ReverseMatch: x.GetReverseMatch(),
							CountryCode:  x.GetCountryCode(),
							Cidr:         make([]*router.CIDR, 0, len(x.GetCidr())),
						}
						for _, c := range x.Cidr {
							if len(c.Ip) == net.IPv4len {
								newEntry.Cidr = append(newEntry.Cidr, c)
							}
						}
						out.Entry = append(out.Entry, newEntry)
					} else {
						out.Entry = append(out.Entry, x)
					}
					kept[u] = true
				}
			}
		}
	}

	return out
}

func SaveGeoIP(in *router.GeoIPList, fn string) error {
	b, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	return os.WriteFile(fn, b, 0600)
}
