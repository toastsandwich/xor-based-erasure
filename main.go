package main

import (
	"bytes"
	"fmt"
)

type Drive struct {
	ID      int
	Storage *bytes.Buffer
	Failed  bool // Simulate drive failure
}

func NewDrive(id int) *Drive {
	return &Drive{
		ID:      id,
		Storage: bytes.NewBuffer([]byte{}),
		Failed:  false,
	}
}

func (d *Drive) Use(data []byte) error {
	if d.Failed {
		fmt.Println("Drive", d.ID, "is failed! Cannot write.")
		return fmt.Errorf("drive %d is failed", d.ID)
	}
	n, err := d.Storage.Write(data)
	fmt.Println("Written", n, "bytes in Drive", d.ID)
	return err
}

func (d *Drive) Read() []byte {
	if d.Failed {
		fmt.Println("Drive", d.ID, "is failed! Cannot read.")
		return nil
	}
	return d.Storage.Bytes()
}

func (d *Drive) Fail() {
	fmt.Println("Drive", d.ID, "has failed!")
	d.Failed = true
	d.Storage.Reset() // Clear data to simulate failure
}

func main() {
	data := "This is the data that will be distributed and then we will destroy one drive once done! we will recover data"
	drives := []*Drive{}
	for i := 0; i < 4; i++ {
		drives = append(drives, NewDrive(i))
	}

	fmt.Println("Creating 4 drives for storage")

	buffer := bytes.NewBufferString(data)
	buflen := buffer.Len()
	fmt.Println("Buffer length found:", buflen)

	chunkSize := (buflen + 2) / 3 // Ensure all chunks cover full data
	buf := buffer.Bytes()

	// Create a parity buffer
	parity := make([]byte, chunkSize)

	// Distribute data across 3 drives
	for i := 0; i < 3; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > buflen {
			end = buflen
		}

		shard := make([]byte, chunkSize) // Fixed-size buffer
		copy(shard, buf[start:end])      // Copy actual data

		// Store shard in the corresponding drive
		if err := drives[i].Use(shard); err != nil {
			panic(err)
		}

		// Compute XOR for parity
		for j := 0; j < len(shard); j++ {
			parity[j] ^= shard[j]
		}
	}

	// Store the parity in Drive 3
	if err := drives[3].Use(parity); err != nil {
		panic(err)
	}

	fmt.Println("\n=== Simulating Drive Failure ===")
	failedDrive := 1 // Simulate failure of Drive 1
	drives[failedDrive].Fail()

	fmt.Println("\n=== Recovering Lost Data ===")
	recoveredData := make([]byte, chunkSize)

	// Recompute missing data using XOR
	for i := 0; i < chunkSize; i++ {
		recoveredData[i] = drives[3].Read()[i] // Start with parity
		for j := 0; j < 3; j++ {
			if j != failedDrive { // Exclude failed drive
				recoveredData[i] ^= drives[j].Read()[i]
			}
		}
	}

	fmt.Println("Recovered Data:", string(recoveredData))
}
