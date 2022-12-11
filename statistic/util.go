package statistic

import (
	"log"
	"net"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/shadowsocks/go-shadowsocks2/socks"
)

func parseSocksAddr(target socks.Addr) *Metadata {
	metadata := &Metadata{}

	switch target[0] {
	case socks.AtypDomainName:
		// trim for FQDN
		metadata.Host = strings.TrimRight(string(target[2:2+target[1]]), ".")
		metadata.DstPort = strconv.Itoa((int(target[2+target[1]]) << 8) | int(target[2+target[1]+1]))
	case socks.AtypIPv4:
		ip := net.IP(target[1 : 1+net.IPv4len])
		metadata.DstIP = ip
		metadata.DstPort = strconv.Itoa((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
	case socks.AtypIPv6:
		ip := net.IP(target[1 : 1+net.IPv6len])
		metadata.DstIP = ip
		metadata.DstPort = strconv.Itoa((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
	}

	return metadata
}

func parseAddr(addr string) (net.IP, string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, "", err
	}

	ip := net.ParseIP(host)
	return ip, port, nil
}

func defaultPath() string {
	u, err := user.Current()
	if err != nil {
		log.Println("获取用户出错")
		return "./"
	}
	return path.Join(u.HomeDir, "orgweb")
}
