package main
	
import (
	"fmt"
	"idgenerator"
	"log"
	"sync"
)

	func main() {
	// Create a new ID generator with a simple machine ID function
	gen, err := idgenerator.New(idgenerator.Settings{
		MachineID: func() (uint16, error) {
			// In a real application, this could be derived from the hostname,
			// MAC address, or some other unique machine identifier
			return 42, nil
		},
		CheckMachineID: func(id uint16) bool {
			// Optional validation function
			return id > 0 && id < 1000
		},
	})
	if err != nil {
		log.Fatalf("Failed to create ID generator: %v", err)
	}

	fmt.Println("ID Generator Demo")
	fmt.Println("================")

	// Generate some IDs sequentially
	fmt.Println("\nSequential ID generation:")
	for i := 0; i < 10; i++ {
		id, err := gen.NextID()
		if err != nil {
			log.Fatalf("Failed to generate ID: %v", err)
		}
		fmt.Printf("ID %2d: %d\n", i+1, id)
	}

	// Test concurrent generation
	fmt.Println("\nConcurrent ID generation (10 goroutines, 5 IDs each):")
	var wg sync.WaitGroup
	idChan := make(chan uint64, 50)

	// Launch 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(routineNum int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				id, err := gen.NextID()
				if err != nil {
					log.Printf("Failed to generate ID in routine %d: %v", routineNum, err)
					return
				}
				idChan <- id
			}
		}(i)
	}

	wg.Wait()
	close(idChan)

	// Collect and sort IDs to show they're unique
	ids := make([]uint64, 0, 50)
	for id := range idChan {
		ids = append(ids, id)
	}

	// Simple sort to demonstrate uniqueness
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] > ids[j] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}

	fmt.Printf("Generated %d unique IDs concurrently:\n", len(ids))
	for i, id := range ids {
		fmt.Printf("%2d: %d\n", i+1, id)
	}

	// Demonstrate ID component extraction
	fmt.Println("\nID Component Analysis (last generated ID):")
	if len(ids) > 0 {
		lastID := ids[len(ids)-1]

		// Extract components using bit operations
		sequence := uint16(lastID & idgenerator.MaxSequence)
		machineID := uint16((lastID >> idgenerator.BitLenSequence) & idgenerator.MaxMachineID)
		elapsedTime := int64(lastID >> (idgenerator.BitLenSequence + idgenerator.BitLenMachineID))

		fmt.Printf("Full ID:      %d\n", lastID)
		fmt.Printf("Time:         %d (elapsed 10ms units)\n", elapsedTime)
		fmt.Printf("Machine ID:   %d\n", machineID)
		fmt.Printf("Sequence:     %d\n", sequence)
	}
}
