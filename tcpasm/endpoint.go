package tcpasm

import (
	"fmt"
	"strconv"

	"github.com/google/gopacket"
)

type Endpoint struct {
	Address string
	Port    uint16
}

func NewEndpoint(net, tcp gopacket.Endpoint) (Endpoint, error) {
	addr := net.String()
	port, err := strconv.ParseUint(tcp.String(), 10, 16)
	if err != nil {
		return Endpoint{}, err
	}
	return Endpoint{
		Address: addr,
		Port:    uint16(port),
	}, nil
}

func (e Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Address, e.Port)
}

func (e Endpoint) PortString() string {
	return strconv.FormatUint(uint64(e.Port), 10)
}
