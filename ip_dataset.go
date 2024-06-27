package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

type IPDataset struct {
	filepath    string
	concurrency int
}

func NewIPDataset(filepath string, concurrency int) *IPDataset {
	return &IPDataset{
		filepath:    filepath,
		concurrency: concurrency,
	}
}

func (d *IPDataset) CountUniqueIPs(bitmap *Bitmap) (counter int64, rerr error) {
	fileChunkOffsets, err := d.calculateFileChunkOffsets()
	if err != nil {
		return 0, fmt.Errorf("failed to calculate file chunks: %w", err)
	}

	var (
		uniqueIPCounter int64
		waitGroup       sync.WaitGroup
	)

	waitGroup.Add(d.concurrency)

	for i := 0; i < d.concurrency; i++ {
		start := fileChunkOffsets[i]
		end := fileChunkOffsets[i+1]

		go func() {
			err := d.processFileChunk(start, end, &waitGroup, bitmap, &uniqueIPCounter)
			if err != nil {
				rerr = errors.Join(rerr, fmt.Errorf("failed to process file chunk: %w", err))
			}
		}()
	}

	waitGroup.Wait()

	return uniqueIPCounter, nil
}

func (d *IPDataset) calculateFileChunkOffsets() (offsets []int64, rerr error) {
	file, err := os.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open a file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			rerr = errors.Join(rerr, fmt.Errorf("failed to close a file: %w", err))
		}
	}()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	fileSize := fileStat.Size()
	fileChunkSize := fileSize / int64(d.concurrency)

	var fileChunkOffsets []int64
	for i := int64(0); i < int64(d.concurrency); i++ {
		fileChunkOffsets = append(fileChunkOffsets, i*fileChunkSize)
	}

	fileChunkOffsets = append(fileChunkOffsets, fileSize)

	return fileChunkOffsets, nil
}

func (d *IPDataset) processFileChunk(
	start,
	end int64,
	waitGroup *sync.WaitGroup,
	bitmap *Bitmap,
	uniqueIPCounter *int64,
) (rerr error) {
	defer waitGroup.Done()

	file, err := os.Open(d.filepath)
	if err != nil {
		return fmt.Errorf("failed to open a file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			rerr = errors.Join(rerr, fmt.Errorf("failed to close file: %w", err))
		}
	}()

	// If not the start of the file, seek to the previous newline.
	if start != 0 {
		buffer := make([]byte, 1)

		for start > 0 {
			start--

			_, err := file.Seek(start, 0)
			if err != nil {
				return fmt.Errorf("failed to seek file: %w", err)
			}

			_, err = file.Read(buffer)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if buffer[0] == '\n' {
				start++

				break
			}
		}
	}

	_, err = file.Seek(start, 0)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	reader := bufio.NewReader(file)
	position := start

	for {
		if position >= end {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read string line: %w", err)
		}

		// Remove newline at the end of the line.
		address := line[:len(line)-1]

		ip := net.ParseIP(address)
		if ip == nil {
			log.Printf("invalid IP address: %s", address)

			continue
		}

		if bitmap.Set(ipToUint32(ip)) {
			atomic.AddInt64(uniqueIPCounter, 1)
		}

		position += int64(len(line))
	}

	return nil
}

// ipToUint32 converts an IP address to a uint32 representation.
func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}
