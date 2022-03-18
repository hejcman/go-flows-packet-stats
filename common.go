package go_flows_packet_stats

import (
	"github.com/CN-TU/go-flows/flows"
	"github.com/CN-TU/go-ipfix"
)

// TODO: Allow setting "includeZeroes" based on an argument.

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Common variables and definitions │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

// CesnetPen is the pen of CESNET, as defined by RFC 7013
// https://datatracker.ietf.org/doc/html/draft-ietf-ipfix-ie-doctors#section-10.1
var CesnetPen uint32 = 8057

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Init function │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

func init() {

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PHists │ Sizes (SRC -> DST) │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

	flows.RegisterFeature(
		ipfix.NewBasicList(
			"phistSrcSizes",
			ipfix.NewInformationElement(
				"phistSrcSizes",
				CesnetPen,
				1060,
				ipfix.Unsigned16Type,
				0),
			8),
		"histogram of interpacket sizes, SRC -> DST",
		flows.FlowFeature,
		func() flows.Feature {
			return &phistsSizes{
				phists: makePhists(true, true),
				sizes:  []uint16{0, 0, 0, 0, 0, 0, 0, 0}}
		},
		flows.RawPacket)

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PHists │ Inter Packet Time (SRC -> DST) │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

	flows.RegisterFeature(
		ipfix.NewBasicList(
			"phistSrcInterPacketTime",
			ipfix.NewInformationElement(
				"phistSrcInterPacketTime",
				CesnetPen,
				1061,
				ipfix.Unsigned16Type,
				0),
			8),
		"histogram of interpacket times, SRC -> DST",
		flows.FlowFeature,
		func() flows.Feature {
			return &phistsIpt{
				phists: makePhists(true, true),
				times:  []uint16{0, 0, 0, 0, 0, 0, 0, 0}}
		},
		flows.RawPacket)

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PHists │ Sizes (DST -> SRC) │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

	flows.RegisterFeature(
		ipfix.NewBasicList(
			"phistDstSizes",
			ipfix.NewInformationElement(
				"phistDstSizes",
				CesnetPen,
				1062,
				ipfix.Unsigned16Type,
				0),
			8),
		"histogram of interpacket sizes, DST -> SRC",
		flows.FlowFeature,
		func() flows.Feature {
			return &phistsSizes{
				phists: makePhists(false, true),
				sizes:  []uint16{0, 0, 0, 0, 0, 0, 0, 0}}
		},
		flows.RawPacket)

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PHists │ Inter Packet Time (DST -> SRC) │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

	flows.RegisterFeature(
		ipfix.NewBasicList(
			"phistDstInterPacketTime",
			ipfix.NewInformationElement(
				"phistDstInterPacketTime",
				CesnetPen,
				1063,
				ipfix.Unsigned32Type,
				0),
			8),
		"histogram of interpacket times, DST -> SRC",
		flows.FlowFeature,
		func() flows.Feature {
			return &phistsIpt{
				phists: makePhists(false, true),
				times:  []uint16{0, 0, 0, 0, 0, 0, 0, 0}}
		},
		flows.RawPacket)

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PStats │ Packet Payload Length │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

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

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PStats │ Packet Payload Length │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

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

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PStats │ TCP Packet Flags │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

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

	// ╭╶╶╶╶╶╶╶╶┬╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
	// │ PStats │ Packet directions │
	// ╰╴╴╴╴╴╴╴╴┴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯

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
}
