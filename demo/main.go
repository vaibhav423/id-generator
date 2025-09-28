package main

import (
	"bufio"
	"flag"
	"fmt"
	"idgenerator"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// Define command-line flags
	machineID := flag.Uint("m", 1, "Machine ID (0-1023)")
	numIDs := flag.Int("n", 100, "Number of IDs to generate")
	outputFile := flag.String("o", "generated_ids.txt", "Output file name")
	flag.Parse()

	// Validate Machine ID
	if *machineID > idgenerator.MaxMachineID {
		log.Fatalf("Error: Machine ID %d is out of the valid range (0-%d).\n", *machineID, idgenerator.MaxMachineID)
	}

	// 1. Initialize the ID Generator with the provided machine ID
	st := idgenerator.Settings{
		StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			return uint16(*machineID), nil
		},
	}
	gen, err := idgenerator.New(st)
	if err != nil {
		log.Fatalf("Failed to create ID generator: %v", err)
	}

	// 2. Create the output file
	file, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	fmt.Printf("Generating %d IDs with Machine ID %d into file '%s'...\n", *numIDs, *machineID, *outputFile)

	// 3. Generate and write IDs to the file
	start := time.Now()
	for i := 0; i < *numIDs; i++ {
		id, err := gen.NextID()
		if err != nil {
			log.Printf("Error generating ID: %v", err)
			continue
		}
		// Write the ID to the file, followed by a newline
		_, err = writer.WriteString(strconv.FormatUint(id, 10) + "\n")
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}

	// Ensure all buffered operations have been applied to the underlying writer
	writer.Flush()
	duration := time.Since(start)

	fmt.Printf("Successfully generated and saved %d IDs in %v.\n", *numIDs, duration)
	fmt.Printf("Check the '%s' file for the results.\n", *outputFile)
}
