package go_flows_packet_stats

import (
	"github.com/google/gopacket/layers"
	"time"

	"github.com/CN-TU/go-flows/flows"
	"github.com/CN-TU/go-flows/packet"
	"github.com/CN-TU/go-ipfix"
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
	p.skipZeroes = false
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

func (p *pktLengths) Event(new interface{}, context *flows.EventContext, _ interface{}) {
	if uint(len(p.pktLengths)) != p.maxElemCount {
		buf := new.(packet.Buffer)
		p.pktLengths = append(p.pktLengths, uint16(buf.PayloadLength()))
	}
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
	if uint(len(p.pktTimes)) != p.maxElemCount {
		buf := new.(packet.Buffer)
		p.pktTimes = append(p.pktTimes, buf.Metadata().Timestamp)
	}
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

func (p *pktDirections) Event(new interface{}, context *flows.EventContext, _ interface{}) {
	if uint(len(p.pktDirections)) == p.maxElemCount {
		return
	}

	// TODO: Use context?
	buf := new.(packet.Buffer)
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
	pktFlags []uint8
}

func (p *pktFlags) Event(new interface{}, _ *flows.EventContext, _ interface{}) {
	if uint(len(p.pktFlags)) == p.maxElemCount {
		return
	}

	buf := new.(packet.Buffer)

	if layer := buf.Layer(layers.LayerTypeTCP); layer != nil {

	}
}

func (p *pktFlags) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.pktFlags, context, p)
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Init function │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

func init() {

	// https://github.com/CESNET/libfds/blob/0302f9c3583bd14b96680100e49a99dadf64e13b/config/system/elements/cesnet.xml#L421
	flows.RegisterFeature(
		ipfix.NewBasicList(
			"packetLength",
			ipfix.NewInformationElement(
				"packetLength",
				CesnetPen,
				1013,
				ipfix.Signed16Type,
				0),
			0),
		"sizes of the first packets",
		flows.FlowFeature,
		func() flows.Feature { return &pktLengths{pstats: makePstats()} },
		flows.RawPacket)

	// https://github.com/CESNET/libfds/blob/0302f9c3583bd14b96680100e49a99dadf64e13b/config/system/elements/cesnet.xml#L428
	flows.RegisterFeature(
		ipfix.NewBasicList(
			"packetTime",
			ipfix.NewInformationElement(
				"packetTime",
				CesnetPen,
				1014,
				ipfix.DateTimeMillisecondsType,
				0),
			0),
		"timestamps of the first packets",
		flows.FlowFeature,
		func() flows.Feature { return &pktTimes{pstats: makePstats()} },
		flows.RawPacket)

	// https://github.com/CESNET/libfds/blob/0302f9c3583bd14b96680100e49a99dadf64e13b/config/system/elements/cesnet.xml#L435
	flows.RegisterFeature(
		ipfix.NewBasicList(
			"packetFlag",
			ipfix.NewInformationElement(
				"packetFlag",
				CesnetPen,
				1015,
				ipfix.Unsigned8Type,
				0),
			0),
		"TCP flags for each packet",
		flows.FlowFeature,
		func() flows.Feature { return &pktFlags{pstats: makePstats()} },
		flows.RawPacket)

	// https://github.com/CESNET/libfds/blob/0302f9c3583bd14b96680100e49a99dadf64e13b/config/system/elements/cesnet.xml#L442
	flows.RegisterFeature(
		ipfix.NewBasicList(
			"packetDirection",
			ipfix.NewInformationElement(
				"packetDirection",
				CesnetPen,
				1016,
				ipfix.Signed8Type,
				0),
			0),
		"directions of the first packets",
		flows.FlowFeature,
		func() flows.Feature { return &pktDirections{pstats: makePstats()} },
		flows.RawPacket)

	//flows.RegisterTemporaryFeature(
	//	"__ppiPktLengths",
	//	"sizes of the first packets",
	//	ipfix.BasicListType,
	//	0,
	//	flows.FlowFeature,
	//	func() flows.Feature { return &pktLengths{pstats: makePstats()} },
	//	flows.RawPacket)
	//flows.RegisterTemporaryFeature(
	//	"__ppiPktTimes",
	//	"timestamps of the first packets",
	//	ipfix.BasicListType,
	//	0,
	//	flows.FlowFeature,
	//	func() flows.Feature { return &pktTimes{pstats: makePstats()} },
	//	flows.RawPacket)
	//flows.RegisterTemporaryFeature(
	//	"__ppiPktDirections",
	//	"directions of the first packets",
	//	ipfix.BasicListType,
	//	0,
	//	flows.FlowFeature,
	//	func() flows.Feature { return &pktDirections{pstats: makePstats()} },
	//	flows.RawPacket)
	//flows.RegisterTemporaryFeature(
	//	"__ppiPktFlags",
	//	"TCP flags for each packet",
	//	ipfix.BasicListType,
	//	0,
	//	flows.FlowFeature,
	//	func() flows.Feature { return &pktFlags{pstats: makePstats()} },
	//	flows.RawPacket)
}
