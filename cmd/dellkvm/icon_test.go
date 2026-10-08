//go:build linux || windows

package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"testing"
)

func TestTrayIconAssets(t *testing.T) {
	config, err := png.DecodeConfig(bytes.NewReader(trayIconPNG))
	if err != nil || config.Width != 32 || config.Height != 32 {
		t.Fatalf("tray PNG = %dx%d, %v", config.Width, config.Height, err)
	}
	if len(trayIconICO) < 6 || binary.LittleEndian.Uint16(trayIconICO[2:4]) != 1 {
		t.Fatal("invalid tray ICO header")
	}
	count := int(binary.LittleEndian.Uint16(trayIconICO[4:6]))
	if count != 5 || len(trayIconICO) < 6+16*count {
		t.Fatalf("tray ICO contains %d sizes", count)
	}
	for i := range count {
		entry := trayIconICO[6+i*16 : 6+(i+1)*16]
		size := int(binary.LittleEndian.Uint32(entry[8:12]))
		offset := int(binary.LittleEndian.Uint32(entry[12:16]))
		if offset < 6+16*count || size <= 0 || offset+size > len(trayIconICO) {
			t.Fatalf("ICO entry %d has invalid offset or size", i)
		}
		if _, err := png.DecodeConfig(bytes.NewReader(trayIconICO[offset : offset+size])); err != nil {
			t.Fatalf("ICO entry %d is not a PNG: %v", i, err)
		}
	}
}
