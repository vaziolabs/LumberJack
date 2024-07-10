package _01

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/gob"
	"engine"
	"fmt"
)

type Unit struct {
	Name  string
	Value interface{}
	Hash  string
	Size  int
}

func (u *Unit) Log() {
	engine.Log(engine.DebugLevel, "Unit:")
	engine.Log(engine.DebugLevel, "\tName: %s", u.Name)
	engine.Log(engine.DebugLevel, "\tValue: %v", u.Value)
	engine.Log(engine.DebugLevel, "\tHash: %s", u.Hash)
	engine.Log(engine.DebugLevel, "\tSize: %d", u.Size)
}

func NewUnit(name string, value interface{}) *Unit {
	return &Unit{
		Name:  name,
		Value: value,
	}
}

func (u *Unit) CalculateHash() {
	// Calculate hash based on the unit's name and value
	NameValueHash := sha256.New()
	NameValueHash.Write([]byte(u.Name))
	NameValueHash.Write([]byte(fmt.Sprintf("%v", u.Value)))
	u.Hash = fmt.Sprintf("%x", NameValueHash.Sum(nil))
}

func (u *Unit) Compress() ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if err := gob.NewEncoder(zw).Encode(u); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	u.Size = len(buf.Bytes())
	return buf.Bytes(), nil
}

func DecompressUnit(compressed []byte) (*Unit, error) {
	buf := bytes.NewBuffer(compressed)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	var u Unit
	if err := gob.NewDecoder(zr).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

//
