package internal

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"forestree"
	"io"
	"log"
	"os"
)

// loadFromFile loads the forest data from the file.
func (app *App) loadFromFile(filename string) error {
	// Here, load and verify the data, decompress, and unmarshal the forest data from the file
	var loadedForest forestree.Node

	err := app.ReadChangesFromFile(filename, &loadedForest)
	if err != nil {
		return fmt.Errorf("error loading forest from file: %v", err)
	}

	// Assign the loaded forest data to the app's forest
	app.forest = &loadedForest
	return nil
}

// LoadStateFromFile loads the forest state from a file
func (app *App) LoadStateFromFile(filename string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open state file: %v", err)
	}
	defer file.Close()

	// Create a new decoder for reading the JSON data
	decoder := json.NewDecoder(file)

	// Decode into the app's forest
	if err := decoder.Decode(&app.forest); err != nil {
		return fmt.Errorf("failed to decode state file: %v", err)
	}

	log.Println("Successfully loaded forest from file.")
	return nil
}

// TODO: Encrypt this
// Function to write changes to an encrypted state file
func (app *App) WriteChangesToFile(data interface{}, filename string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create a hash of the JSON data
	hash := sha256.New()
	hash.Write(jsonData)
	hashedData := hash.Sum(nil)

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the hash to the file header
	_, err = file.Write(hashedData)
	if err != nil {
		return err
	}

	// Create a gzip writer to compress the JSON data
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// Write the JSON data to the gzip writer
	_, err = gzipWriter.Write(jsonData)
	return err
}

// ReadChangesFromFile reads the hash and compressed data from the file, decompresses and validates it.
func (app *App) ReadChangesFromFile(filename string, data interface{}) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the hash from the file header
	hash := make([]byte, sha256.Size)
	_, err = file.Read(hash)
	if err != nil {
		return err
	}

	// Create a gzip reader to decompress the data
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Decompress the data
	var decompressedData []byte
	decompressedData, err = io.ReadAll(gzipReader)
	if err != nil {
		return err
	}

	// Verify the hash
	dataHash := sha256.New()
	dataHash.Write(decompressedData)
	if !compareHashes(hash, dataHash.Sum(nil)) {
		return fmt.Errorf("data hash mismatch, file may be corrupted")
	}

	// Unmarshal the decompressed data back into the provided data structure
	err = json.Unmarshal(decompressedData, data)
	if err != nil {
		return err
	}

	return nil
}
