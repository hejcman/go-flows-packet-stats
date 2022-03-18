package go_flows_packet_stats

import (
	"github.com/google/gopacket/layers"
	"time"

	"github.com/CN-TU/go-flows/flows"
	"github.com/CN-TU/go-flows/packet"
)

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Common pstats │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

// pstats represents the basic struct used by all the PPI* features.
type pstats struct {
	flows.BaseFeature

	// maxElemCount contains the number of packets, for which stats should be gathered.
	// In essence, it is the length of all the IPFIX lists outputed by this feature.
	maxElemCount uint
	// skipZeroes discards packets with a payload length of 0, so that the output lists
	// will not contain any 0 elements.
	skipZeroes bool
	// skipDup skips duplicated (retransmitted) TCP packets.
	skipDup bool
}

// makePstats is used to initialize a common base structure for all PPI* features.
func makePstats() pstats {
	p := pstats{}
	p.maxElemCount = 30
	p.skipZeroes = true
	p.skipDup = false
	return p
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Packet Lengths  │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

type pktLengths struct {
	pstats
	pktLengths []uint16
}

func (p *pktLengths) Event(new interface{}, _ *flows.EventContext, _ interface{}) {

	if uint(len(p.pktLengths)) == p.maxElemCount {
		return
	}

	buf := new.(packet.Buffer)
	if p.skipZeroes && buf.PayloadLength() == 0 {
		return
	}

	p.pktLengths = append(p.pktLengths, uint16(buf.PayloadLength()))
}

func (p *pktLengths) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.pktLengths, context, p)
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Packet Times │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

type pktTimes struct {
	pstats
	pktTimes []time.Time
}

func (p *pktTimes) Event(new interface{}, _ *flows.EventContext, _ interface{}) {
	if uint(len(p.pktTimes)) == p.maxElemCount {
		return
	}

	buf := new.(packet.Buffer)
	if p.skipZeroes && buf.PayloadLength() == 0 {
		return
	}

	p.pktTimes = append(p.pktTimes, buf.Metadata().Timestamp)
}

func (p *pktTimes) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.pktTimes, context, p)
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Packet directions │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

type pktDirections struct {
	pstats

	// pktDirections holds information about the direction of the first packets. The following directions are possible:
	// Client ⟶ Server = 1
	// Server ⟶ Client = -1
	pktDirections []int8
}

func (p *pktDirections) Event(new interface{}, _ *flows.EventContext, _ interface{}) {
	if uint(len(p.pktDirections)) == p.maxElemCount {
		return
	}

	buf := new.(packet.Buffer)
	if p.skipZeroes && buf.PayloadLength() == 0 {
		return
	}

	// TODO: Use context?
	var dir int8
	if buf.LowToHigh() {
		dir = -1
	} else {
		dir = 1
	}
	p.pktDirections = append(p.pktDirections, dir)
}

func (p *pktDirections) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.pktDirections, context, p)
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╮
│ TCP Flags │
╰╴╴╴╴╴╴╴╴╴╴╴╯
*/

type pktFlags struct {
	pstats
	pktFlags []uint16
}

func (p *pktFlags) Event(new interface{}, _ *flows.EventContext, _ interface{}) {
	if uint(len(p.pktFlags)) == p.maxElemCount {
		return
	}

	buf := new.(packet.Buffer)
	if p.skipZeroes && buf.PayloadLength() == 0 {
		return
	}

	tcp := buf.TransportLayer()
	if tcp == nil {
		return
	}

	tcpContents, _ := tcp.(*layers.TCP)
	if tcpContents == nil {
		return
	}
	var tmpFlag uint8 = 0

	// ╭╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶┬╶╶╶╶╶╮
	// │ Flag │ Binary    │ Dec │
	// ├╶╶╶╶╶╶┼╶╶╶╶╶╶╶╶╶╶╶┼╶╶╶╶╶┤
	// │ FIN  │ 0000 0001 │   1 │
	// │ SYN  │ 0000 0010 │   2 │
	// │ RST  │ 0000 0100 │   4 │
	// │ PSH  │ 0000 1000 │   8 │
	// │ ACK  │ 0001 0000 │  16 │
	// │ URG  │ 0010 0000 │  32 │
	// │ ECE  │ 0100 0000 │  64 │
	// │ CWR  │ 1000 0000 │ 128 │
	// ╰╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴┴╴╴╴╴╴╯

	if tcpContents.FIN {
		tmpFlag |= 1
	}
	if tcpContents.SYN {
		tmpFlag |= 2
	}
	if tcpContents.RST {
		tmpFlag |= 4
	}
	if tcpContents.PSH {
		tmpFlag |= 8
	}
	if tcpContents.ACK {
		tmpFlag |= 16
	}
	if tcpContents.URG {
		tmpFlag |= 32
	}
	if tcpContents.ECE {
		tmpFlag |= 64
	}
	if tcpContents.CWR {
		tmpFlag |= 128
	}
	p.pktFlags = append(p.pktFlags, uint16(tmpFlag))
}

func (p *pktFlags) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.pktFlags, context, p)
}
