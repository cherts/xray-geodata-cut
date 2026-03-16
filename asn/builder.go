package asn

import (
	"fmt"
	"math"
	"net/netip"

	"github.com/xtls/xray-core/app/router"
)

func BuildGeoIp(asn []int32, trimIpv6 bool) (*router.GeoIPList, error) {
	result := &router.GeoIPList{}
	for _, x := range asn {
		data, err := GetAsnData(x)
		if err != nil {
			return nil, err
		}
		entry := &router.GeoIP{
			CountryCode: fmt.Sprintf("AS%d", x),
			Cidr:        make([]*router.CIDR, 0),
		}
		if data.Prefixes != nil && data.Prefixes.Ipv4 != nil {
			for _, y := range data.Prefixes.Ipv4 {
				ip, e1 := netip.ParsePrefix(y)
				if e1 != nil {
					return nil, e1
				}
				b, e2 := ip.Addr().MarshalBinary()
				if e2 != nil {
					return nil, e2
				}
				bits := ip.Bits()
				if bits < 0 || bits > math.MaxUint32 {
					return nil, fmt.Errorf("invalid prefix bits: %d", bits)
				}
				entry.Cidr = append(entry.Cidr, &router.CIDR{
					Ip:     b,
					Prefix: uint32(bits),
				})
			}
		}
		if !trimIpv6 {
			if data.Prefixes != nil && data.Prefixes.Ipv6 != nil {
				for _, y := range data.Prefixes.Ipv6 {
					ip, e1 := netip.ParsePrefix(y)
					if e1 != nil {
						return nil, e1
					}
					b, e2 := ip.Addr().MarshalBinary()
					if e2 != nil {
						return nil, e2
					}
					bits := ip.Bits()
					if bits < 0 || bits > math.MaxUint32 {
						return nil, fmt.Errorf("invalid prefix bits: %d", bits)
					}
					entry.Cidr = append(entry.Cidr, &router.CIDR{
						Ip:     b,
						Prefix: uint32(bits),
					})
				}
			}
		}
		result.Entry = append(result.Entry, entry)
	}
	return result, nil
}
