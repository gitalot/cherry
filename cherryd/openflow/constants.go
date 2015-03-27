/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package openflow

// ofp_type
const (
	OFPT_HELLO = iota
	OFPT_ERROR
	OFPT_ECHO_REQUEST
	OFPT_ECHO_REPLY
	OFPT_VENDOR
	OFPT_FEATURES_REQUEST
	OFPT_FEATURES_REPLY
	OFPT_GET_CONFIG_REQUEST
	OFPT_GET_CONFIG_REPLY
	OFPT_SET_CONFIG
	OFPT_PACKET_IN
	OFPT_FLOW_REMOVED
	OFPT_PORT_STATUS
	OFPT_PACKET_OUT
	OFPT_FLOW_MOD
	OFPT_PORT_MOD
	OFPT_STATS_REQUEST
	OFPT_STATS_REPLY
	OFPT_BARRIER_REQUEST
	OFPT_BARRIER_REPLY
	OFPT_QUEUE_GET_CONFIG_REQUEST
	OFPT_QUEUE_GET_CONFIG_REPLY
)

// ofp_error_type
const (
	OFPET_HELLO_FAILED = iota
	OFPET_BAD_REQUEST
	OFPET_BAD_ACTION
	OFPET_FLOW_MOD_FAILED
	OFPET_PORT_MOD_FAILED
	OFPET_QUEUE_OP_FAILED
)

// ofp_hello_failed_code
const (
	OFPHFC_INCOMPATIBLE = iota
	OFPHFC_EPERM
)

// ofp_bad_request_code
const (
	OFPBRC_BAD_VERSION = iota
	OFPBRC_BAD_TYPE
	OFPBRC_BAD_STAT
	OFPBRC_BAD_VENDOR
	OFPBRC_BAD_SUBTYPE
	OFPBRC_EPERM
	OFPBRC_BAD_LEN
	OFPBRC_BUFFER_EMPTY
	OFPBRC_BUFFER_UNKNOWN
)

// ofp_bad_action_code
const (
	OFPBAC_BAD_TYPE = iota
	OFPBAC_BAD_LEN
	OFPBAC_BAD_VENDOR
	OFPBAC_BAD_VENDOR_TYPE
	OFPBAC_BAD_OUT_PORT
	OFPBAC_BAD_ARGUMENT
	OFPBAC_EPERM
	OFPBAC_TOO_MANY
	OFPBAC_BAD_QUEUE
)

// ofp_flow_mod_failed_code
const (
	OFPFMFC_ALL_TABLES_FULL = iota
	OFPFMFC_OVERLAP
	OFPFMFC_EPERM
	OFPFMFC_BAD_EMERG_TIMEOUT
	OFPFMFC_BAD_COMMAND
	OFPFMFC_UNSUPPORTED
)

// ofp_port_mod_failed_code
const (
	FPPMFC_BAD_PORT = iota
	OFPPMFC_BAD_HW_ADDR
)

// ofp_queue_op_failed_code
const (
	OFPQOFC_BAD_PORT = iota
	OFPQOFC_BAD_QUEUE
	OFPQOFC_EPERM
)

// ofp_port_config
const (
	OFPPC_PORT_DOWN    = 1 << 0
	OFPPC_NO_STP       = 1 << 1
	OFPPC_NO_RECV      = 1 << 2
	OFPPC_NO_RECV_STP  = 1 << 3
	OFPPC_NO_FLOOD     = 1 << 4
	OFPPC_NO_FWD       = 1 << 5
	OFPPC_NO_PACKET_IN = 1 << 6
)

// ofp_port_state
const (
	OFPPS_LINK_DOWN = 1 << 0
)

// ofp_port_features
const (
	OFPPF_10MB_HD    = 1 << 0  /* 10 Mb half-duplex rate support. */
	OFPPF_10MB_FD    = 1 << 1  /* 10 Mb full-duplex rate support. */
	OFPPF_100MB_HD   = 1 << 2  /* 100 Mb half-duplex rate support. */
	OFPPF_100MB_FD   = 1 << 3  /* 100 Mb full-duplex rate support. */
	OFPPF_1GB_HD     = 1 << 4  /* 1 Gb half-duplex rate support. */
	OFPPF_1GB_FD     = 1 << 5  /* 1 Gb full-duplex rate support. */
	OFPPF_10GB_FD    = 1 << 6  /* 10 Gb full-duplex rate support. */
	OFPPF_COPPER     = 1 << 7  /* Copper medium. */
	OFPPF_FIBER      = 1 << 8  /* Fiber medium. */
	OFPPF_AUTONEG    = 1 << 9  /* Auto-negotiation. */
	OFPPF_PAUSE      = 1 << 10 /* Pause. */
	OFPPF_PAUSE_ASYM = 1 << 11 /* Asymmetric pause. */
)

// ofp_capabilities
const (
	OFPC_FLOW_STATS   = 1 << 0 /* Flow statistics. */
	OFPC_TABLE_STATS  = 1 << 1 /* Table statistics. */
	OFPC_PORT_STATS   = 1 << 2 /* Port statistics. */
	OFPC_STP          = 1 << 3 /* 802.1d spanning tree. */
	OFPC_RESERVED     = 1 << 4 /* Reserved, must be zero. */
	OFPC_IP_REASM     = 1 << 5 /* Can reassemble IP fragments. */
	OFPC_QUEUE_STATS  = 1 << 6 /* Queue statistics. */
	OFPC_ARP_MATCH_IP = 1 << 7 /* Match IP addresses in ARP pkts. */
)

// ofp_action_type
const (
	OFPAT_OUTPUT       = iota /* Output to switch port. */
	OFPAT_SET_VLAN_VID        /* Set the 802.1q VLAN id. */
	OFPAT_SET_VLAN_PCP        /* Set the 802.1q priority. */
	OFPAT_STRIP_VLAN          /* Strip the 802.1q header. */
	OFPAT_SET_DL_SRC          /* Ethernet source address. */
	OFPAT_SET_DL_DST          /* Ethernet destination address. */
	OFPAT_SET_NW_SRC          /* IP source address. */
	OFPAT_SET_NW_DST          /* IP destination address. */
	OFPAT_SET_NW_TOS          /* IP ToS (DSCP field, 6 bits). */
	OFPAT_SET_TP_SRC          /* TCP/UDP source port. */
	OFPAT_SET_TP_DST          /* TCP/UDP destination port. */
	OFPAT_ENQUEUE             /* Output to queue. */
	OFPAT_VENDOR       = 0xffff
)
