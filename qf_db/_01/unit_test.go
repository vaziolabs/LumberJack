package _01

import (
	"testing"
	. "fmt"
)

func TestUnit(t *testing.T) {
	u := NewUnit("test", 123)
	if u.Name != "test" {
		t.Errorf("Expected name to be 'test', got %s", u.Name)
	}
	if u.Value != 123 {
		t.Errorf("Expected value to be 123, got %v", u.Value)
	}
	Println(" > UnitTest")
	Println(u)
	Println(" ")
}

func TestHash(t *testing.T) {
	u := NewUnit("test", 123)
	u.CalculateHash()
	if u.Hash == "" {
		t.Error("Expected hash to be calculated")
	}
	Println(" > HashTest")
	Println(u)
	Println(" ")
}

func TestCompress(t *testing.T) {
	u := NewUnit("test", 123)
	u.CalculateHash()
	compressed, err := u.Compress()

	if err != nil {
		t.Error("Compression error:", err)
	}

	Println(" > CompressTest")
	Println(u)
	Println(" ")
	Println("Compressed size:", len(compressed))
	Println(" ")
	Println("Compressed:", compressed)
	Println(" ")
}