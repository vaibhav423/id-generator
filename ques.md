# ID Generator Implementation Challenge

## Overview
Your task is to implement a **concurrent-safe ID generator** in Go that produces unique 64-bit IDs based on specific rules. The generator should be initialized using a `Settings` struct and provide methods for creating new generators and generating IDs.

This project includes a starter file `idgenerator.go` with **function stubs and comments** guiding your implementation.

---

## Problem Description

### Goal
Implement the logic for generating **unique 64-bit IDs** composed of multiple parts:
- **Elapsed Time** since a given start epoch, measured in 10ms units.
- **Machine ID**, which identifies the machine generating the ID.
- **Sequence Number**, which ensures uniqueness when multiple IDs are generated in the same time unit.

---

### What You Need to Implement

#### 1. `Settings` struct
The `Settings` struct configures the generator. It must include:
- `StartTime` (`time.Time`): The epoch from which elapsed time is calculated.
- `MachineID` (`func() (uint16, error)`): A function that returns a unique machine identifier.
- `CheckMachineID` (`func(uint16) bool`, optional): A function to validate the machine ID.

---

#### 2. `IDGenerator` struct
Represents the state of the generator. It should contain:
- `startTime` (int64): The start time in generator units.
- `elapsedTime` (int64): The current elapsed time units.
- `sequence` (uint16): The sequence counter for IDs within the same time unit.
- `machineID` (uint16): The machine ID obtained from `Settings`.
- A mutex to ensure **concurrency safety**.

---

#### 3. Error Conditions
You must define and return appropriate errors in the following cases:
- **`ErrStartTimeAhead`**: When `StartTime` is set in the future compared to the current time.
- **`ErrNoMachineIDProvided`**: When `MachineID` function is not provided in `Settings`.
- **`ErrInvalidMachineID`**: When `CheckMachineID` fails for the returned machine ID.

---

### Functions to Implement

#### `func New(st Settings) (*IDGenerator, error)`
Creates and initializes an `IDGenerator` based on `Settings`.  
Steps:
1. Validate `StartTime`. If it's after the current time, return `ErrStartTimeAhead`.
2. If `StartTime` is zero, use the default epoch: **2014-09-01 00:00:00 UTC**.
3. Validate `MachineID` is provided; return `ErrNoMachineIDProvided` if not.
4. Call `MachineID()` and handle any error.
5. If `CheckMachineID` is provided and returns `false`, return `ErrInvalidMachineID`.
6. Initialize `IDGenerator` fields properly and return it.

---

#### `func (g *IDGenerator) NextID() (uint64, error)`
Generates the next unique ID in a **concurrency-safe** manner.  
Hints:
- Use a mutex to lock the critical section.
- Compute the current elapsed time using `currentElapsedTime()`.
- If the sequence wraps around (goes back to zero), increment elapsed time and sleep using `sleepTime()`.
- Use `g.toID()` to pack fields into a 64-bit ID.

---

### Helper Functions
- `toGeneratorTime(t time.Time) int64`: Convert a `time.Time` to generator time units (10ms units).
- `currentElapsedTime(startTime int64) int64`: Compute elapsed time since `startTime`.
- `sleepTime(overtime int64) time.Duration`: Calculate how long to sleep when sequence overflows within a single time unit.
- `toID() (uint64, error)`: Combine internal fields into a single 64-bit ID.

---

## Bit Allocation (Example)
Although exact allocation can vary, a typical layout could be:
- **Elapsed Time**: 39 bits
- **Machine ID**: 10 bits
- **Sequence**: 14 bits

---

## Constraints
- The generator must be **thread-safe**.
- The IDs must be **monotonically increasing** for the same machine.
- Sleep when necessary to maintain uniqueness if IDs are requested too quickly.

---

## Testing
- Create multiple goroutines calling `NextID()` and verify all IDs are unique.
- Test edge cases such as:
  - Start time in the future.
  - Sequence overflow handling.
  - Machine ID validation failure.

---

### Deliverables
- Complete implementation in `idgenerator.go`.
- Ensure no **race conditions** or uniqueness issues.
