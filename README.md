# xray-geodata-cut

Cut unneeded data from geoip.dat or geosite.dat, or build geoip.dat from ASNs

### Quick start

Download the archive from [releases](https://github.com/cherts/xray-geodata-cut/releases). Unpack the archive.

For Linux:
```bash
wget -qO- https://github.com/cherts/xray-geodata-cut/releases/download/v1.0.2/xray-geodata-cut_1.0.2_linux_$(uname -m).tar.gz | tar xzf - -C /tmp && \
mv /tmp/xray-geodata-cut /usr/sbin
```

For macOS (install to /opt):
```bash
wget -qO- https://github.com/cherts/xray-geodata-cut/releases/download/v1.0.2/xray-geodata-cut_1.0.2_darwin_$(uname -m).tar.gz | tar xzf - -C /tmp && \
sudo mv /tmp/xray-geodata-cut /opt
```

For Windows (install to C:\Windows, run as Administrator):
```bash
wget -qO- https://github.com/cherts/xray-geodata-cut/releases/download/v1.0.2/xray-geodata-cut_1.0.2_windows_x86_64.tar.gz | tar xzf - -C "C:\Windows\Temp" && \
move "C:\Windows\Temp\xray-geodata-cut.exe" "C:\Windows\"
```

Usage options:
```bash
Usage of xray-geodata-cut:
  -in string
        Path to GeoData file / ASNs split by comma
  -keep string
        GeoIP or GeoSite codes to keep (private is always kept for GeoIP) (default "cn,private,geolocation-!cn")
  -out string
        Path to processed file
  -search string
        Search GeoIP or GeoSite Item
  -show
        Print codes in GeoIP or GeoSite file
  -trimipv6
        Trim all IPv6 ranges in GeoIP file
  -type string
        ASN (asn), GeoIP (geoip) or GeoSite (geosite)

```

ASN information comes from [https://github.com/ipverse/as-ip-blocks/](https://github.com/ipverse/as-ip-blocks/)

Examples for search: 

```bash
xray-geodata-cut -type asn -in 24429,4134 -search 106.124.1.2
AS4134

xray-geodata-cut -in /usr/local/share/xray/geoip.dat -type geoip -search 114.114.114.114
CN

xray-geodata-cut -in /usr/local/share/xray/geoip.dat -type geoip -search 192.0.2.1
PRIVATE

xray-geodata-cut -in /usr/local/share/xray/geoip.dat -type geoip -search 127.0.0.1
PRIVATE
TEST

xray-geodata-cut -in /usr/local/share/xray/geosite.dat -type geosite -search bilibili.com
BILIBILI
CN
GEOLOCATION-CN

xray-geodata-cut -in /usr/local/share/xray/geosite.dat -type geosite -search baidu.com
BAIDU
CN
GEOLOCATION-CN

xray-geodata-cut -in /usr/local/share/xray/geosite.dat -type geosite -search youtube.com
CATEGORY-COMPANIES
GEOLOCATION-!CN
GOOGLE
YOUTUBE

xray-geodata-cut -in /usr/local/share/xray/geosite.dat -type geosite -search www.netflix.com
CATEGORY-ENTERTAINMENT
GEOLOCATION-!CN
NETFLIX
```
