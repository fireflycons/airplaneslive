package airplaneslive

import (
	"encoding/json"
	"fmt"
)

// Pitch describes an aircraft's vertical direction.
// Actual values can be obtained from [Aircraft] properties.
type Pitch int

const (
	// PitchDecending indicates a negative barometric climb rate.
	PitchDecending = iota - 1
	// PitchLevel indicates no barometric climb or descent.
	PitchLevel
	// PitchAscending indicates a positive barometric climb rate.
	PitchAscending
)

// String returns the textual description of the pitch direction.
func (p Pitch) String() string {
	switch p {
	case PitchAscending:
		return "ascending"
	case PitchLevel:
		return "level"
	case PitchDecending:
		return "descending"
	default:
		return "unknown"
	}
}

// Bank describes an aircraft's turn state (roll).
// Actual values can be obtained from [Aircraft] properties.
type Bank int

const (
	// BankLeft indicates a left roll.
	BankLeft = iota - 1
	// BankNone indicates no roll.
	BankNone
	// BankRight indicates a right roll.
	BankRight
)

// String returns the textual description of the aircraft's bank direction.
func (b Bank) String() string {

	switch b {
	case BankLeft:
		return "left"
	case BankNone:
		return "none"
	case BankRight:
		return "right"
	default:
		return "unknown"
	}
}

// Aircraft represents an individual aircraft record containing position, identification, and telemetry data.
type Aircraft struct {
	// 24-bit Mode S / ICAO address (6 hex digits).
	Hex string `json:"hex"`

	// Message/position source, e.g. adsb_icao, mlat.
	Type string `json:"type,omitempty"`

	// Callsign / flight id, up to 8 chars (DO-260B 2.2.8.2.6).
	Flight string `json:"flight,omitempty"`

	// Registration (from the aircraft DB).
	Registration string `json:"r,omitempty"`

	// ICAO type code (from the aircraft DB).
	TypeCode string `json:"t,omitempty"`

	// Long type description (DB, optional).
	Desc string `json:"desc,omitempty"`

	// Bitfield: 1=military, 2=interesting, 4=PIA, 8=LADD.
	DBFlags int `json:"dbFlags,omitempty"`

	// Barometric altitude in feet, or "ground".
	AltBaro AltBaro `json:"alt_baro,omitempty"`

	// Geometric altitude in feet (WGS84).
	AltGeom int `json:"alt_geom,omitempty"`

	// Ground speed in knots.
	GS float64 `json:"gs,omitempty"`

	// Indicated air speed in knots.
	IAS int `json:"ias,omitempty"`

	// True air speed in knots.
	TAS int `json:"tas,omitempty"`

	// Mach number.
	Mach float64 `json:"mach,omitempty"`

	// True track over ground, degrees.
	Track float64 `json:"track,omitempty"`

	// Rate of change of track, degrees/second.
	TrackRate float64 `json:"track_rate,omitempty"`

	// Roll angle, degrees; negative is left roll.
	Roll float64 `json:"roll,omitempty"`

	// Heading, degrees clockwise from magnetic north.
	MagHeading float64 `json:"mag_heading,omitempty"`

	// Heading, degrees clockwise from true north; usually only transmitted on the ground, otherwise derived from magnetic heading using WMM2020.
	TrueHeading float64 `json:"true_heading,omitempty"`

	// Barometric climb rate, ft/min.
	BaroRate int `json:"baro_rate,omitempty"`

	// Geometric climb rate, ft/min.
	GeomRate int `json:"geom_rate,omitempty"`

	// Mode A code, 4 octal digits.
	Squawk string `json:"squawk,omitempty"`

	// ADS-B emergency/priority status; a superset of the 7x00 squawks (DO-260B 2.2.3.2.7.8.1.1).
	// Allowed values: none, general, lifeguard, minfuel, nordo, unlawful, downed, reserved.
	Emergency string `json:"emergency,omitempty"`

	// Emitter category A0-D7 (DO-260B 2.2.3.2.5.2).
	Category string `json:"category,omitempty"`

	// Altimeter setting (QFE or QNH/QNE), hPa.
	NavQNH float64 `json:"nav_qnh,omitempty"`

	// Selected altitude from the MCP/FCU or equivalent, feet.
	NavAltitudeMCP int `json:"nav_altitude_mcp,omitempty"`

	// Selected altitude from the Flight Management System, feet (DO-260B 2.2.3.2.7.1.3.3).
	NavAltitudeFMS int `json:"nav_altitude_fms,omitempty"`

	// Selected heading, degrees (usually magnetic) (DO-260B 2.2.3.2.7.1.3.7).
	NavHeading float64 `json:"nav_heading,omitempty"`

	// Engaged automation modes. Allowed values: autopilot, vnav, althold, approach, lnav, tcas.
	NavModes []string `json:"nav_modes,omitempty"`

	// Latitude, decimal degrees.
	Lat float64 `json:"lat,omitempty"`

	// Longitude, decimal degrees.
	Lon float64 `json:"lon,omitempty"`

	// Navigation Integrity Category (DO-260B 2.2.3.2.7.2.6).
	Nic int `json:"nic,omitempty"`

	// Radius of containment, meters; position-integrity measure derived from NIC and supplementary bits (DO-260B 2.2.3.2.7.2.6, Table 2-69).
	Rc int `json:"rc,omitempty"`

	// Age of the position, seconds before "now".
	SeenPos float64 `json:"seen_pos,omitempty"`

	// ADS-B version number (0, 1 or 2; 3-7 reserved) (DO-260B 2.2.3.2.7.5).
	Version int `json:"version,omitempty"`

	// Navigation Integrity Category for barometric altitude (DO-260B 2.2.5.1.35).
	NicBaro int `json:"nic_baro,omitempty"`

	// Navigation Accuracy for Position (DO-260B 2.2.5.1.35).
	NacP int `json:"nac_p,omitempty"`

	// Navigation Accuracy for Velocity (DO-260B 2.2.5.1.19).
	NacV int `json:"nac_v,omitempty"`

	// Source Integrity Level (DO-260B 2.2.5.1.40).
	Sil int `json:"sil,omitempty"`

	// Interpretation of SIL. Allowed values: unknown, perhour, persample.
	SilType string `json:"sil_type,omitempty"`

	// Geometric Vertical Accuracy (DO-260B 2.2.3.2.7.2.8).
	Gva int `json:"gva,omitempty"`

	// System Design Assurance (DO-260B 2.2.3.2.7.2.4.6).
	Sda int `json:"sda,omitempty"`

	// Names of fields derived from MLAT data.
	Mlat []string `json:"mlat,omitempty"`

	// Names of fields derived from TIS-B data.
	Tisb []string `json:"tisb,omitempty"`

	// Total Mode S messages from this aircraft.
	Messages int `json:"messages,omitempty"`

	// Age of the last message, seconds before "now".
	Seen float64 `json:"seen,omitempty"`

	// Recent average signal power, dBFS (negative).
	Rssi float64 `json:"rssi,omitempty"`

	// Flight status alert bit (DO-260B 2.2.3.2.3.2).
	Alert int `json:"alert,omitempty"`

	// Flight status special position identification bit (DO-260B 2.2.3.2.3.2).
	Spi int `json:"spi,omitempty"`

	// Wind direction, degrees (calculated).
	Wd int `json:"wd,omitempty"`

	// Wind speed, knots (calculated).
	Ws int `json:"ws,omitempty"`

	// Outer/static air temperature, degrees C, calculated from mach and true airspeed; inaccurate at low altitude/mach and inhibited for mach < 0.395.
	Oat int `json:"oat,omitempty"`

	// Total air temperature, degrees C, calculated from mach and true airspeed; inaccurate at low altitude/mach and inhibited for mach < 0.395.
	Tat int `json:"tat,omitempty"`

	// ACAS resolution advisory (experimental, subject to change).
	AcasRA any `json:"acas_ra,omitempty"`

	// Experimental: timestamp before which GPS was working well; shown for 15 min after GPS is lost or degraded.
	GpsOkBefore float64 `json:"gpsOkBefore,omitempty"`

	// Distance from the query point, NM (point queries only).
	Dst float64 `json:"dst,omitempty"`

	// Bearing from the query point, degrees (point queries only).
	Dir float64 `json:"dir,omitempty"`

	// Last known position when the live position is older than 60s.
	LastPosition *Pos `json:"lastPosition,omitempty"`

	// Rough estimated latitude when no ADS-B or MLAT position is available.
	RRLat float64 `json:"rr_lat,omitempty"`

	// Rough estimated longitude when no ADS-B or MLAT position is available.
	RRLon float64 `json:"rr_lon,omitempty"`
}

// AirlineIcao returns the ICAO airline code derived from the first three characters of the flight callsign.
func (a Aircraft) AirlineIcao() string {
	if len(a.Flight) >= 3 {
		return a.Flight[:3]
	}
	return ""
}

// Location returns the aircraft's lat/long position as a coordinate.
func (a Aircraft) Location() Coordinate {
	return Coordinate{lat: a.Lat, lon: a.Lon}
}

// DistanceFrom calculates the distance the aircarft is from the given point.
func (a Aircraft) DistanceFrom(point Coordinate) float64 {

	return point.DistanceTo(a.Location())
}

// BearingFrom calculates the initial bearing from this point to the aircraft.
func (a Aircraft) BearingFrom(point Coordinate) float64 {

	return point.HeadingTo(a.Location())
}

// Pitch returns the aircraft's vertical direction based on its barometric climb rate.
func (a Aircraft) Pitch() Pitch {
	switch {
	case a.BaroRate < 0:
		return PitchDecending
	case a.BaroRate > 0:
		return PitchAscending
	default:
		return PitchLevel
	}
}

// Bank returns the aircraft's turn direction based on its roll angle.
func (a Aircraft) Bank() Bank {
	switch {
	case a.Roll < 0:
		return BankLeft
	case a.Roll > 0:
		return BankRight
	default:
		return BankNone
	}
}

// Pos represents last known position object details.
type Pos struct {
	// Latitude, decimal degrees.
	Lat float64 `json:"lat,omitempty"`

	// Longitude, decimal degrees.
	Lon float64 `json:"lon,omitempty"`

	// Navigation Integrity Category (DO-260B 2.2.3.2.7.2.6).
	Nic int `json:"nic,omitempty"`

	// Radius of containment, meters (DO-260B 2.2.3.2.7.2.6, Table 2-69).
	Rc int `json:"rc,omitempty"`

	// Age of the position, seconds before "now".
	SeenPos float64 `json:"seen_pos,omitempty"`
}

// AltBaro represents barometric altitude, which may be an integer altitude or the value "ground".
type AltBaro struct {
	// Altitude is the barometric altitude in feet; it is zero when IsGround is true.
	Altitude int
	// IsGround reports whether the aircraft is on the ground.
	IsGround bool
}

// UnmarshalJSON decodes an integer altitude or the string "ground".
func (a *AltBaro) UnmarshalJSON(b []byte) error {
	var num int
	if err := json.Unmarshal(b, &num); err == nil {
		a.Altitude = num
		a.IsGround = false
		return nil
	}

	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		if str == "ground" {
			a.IsGround = true
			a.Altitude = 0
			return nil
		}
	}

	return fmt.Errorf("invalid alt_baro value: %s", string(b))
}

// MarshalJSON encodes the altitude as an integer or the string "ground".
func (a AltBaro) MarshalJSON() ([]byte, error) {
	if a.IsGround {
		return json.Marshal("ground")
	}
	return json.Marshal(a.Altitude)
}
