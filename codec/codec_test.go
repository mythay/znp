package codec

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecoder_decodeHeader(t *testing.T) {
	tests := []struct {
		name    string
		dec     *Decoder
		want    *Cmd
		want1   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"wrong EOF header", NewDecoder(bytes.NewBuffer([]byte{})), nil, nil, true},
		{"wrong SOF header, need to consume at least one byte to avoid dead loop", NewDecoder(bytes.NewBuffer([]byte{0xee})), nil, nil, true},
		{"wrong FCS", NewDecoder(bytes.NewBuffer([]byte{0xfe, 00, 00, 00, 11, 00})), nil, nil, true},
		{"too short, no HDR", NewDecoder(bytes.NewBuffer([]byte{0xfe, 00})), nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.dec.decodeHeader()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.decodeHeader() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Decoder.decodeHeader() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		dec     *Decoder
		want    *Cmd
		want1   interface{}
		wantErr bool
	}{
		// TODO: Add test cases.

		{" simple nil payload", NewDecoder(bytes.NewBuffer(makeRawData([]byte{MT_RPC_SYS_SYS | MT_RPC_CMD_SREQ, MT_SYS_PING}))), &Cmd{MT_RPC_SYS_SYS | MT_RPC_CMD_SREQ, MT_SYS_PING}, nil, false},
		{"payload struct, not match", NewDecoder(bytes.NewBuffer(makeRawData([]byte{MT_RPC_SYS_SYS | MT_RPC_CMD_SRSP, MT_SYS_PING, 1}))), &Cmd{MT_RPC_SYS_SYS | MT_RPC_CMD_SRSP, MT_SYS_PING}, nil, true},
		{"payload struct", NewDecoder(bytes.NewBuffer(makeRawData([]byte{MT_RPC_SYS_SYS | MT_RPC_CMD_SRSP, MT_SYS_PING, 0, 1}))), &Cmd{MT_RPC_SYS_SYS | MT_RPC_CMD_SRSP, MT_SYS_PING}, &PingSrspFormat{0x100}, false},
		{"payload struct has lenght tag", NewDecoder(bytes.NewBuffer(makeRawData([]byte{MT_RPC_SYS_SYS | MT_RPC_CMD_SREQ, MT_SYS_RAM_WRITE, 11, 00, 2, 2, 3, 4}))), &Cmd{MT_RPC_SYS_SYS | MT_RPC_CMD_SREQ, MT_SYS_RAM_WRITE}, &RamWriteFormat{11, []uint8{2, 3}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.dec.Decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Decode() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Decoder.Decode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}