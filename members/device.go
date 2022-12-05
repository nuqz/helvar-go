package members

// DeviceState a decimal value, that when broken down into its binary form,
// represents the states where each state is represented by 1 or 0.
type DeviceState int64

const (
	OK DeviceState = 0x00000000

	// NSDisabled - device or subdevice has been disabled, usually an IR
	// subdevice or a DMX channel.
	NSDisabled DeviceState = 0x00000001

	// NSLampFailure - unspecified lamp problem.
	NSLampFailure DeviceState = 0x00000002

	// NSMissing - the device previously existed but is not currently present.
	NSMissing DeviceState = 0x00000004

	// NSFaulty - ran out of addresses (DALI subnet) / unknown Digidim control
	// device / DALI load. For example, dimmers, relay units, ballasts etc.
	// Receives messages from control devices and performs the relevant action
	// e.g.  sets the lighting it controls to the relevant level.  Some
	// control may be possible at the device itself. Also known as Control
	// Gear or LIU, see entries for both. that keeps responding with
	// multi-replies.
	NSFaulty DeviceState = 0x00000008

	// NSRefreshing - DALI subnet, DALI load or Digidim control device is
	// being discovered.
	NSRefreshing DeviceState = 0x00000010

	// NSReserved1 - internal use only.
	NSReserved1 DeviceState = 0x00000020

	// NSReserved2 - internal use only.
	NSReserved2 DeviceState = 0x00000040

	// NSReserved3 - internal use only.
	NSReserved3 DeviceState = 0x00000080

	// NSEMResting - the load is intentionally off whilst the control gear.
	// For example, dimmers, relay units, ballasts etc.  Receives messages
	// from control devices, via the router, and performs the relevant action
	// e.g.  sets the lighting (lamps) it controls to the relevant level.
	// Some control may be possible at the device itself. Also known as LIU or
	// Load, see entries for both. is being powered by the emergency supply.
	NSEMResting DeviceState = 0x00000100

	// NSEMReserved - no description given.
	NSEMReserved DeviceState = 0x00000200

	// NSEMInEmergency - no mains power is being supplied.
	NSEMInEmergency DeviceState = 0x00000400

	// NSEMInProlong - mains has been restored but device is still using the
	// emergency supply.
	NSEMInProlong DeviceState = 0x00000800

	// NSEMFTInProgress - the Functional Test is in progress (brief test where
	// the control gear is being powered by the emergency supply).
	NSEMFTInProgress DeviceState = 0x00001000

	// NSEMDTInProgress - the Duration Test is in progress. This test involves
	// operating the control gear using the battery until the battery is
	// completely discharged. The duration that the control gear was
	// operational for is recorded, and then the battery recharges itself
	// from the mains supply.
	NSEMDTInProgress DeviceState = 0x00002000

	// NSEMReserved1 - no description given.
	NSEMReserved1 DeviceState = 0x00004000

	// NSEMReserved2 - no description given.
	NSEMReserved2 DeviceState = 0x00008000

	// NSEMDTPending - the Duration Test has been requested but has not yet
	// commenced. The test can be delayed if the battery is not fully charged.
	NSEMDTPending DeviceState = 0x00010000

	// NSEMending - the Functional Test has been requested but has not yet
	// commenced. The test can be delayed if there is not enough charge in the
	// battery.
	NSEMFTPending DeviceState = 0x00020000

	// NSEMBatteryFail - battery has failed.
	NSEMBatteryFail DeviceState = 0x00040000

	// NSReserved4 - internal use only.
	NSReserved4 DeviceState = 0x00080000

	// NSReserved5 - internal use only.
	NSReserved5 DeviceState = 0x00100000

	// NSEMInhibit - prevents an emergency fitting from going into emergency
	// mode.
	NSEMInhibit DeviceState = 0x00200000

	// NSEMFTRequested - emergency Function Test has been requested.
	NSEMFTRequested DeviceState = 0x00400000

	// NSEMDTRequested - emergency Duration Test has been requested.
	NSEM_DTRequested DeviceState = 0x00800000

	// NSEMUnknown - initial state of an emergency fitting.
	NSEM_Unknown DeviceState = 0x01000000

	// NSOverTemperature - load is over temperature/heating.
	NSOverTemperature DeviceState = 0x02000000

	// NSOverCurrent - too much current is being drawn by the load.
	NSOverCurrent DeviceState = 0x04000000

	// NSCommsError - communications error.
	NSCommsError DeviceState = 0x08000000

	// NSSevereError - indicates that a load is either over temperature or
	// drawing too much current, or both.
	NSSevereError DeviceState = 0x10000000

	// NSBadReply - indicates that a reply to a query was malformed.
	NSBadReply DeviceState = 0x20000000

	// NSReserved6 - no description given.
	NSReserved6 DeviceState = 0x40000000

	// NSDeviceMismatch - the actual load type does not match An attempt to
	// match corresponding items in an Upload Design and a Workgroup Design /
	// Real Workgroup the expected type.
	NSDeviceMismatch DeviceState = 0x80000000
)

// Device represents a device as defined in HelvarNET protocol.
type Device struct {
	Address string
	Name    string
	State   DeviceState
}

// TODO: shortcuts for state checks: IsXxx() bool { ... }
