package models

import (
	"encoding/hex"
	"testing"
)

func TestAddBloomFilter(t *testing.T) {
	bf := CreateBloomFilter(10, 5, 99)
	item := "Hello World"
	bf.Add(item)
	res := bf.FilterBytes()
	get := hex.EncodeToString(res)
	want := "0000000a080000000140"
	if get != want {
		t.Errorf("Expected the bloom filter to be %s but got %s", want, get)
	}
	item = "Goodbye!"
	bf.Add(item)
	res = bf.FilterBytes()
	get = hex.EncodeToString(res)
	want = "4000600a080000010940"
	if get != want {
		t.Errorf("Expected the bloom filter to be %s but got %s", want, get)
	}
}

func TestBloomFilterFilterLoad(t *testing.T) {
	bf := CreateBloomFilter(10, 5, 99)
	item := "Hello World"
	bf.Add(item)
	item = "Goodbye!"
	bf.Add(item)
	get := bf.FilterLoad(1)
	want := "0a4000600a080000010940050000006300000001"
	if get != want {
		t.Errorf("Expected the filterload mesage to be %s but got %s", want, get)
	}
}
