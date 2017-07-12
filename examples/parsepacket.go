package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func main() {
	r, err := pcapgo.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var (
		eth     layers.Ethernet
		ip4     layers.IPv4
		ip6     layers.IPv6
		tcp     layers.TCP
		payload gopacket.Payload
	)
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet,
		&eth, &ip4, &ip6, &tcp, &payload)
	decoded := make([]gopacket.LayerType, 0, 4)

	for {
		d, ci, err := r.ReadPacketData()
		if err != nil {
			log.Fatal(err)
		}
		err = parser.DecodeLayers(d, &decoded)
		if err != nil {
			log.Fatal(err)
		}

		// TODO:use ip4 & ip6
		_, _ = d, ci

		fmt.Printf("%+v\n", ci)
		for _, typ := range decoded {
			switch typ {
			case layers.LayerTypeEthernet:
				fmt.Println("  Eth ", eth.SrcMAC, eth.DstMAC)
			case layers.LayerTypeIPv4:
				fmt.Println("  IP4 ", ip4.SrcIP, ip4.DstIP)
			case layers.LayerTypeIPv6:
				fmt.Println("  IP6 ", ip6.SrcIP, ip6.DstIP)
			case layers.LayerTypeTCP:
				fmt.Println("  TCP ", tcp.SrcPort, tcp.DstPort)
			}
		}
	}
}
