package airplaneslive

import (
	"encoding/json"
	"time"
)

// V2Response represents the top-level aircraft v2 response envelope.
type V2Response struct {
	// Array of aircraft records. Keys are omitted when data is unavailable.
	Ac []Aircraft `json:"ac"`

	// "No error" on success.
	Msg string `json:"msg"`

	// Cache time as a UTC time.Time (parsed from milliseconds since Unix epoch).
	Now time.Time `json:"now"`

	// Number of aircraft in ac.
	Total int `json:"total"`

	// Optional cache time as a UTC time.Time pointer (parsed from milliseconds since Unix epoch).
	CTime *time.Time `json:"ctime,omitempty"`

	// Time to build the response, milliseconds.
	Ptime int `json:"ptime,omitempty"`
}

// ClosestTo returns the nearest aircraft, whether airborne or on the ground, to coord.
// It returns nil if the response contains no aircraft.
func (r V2Response) ClosestTo(coord Coordinate) *Aircraft {

	return r.closestTo(coord, false)
}

// ClosestToAirborne returns the nearest airborne aircraft to coord.
// It returns nil if the response contains no airborne aircraft.
func (r V2Response) ClosestToAirborne(coord Coordinate) *Aircraft {

	return r.closestTo(coord, true)
}

func (r V2Response) closestTo(coord Coordinate, airborne bool) *Aircraft {

	const impossibleDistance float64 = 1e6

	if len(r.Ac) == 0 {
		return nil
	}

	var result Aircraft

	min_dist := impossibleDistance

	for _, ac := range r.Ac {
		if airborne && ac.AltBaro.IsGround {
			continue
		}

		dist := coord.DistanceTo(ac.Location())
		if dist < min_dist {
			min_dist = dist
			result = ac
		}
	}

	if min_dist < impossibleDistance {
		return &result
	}

	return nil
}

// UnmarshalJSON decodes the response and converts millisecond timestamps to UTC times.
func (r *V2Response) UnmarshalJSON(data []byte) error {
	type Alias V2Response
	aux := &struct {
		Now   int64  `json:"now"`
		CTime *int64 `json:"ctime,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Now = time.UnixMilli(aux.Now).UTC()

	if aux.CTime != nil {
		t := time.UnixMilli(*aux.CTime).UTC()
		r.CTime = &t
	} else {
		r.CTime = nil
	}

	return nil
}

// MarshalJSON encodes UTC times as millisecond Unix timestamps.
func (r V2Response) MarshalJSON() ([]byte, error) {
	type Alias V2Response

	var ctimeMs *int64
	if r.CTime != nil {
		ms := r.CTime.UnixMilli()
		ctimeMs = &ms
	}

	return json.Marshal(&struct {
		Now   int64  `json:"now"`
		CTime *int64 `json:"ctime,omitempty"`
		*Alias
	}{
		Now:   r.Now.UnixMilli(),
		CTime: ctimeMs,
		Alias: (*Alias)(&r),
	})
}
