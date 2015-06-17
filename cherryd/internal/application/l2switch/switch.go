/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package l2switch

import (
	"bytes"
	"fmt"
	"git.sds.co.kr/cherry.git/cherryd/internal/log"
	"git.sds.co.kr/cherry.git/cherryd/internal/network"
	"git.sds.co.kr/cherry.git/cherryd/openflow"
	"git.sds.co.kr/cherry.git/cherryd/protocol"
	"github.com/dlintw/goconf"
	"net"
)

type L2Switch struct {
	conf *goconf.ConfigFile
	log  log.Logger
}

func New(conf *goconf.ConfigFile, log log.Logger) *L2Switch {
	return &L2Switch{
		conf: conf,
		log:  log,
	}
}

func (r *L2Switch) Name() string {
	return "L2Switch"
}

func flood(f openflow.Factory, ingress *network.Port, packet []byte) error {
	inPort := openflow.NewInPort()
	inPort.SetValue(ingress.Number())

	outPort := openflow.NewOutPort()
	outPort.SetFlood()

	action, err := f.NewAction()
	if err != nil {
		return err
	}
	action.SetOutPort(outPort)

	out, err := f.NewPacketOut()
	if err != nil {
		return err
	}
	out.SetInPort(inPort)
	out.SetAction(action)
	out.SetData(packet)

	return ingress.Device().SendMessage(out)
}

func packetout(f openflow.Factory, egress *network.Port, packet []byte) error {
	inPort := openflow.NewInPort()
	inPort.SetController()

	outPort := openflow.NewOutPort()
	outPort.SetValue(egress.Number())

	action, err := f.NewAction()
	if err != nil {
		return err
	}
	action.SetOutPort(outPort)

	out, err := f.NewPacketOut()
	if err != nil {
		return err
	}
	out.SetInPort(inPort)
	out.SetAction(action)
	out.SetData(packet)

	return egress.Device().SendMessage(out)
}

func isBroadcast(eth *protocol.Ethernet) bool {
	return bytes.Compare(eth.DstMAC, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}) == 0
}

type flowParam struct {
	device    *network.Device
	etherType uint16
	inPort    uint32
	outPort   uint32
	srcMAC    net.HardwareAddr
	dstMAC    net.HardwareAddr
}

func installFlow(f openflow.Factory, p flowParam) error {
	inPort := openflow.NewInPort()
	inPort.SetValue(p.inPort)
	match, err := f.NewMatch()
	if err != nil {
		return err
	}
	match.SetInPort(inPort)
	match.SetEtherType(p.etherType)
	match.SetSrcMAC(p.srcMAC)
	match.SetDstMAC(p.dstMAC)

	outPort := openflow.NewOutPort()
	outPort.SetValue(p.outPort)
	action, err := f.NewAction()
	if err != nil {
		return err
	}
	action.SetOutPort(outPort)
	inst, err := f.NewInstruction()
	if err != nil {
		return err
	}
	inst.ApplyAction(action)

	flow, err := f.NewFlowMod(openflow.FlowAdd)
	if err != nil {
		return err
	}
	flow.SetTableID(p.device.FlowTableID())
	flow.SetIdleTimeout(30)
	flow.SetPriority(10)
	flow.SetFlowMatch(match)
	flow.SetFlowInstruction(inst)

	return p.device.SendMessage(flow)
}

func setFlowRule(f openflow.Factory, p flowParam) error {
	// Forward
	if err := installFlow(f, p); err != nil {
		return err
	}
	// Backward
	return installFlow(f, p)
}

type switchParam struct {
	factory   openflow.Factory
	finder    network.Finder
	ethernet  *protocol.Ethernet
	ingress   *network.Port
	egress    *network.Port
	rawPacket []byte
}

func (r *L2Switch) switching(p switchParam) error {
	// Find path between the ingress device and the other one that has that destination node
	path := p.finder.Path(p.ingress.Device().ID(), p.egress.Device().ID())
	if path == nil || len(path) == 0 {
		r.log.Debug(fmt.Sprintf("Not found a path from %v to %v", p.ethernet.SrcMAC, p.ethernet.DstMAC))
		return nil
	}

	inPort := p.ingress.Number()
	// Install bi-directional flow rules into all devices on the path
	for _, v := range path {
		param := flowParam{
			device:    v[0].Device(),
			etherType: p.ethernet.Type,
			inPort:    inPort,
			outPort:   v[0].Number(),
			srcMAC:    p.ethernet.SrcMAC,
			dstMAC:    p.ethernet.DstMAC,
		}
		if err := setFlowRule(p.factory, param); err != nil {
			return err
		}
		inPort = v[1].Number()
	}

	// Set final flow rule on the destination device
	param := flowParam{
		device:    p.egress.Device(),
		etherType: p.ethernet.Type,
		inPort:    inPort,
		outPort:   p.egress.Number(),
		srcMAC:    p.ethernet.SrcMAC,
		dstMAC:    p.ethernet.DstMAC,
	}
	if err := setFlowRule(p.factory, param); err != nil {
		return err
	}

	// Send this ethernet packet directly to the destination node
	return packetout(p.factory, p.egress, p.rawPacket)
}

func (r *L2Switch) localSwitching(p switchParam) error {
	param := flowParam{
		device:    p.ingress.Device(),
		etherType: p.ethernet.Type,
		inPort:    p.ingress.Number(),
		outPort:   p.egress.Number(),
		srcMAC:    p.ethernet.SrcMAC,
		dstMAC:    p.ethernet.DstMAC,
	}
	if err := setFlowRule(p.factory, param); err != nil {
		return err
	}

	// Send this ethernet packet directly to the destination node
	return packetout(p.factory, p.egress, p.rawPacket)
}

func (r *L2Switch) ProcessPacket(factory openflow.Factory, finder network.Finder, eth *protocol.Ethernet, ingress *network.Port) (drop bool, err error) {
	packet, err := eth.MarshalBinary()
	if err != nil {
		return false, err
	}

	dstNode := finder.Node(eth.DstMAC)
	// Unknown node or broadcast request?
	if dstNode == nil || isBroadcast(eth) {
		r.log.Debug(fmt.Sprintf("Broadcasting (dstMAC=%v)", eth.DstMAC))
		return true, flood(factory, ingress, packet)
	}

	param := switchParam{
		factory:   factory,
		finder:    finder,
		ethernet:  eth,
		ingress:   ingress,
		egress:    dstNode.Port(),
		rawPacket: packet,
	}
	// Two nodes on a same switch device?
	if ingress.Device().ID() == dstNode.Port().Device().ID() {
		err = r.localSwitching(param)
	} else {
		err = r.switching(param)
	}
	if err != nil {
		return false, fmt.Errorf("failed to switch a packet: %v", err)
	}

	return true, nil
}

func (r *L2Switch) ProcessEvent(factory openflow.Factory, finder network.Finder, device *network.Device, status openflow.PortStatus) error {
	if status.Port().IsPortDown() || status.Port().IsLinkDown() {
		port := device.Port(status.Port().Number())
		if port == nil {
			return fmt.Errorf("failed to find a port %v on %v", status.Port().Number(), device.ID())
		}
		return r.cleanup(factory, finder, port)
	}

	return nil
}

func (r *L2Switch) cleanup(factory openflow.Factory, finder network.Finder, port *network.Port) error {
	r.log.Debug(fmt.Sprintf("Cleaning up for %v..", port.ID()))

	// We should remove all edges from all switch devices if the port is an edge among two switches.
	// Otherwise, remaining flow rules in switches may result in incorrect packet routing to the
	// disconnected port.
	if finder.IsEdge(port) {
		return r.removeAllFlows(factory, finder)
	}

	nodes := port.Nodes()
	// Remove all flows related with the nodes that are connected to this port
	for _, n := range nodes {
		r.log.Debug(fmt.Sprintf("Removing all flows related with a node %v..", n.MAC()))

		if err := r.removeFlowRules(factory, finder, n.MAC()); err != nil {
			r.log.Err(fmt.Sprintf("Failed to remove flows related with %v: %v", n.MAC(), err))
			continue
		}
	}

	return nil
}

func (r *L2Switch) removeAllFlows(factory openflow.Factory, finder network.Finder) error {
	r.log.Debug("Removing all flows from all devices..")

	// Wildcard match
	match, err := factory.NewMatch()
	if err != nil {
		return err
	}

	devices := finder.Devices()
	for _, d := range devices {
		if err := r.removeFlow(factory, d, match); err != nil {
			r.log.Err(fmt.Sprintf("Failed to remove flows on %v: %v", d.ID(), err))
			continue
		}
	}

	return nil
}

func (r *L2Switch) removeFlowRules(factory openflow.Factory, finder network.Finder, mac net.HardwareAddr) error {
	devices := finder.Devices()
	for _, d := range devices {
		r.log.Debug(fmt.Sprintf("Removing all flows related with a node %v on device %v..", mac, d.ID()))

		// Remove all flow rules whose source MAC address is mac in its flow match
		match, err := factory.NewMatch()
		if err != nil {
			return err
		}
		match.SetSrcMAC(mac)
		if err := r.removeFlow(factory, d, match); err != nil {
			return err
		}

		// Remove all flow rules whose destination MAC address is mac in its flow match
		match, err = factory.NewMatch()
		if err != nil {
			return err
		}
		match.SetDstMAC(mac)
		if err := r.removeFlow(factory, d, match); err != nil {
			return err
		}
	}

	return nil
}

func (r *L2Switch) removeFlow(f openflow.Factory, d *network.Device, match openflow.Match) error {
	r.log.Debug(fmt.Sprintf("Removing flows on device %v..", d.ID()))

	flowmod, err := f.NewFlowMod(openflow.FlowDelete)
	if err != nil {
		return err
	}
	// Remove flows except the table miss flows (Note that MSB of the cookie is a marker)
	flowmod.SetCookieMask(0x1 << 63)
	flowmod.SetTableID(0xFF) // ALL
	flowmod.SetFlowMatch(match)

	return d.SendMessage(flowmod)
}
