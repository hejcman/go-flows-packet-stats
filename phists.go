package go_flows_packet_stats

import (
	"github.com/CN-TU/go-flows/flows"
	"github.com/CN-TU/go-flows/packet"
)

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Common phists │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

// phists represents the basic struct used by all the *Phists* features.
type phists struct {
	flows.BaseFeature

	// direction is based on IANA IPFIX flowDirection element (ID 61). 0x00 signifies ingress, 0x01 signifies egress.
	// Thus, false means ingress (incoming packets) and true means egress (outgoing packets).
	direction bool
}

func makePhists(direction bool) phists {
	p := phists{}
	p.direction = direction
	return p
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Inter packet times │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

type phistsIpt struct {
	phists

	lastEvent flows.DateTimeMilliseconds
	times     []uint16
}

func (p *phistsIpt) Start(context *flows.EventContext) {
	p.BaseFeature.Start(context)
	p.lastEvent = flows.DateTimeMilliseconds(context.When() / 1000000)
}

func (p *phistsIpt) Event(new interface{}, context *flows.EventContext, _ interface{}) {

	if new.(packet.Buffer).PayloadLength() == 0 && includeZeroes == false {
		return
	}
	if new.(packet.Buffer).LowToHigh() == p.direction {
		return
	}

	// Measure the time difference since the last packet.
	time := uint32(context.When()/1000000) - uint32(p.lastEvent)

	// Fit it into a bin based on the duration.
	// Based on bins defined here: https://github.com/CESNET/ipfixprobe#phists
	if time <= 15 {
		p.times[0] = p.times[0] + 1
	} else if 15 < time && time <= 31 {
		p.times[1] = p.times[1] + 1
	} else if 31 < time && time <= 63 {
		p.times[2] = p.times[2] + 1
	} else if 63 < time && time <= 127 {
		p.times[3] = p.times[3] + 1
	} else if 127 < time && time <= 255 {
		p.times[4] = p.times[4] + 1
	} else if 255 < time && time <= 511 {
		p.times[5] = p.times[5] + 1
	} else if 511 < time && time <= 1023 {
		p.times[6] = p.times[6] + 1
	} else {
		p.times[7] = p.times[7] + 1
	}

	// Update the last event time.
	p.lastEvent = flows.DateTimeMilliseconds(context.When() / 1000000)
}

func (p *phistsIpt) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.times, context, p)
}

/*
╭╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╶╮
│ Inter packet sizes │
╰╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╴╯
*/

type phistsSizes struct {
	phists

	sizes []uint16
}

func (p *phistsSizes) Event(new interface{}, _ *flows.EventContext, _ interface{}) {

	if new.(packet.Buffer).PayloadLength() == 0 && includeZeroes == false {
		return
	}
	if new.(packet.Buffer).LowToHigh() == p.direction {
		return
	}

	size := new.(packet.Buffer).PayloadLength()

	if size <= 15 {
		p.sizes[0] = p.sizes[0] + 1
	} else if 15 < size && size <= 31 {
		p.sizes[1] = p.sizes[1] + 1
	} else if 31 < size && size <= 63 {
		p.sizes[2] = p.sizes[2] + 1
	} else if 63 < size && size <= 127 {
		p.sizes[3] = p.sizes[3] + 1
	} else if 127 < size && size <= 255 {
		p.sizes[4] = p.sizes[4] + 1
	} else if 255 < size && size <= 511 {
		p.sizes[5] = p.sizes[5] + 1
	} else if 511 < size && size <= 1023 {
		p.sizes[6] = p.sizes[6] + 1
	} else {
		p.sizes[7] = p.sizes[7] + 1
	}
}

func (p *phistsSizes) Stop(_ flows.FlowEndReason, context *flows.EventContext) {
	p.BaseFeature.SetValue(p.sizes, context, p)
}
