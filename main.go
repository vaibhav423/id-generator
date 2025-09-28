package idgenerator

import (
	"errors"
	"sync"
	"time"
)

// Bit allocation for the 64-bit ID
const (
	BitLenSequence  = 14
	BitLenMachineID = 10
	BitLenTime      = 39
)

// Maximum values for each component
const (
	MaxSequence  = (1 << BitLenSequence) - 1
	MaxMachineID = (1 << BitLenMachineID) - 1
	MaxTime      = (1 << BitLenTime) - 1
)

// Generator time unit (10ms)
const generatorTimeUnit = 10 * time.Millisecond

// Settings contains configuration for the ID generator.
type Settings struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

// IDGenerator represents the generator instance.
type IDGenerator struct {
	mutex       sync.Mutex
	startTime   int64  // Start time in generator units
	elapsedTime int64  // Current elapsed time units
	sequence    uint16 // Sequence counter
	machineID   uint16 // Machine ID
}

// Error definitions
var (
	ErrStartTimeAhead       = errors.New("start time is ahead of current time")
	ErrNoMachineIDProvided  = errors.New("machine ID function not provided")
	ErrInvalidMachineID     = errors.New("invalid machine ID")
	ErrSequenceOverflow     = errors.New("sequence overflow")
	ErrTimeOverflow         = errors.New("time overflow")
)

// New creates and initializes an IDGenerator according to the settings.
func New(st Settings) (*IDGenerator, error) {
// 1. Validate StartTime (if StartTime.After(now) -> return ErrStartTimeAhead).
now := time.Now()
startTime := st.StartTime
if !startTime.IsZero() && startTime.After(now) {
return nil, ErrStartTimeAhead
}

// 2. If StartTime.IsZero(), set default epoch to 2014-09-01 00:00:00 UTC.
if startTime.IsZero() {
startTime = time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC)
}

// 3. Validate that st.MachineID is not nil. If nil -> return ErrNoMachineIDProvided.
if st.MachineID == nil {
return nil, ErrNoMachineIDProvided
}

// 4. Call st.MachineID(), return error if it fails.
machineID, err := st.MachineID()
if err != nil {
return nil, err
}

// Check if machine ID is within valid range
if machineID > MaxMachineID {
return nil, ErrInvalidMachineID
}

// 5. If st.CheckMachineID != nil and it returns false, return ErrInvalidMachineID.
if st.CheckMachineID != nil && !st.CheckMachineID(machineID) {
return nil, ErrInvalidMachineID
}

// 6. Initialize and return *IDGenerator with proper starting sequence & times.
startTimeUnits := toGeneratorTime(startTime)
currentElapsed := currentElapsedTime(startTimeUnits)

return &IDGenerator{
startTime:   startTimeUnits,
elapsedTime: currentElapsed,
sequence:    0,
machineID:   machineID,
}, nil
}

// NextID generates the next unique ID. Must be concurrency-safe.
func (g *IDGenerator) NextID() (uint64, error) {
	// Lock g.mutex at the start of the critical section and defer Unlock.
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Compute current elapsed time units using currentElapsedTime(g.startTime).
	currentElapsed := currentElapsedTime(g.startTime)

	// If current time has moved forward, reset sequence and update elapsed time
	if currentElapsed > g.elapsedTime {
		g.elapsedTime = currentElapsed
		g.sequence = 0
	} else if currentElapsed < g.elapsedTime {
		// We're ahead of real time, need to wait
		overtime := g.elapsedTime - currentElapsed
		time.Sleep(sleepTime(overtime))
		// After sleeping, update to current time
		g.elapsedTime = currentElapsedTime(g.startTime)
		g.sequence = 0
	}

	// Check if sequence would exceed maximum before incrementing
	if g.sequence >= MaxSequence {
		// Sequence would overflow, advance to next time unit
		g.elapsedTime++
		overtime := g.elapsedTime - currentElapsedTime(g.startTime)
		if overtime > 0 {
			time.Sleep(sleepTime(overtime))
		}
		g.sequence = 0 // Reset sequence for new time unit
	}

	// Now increment sequence safely
	g.sequence++

	// Call g.toID() to pack fields and return the value.
	return g.toID()
}

// Helpers you must implement:

// toGeneratorTime converts a time.Time into generator units since Unix epoch.
func toGeneratorTime(t time.Time) int64 {
return t.UnixNano() / int64(generatorTimeUnit)
}

// currentElapsedTime returns elapsed time units since startTime.
func currentElapsedTime(startTime int64) int64 {
now := time.Now()
currentTime := toGeneratorTime(now)
return currentTime - startTime
}

// sleepTime returns the duration to sleep when the generator is ahead of the real current time.
func sleepTime(overtime int64) time.Duration {
// overtime represents how many time units ahead we are
// We need to sleep for that many generator time units
return time.Duration(overtime) * generatorTimeUnit
}

// DecomposedID represents the components of a generated ID.
type DecomposedID struct {
	Time      int64
	MachineID uint16
	Sequence  uint16
}

// Decompose breaks down a uint64 ID into its constituent parts.
func Decompose(id uint64) DecomposedID {
	// The masks are created by taking the maximum value for each part.
	const maskSequence = uint64(MaxSequence)
	const maskMachineID = uint64(MaxMachineID)

	// Extract the parts using bitwise operations.
	sequence := uint16(id & maskSequence)
	machineID := uint16((id >> BitLenSequence) & maskMachineID)
	idTime := int64(id >> (BitLenMachineID + BitLenSequence))

	return DecomposedID{
		Time:      idTime,
		MachineID: machineID,
		Sequence:  sequence,
	}
}

// toID packs internal fields into a uint64 ID, or returns an error.
func (g *IDGenerator) toID() (uint64, error) {
// Check for overflow conditions
if g.elapsedTime > MaxTime {
return 0, ErrTimeOverflow
}
if g.sequence > MaxSequence {
return 0, ErrSequenceOverflow
}
if g.machineID > MaxMachineID {
return 0, ErrInvalidMachineID
}

// Pack the ID: [time(39)] [machineID(10)] [sequence(14)]
id := uint64(g.elapsedTime)<<(BitLenMachineID+BitLenSequence) |
uint64(g.machineID)<<BitLenSequence |
uint64(g.sequence)

return id, nil
}
