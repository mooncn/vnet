package obfs

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestObfsTLS_packAuthData(t *testing.T) {
	type args struct {
		clientId []byte
	}
	tests := []struct {
		name string
		otls *ObfsTLS
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test",
			otls: &ObfsTLS{
				Plain: &plain{
					ServerInfo: &serverInfo{
						Key: []byte{0x01, 0x02, 0x03, 0x04},
					},
				},
			},
			args: args{
				clientId: []byte{0x01, 0x02, 0x03, 0x04},
			},
			want: len([]byte{0x5c, 0x6c, 0x1f, 0xf0, 0x9c, 0xac, 0x23, 0x8a, 0x37, 0x55, 0x97, 0xff, 0x9, 0xcc, 0xdc, 0xa5, 0xca, 0x30, 0x24, 0xc8, 0xe3, 0x4b, 0x28, 0xde, 0x5f, 0x8c, 0xa1, 0x6d, 0xf, 0x73, 0x12, 0xa0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.otls.packAuthData(tt.args.clientId)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObfsTLS.packAuthData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hmacsha1(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				key:  []byte{0x01, 0x02, 0x03, 0x04},
				data: []byte{0x01, 0x02, 0x03, 0x04},
			},
			want: []byte{0xc7, 0x4b, 0x87, 0xee, 0x22, 0x96, 0xb8, 0xc9, 0x56, 0x33, 0x4f, 0xc5, 0x54, 0x5c, 0x63, 0x9d, 0x9c, 0x24, 0x45, 0xfc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hmacsha1(tt.args.key, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hmacsha1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleStringToByte() {
	fmt.Println(hex.EncodeToString([]byte("你")))
	//Output:
	//e4bda0
}

func TestObfsTLS_sni(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		otls *ObfsTLS
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			otls: &ObfsTLS{
				Plain: &plain{
					ServerInfo: &serverInfo{
						Key: []byte{0x01, 0x02, 0x03, 0x04},
					},
				},
			},
			args: args{
				host: "baidu.com",
			},
			want: []byte{0x0, 0x0, 0x0, 0xe, 0x0, 0xc, 0x0, 0x0, 0x9, 0x62, 0x61, 0x69, 0x64, 0x75, 0x2e, 0x63, 0x6f, 0x6d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.otls.sni(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObfsTLS.sni() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkObfsTLS_sni(b *testing.B) {
	otls := &ObfsTLS{
		Plain: &plain{
			ServerInfo: &serverInfo{
				Key: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		otls.sni("baidu.com")
	}
	b.ReportAllocs()
}

func TestObfsTLS_ClientEncode(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		otls    *ObfsTLS
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			otls: &ObfsTLS{
				Plain: &plain{
					ServerInfo: &serverInfo{
						Key:           []byte{0x01, 0x02, 0x03, 0x04},
						ProtocolParam: "",
						Host:          "0.0.0.0",
						Port:          3306,
						HeadLen:       30,
						TCPMss:        1460,
						BufferSize:    65535,
					},
				},
				ObfsAuthData: NewObfsAuthData(),
			},
			args: args{
				[]byte{0x01, 0x02, 0x03, 0x04},
			},
			want: []byte{
				0x16, 0x3, 0x1, 0x2, 0x3a, 0x1, 0x0, 0x2, 0x36, 0x3, 0x3, 0x5c, 0x79, 0xf9, 0x5f, 0xc1, 0xe9, 0xd, 0xb8, 0x58, 0xd0, 0xfc, 0xd4, 0x9f, 0xcf, 0x3e, 0x6c, 0x45, 0x2d, 0x58, 0x1e, 0xb8, 0x8d, 0xea, 0xf4, 0x75, 0x85, 0xac, 0x58, 0x4e, 0x5, 0xd5, 0x64, 0x20, 0x2b, 0xbd, 0x4f, 0xc1, 0x5a, 0xc8, 0x4a, 0x1c, 0x5, 0xa1, 0x83, 0x35, 0xee, 0x7d, 0x9d, 0xbf, 0x4c, 0xa9, 0x51, 0xa2, 0xc6, 0x16, 0x72, 0xdb, 0x17, 0x36, 0xf0, 0x9f, 0x5c, 0x1, 0xac, 0xb6, 0x0, 0x1c, 0xc0, 0x2b, 0xc0, 0x2f, 0xcc, 0xa9, 0xcc, 0xa8, 0xcc, 0x14, 0xcc, 0x13, 0xc0, 0xa, 0xc0, 0x14, 0xc0, 0x9, 0xc0,
				0x13, 0x0, 0x9c, 0x0, 0x35, 0x0, 0x2f, 0x0, 0xa, 0x1, 0x0, 0x1, 0xd1, 0xff, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x17, 0x0, 0x0, 0x0, 0x23, 0x1, 0x80,
				0x0, 0x4, 0x4, 0x34, 0x91, 0x6c, 0x6a, 0xf1, 0x7b, 0x6c, 0x47, 0x8e, 0x33, 0xc7, 0x8d, 0xff, 0xce, 0x95, 0x98, 0xc9, 0xc1, 0x6f, 0xc7, 0x1b, 0xc4, 0x1, 0xa7, 0xb7, 0x33, 0x22, 0x7c, 0xc2, 0x88, 0x13, 0x2f, 0xc9, 0x35, 0xe8, 0x74, 0x78, 0xe9, 0x57, 0x27, 0x73, 0xae, 0xcd, 0xd6, 0xd9, 0xc8, 0x5, 0x53, 0xf9, 0x57, 0xe, 0xf0, 0x27, 0x56, 0x8b, 0x5, 0x27, 0x4a, 0xd2, 0x16, 0x33, 0x71, 0xb6, 0xb7, 0x75, 0xc0, 0x8d, 0xe5, 0x1, 0x61, 0xec, 0x36, 0x89, 0xd3, 0xe5, 0x81, 0x23, 0x1b, 0x3f, 0xde, 0x1d, 0x3e, 0x1a, 0xd1, 0x69, 0x20, 0x23, 0x7d, 0x7d, 0xf6, 0xa5, 0xd6, 0x4e, 0x99, 0x36, 0x79, 0xb9, 0xdd, 0xa7, 0x51, 0x87, 0x9a, 0x74, 0xb8, 0x8c, 0xc0, 0xe3, 0x69, 0x5, 0x29, 0xfb, 0xb6, 0xad, 0xed, 0x83, 0x79, 0xcf, 0x58, 0xb8, 0x8e, 0x14, 0xe1, 0xd1, 0x7c, 0x6c, 0x90, 0x24, 0x5f, 0xd4, 0xa2, 0x9b, 0xea, 0x3a, 0x1d, 0xf, 0x48, 0xb8, 0x29, 0xd3, 0x8d, 0x98, 0x7b, 0x9d, 0xd8, 0x48, 0x3a, 0xb1, 0xdb, 0xf3, 0xc4, 0x27, 0x48, 0xf6, 0x28, 0xf4, 0xdb, 0x41,
				0x21, 0xd1, 0x81, 0x28, 0x3e, 0xf, 0xee, 0x94, 0x1a, 0x24, 0xd7, 0x7c, 0x3a, 0xdb, 0x2b, 0xd, 0xf5, 0x7b, 0x7e, 0x1b, 0x8a, 0x1d, 0x5f, 0x6d, 0x8d, 0x78, 0x76, 0x85, 0x65, 0xd9, 0x7, 0x49, 0x5e, 0xf2, 0x78, 0x18, 0x70, 0x35, 0x64, 0xb9, 0xdc, 0x82, 0x4c, 0xe2, 0xa0, 0x81, 0x93, 0xd1,
				0xe0, 0xd5, 0xfe, 0xf0, 0x41, 0xdc, 0x36, 0xed, 0x10, 0xa5, 0x72, 0x8b, 0xb, 0x9e, 0x31, 0x9, 0x72, 0x30, 0xb4, 0x14, 0x54, 0xc7, 0x5a, 0x20, 0x35, 0x8a, 0xe6, 0x98, 0x62, 0x5f, 0x31, 0x5f,
				0x76, 0xc1, 0xea, 0x1c, 0xad, 0x73, 0xce, 0x68, 0x4f, 0xff, 0x31, 0x9d, 0x97, 0x64, 0xda, 0xd1, 0x6f, 0x6e, 0x83, 0xc0, 0xfe, 0x3e, 0xac, 0x61, 0x33, 0x35, 0xca, 0x5, 0xb8, 0x14, 0xa2, 0x34, 0x1f, 0x6c, 0x2b, 0xfa, 0x51, 0x2e, 0x9f, 0x8d, 0x3e, 0xfd, 0x5c, 0xfb, 0x67, 0x11, 0x6a, 0x27, 0xdc, 0x92, 0x34, 0x14, 0x64, 0xcb, 0x4b, 0xa5, 0xa4, 0x9c, 0xd3, 0x87, 0x24, 0x55, 0x4b, 0xba, 0x49, 0x3e, 0xf0, 0x27, 0x72, 0x74, 0x83, 0x2, 0xb9, 0xbd, 0x34, 0x65, 0x25, 0xe8, 0x5d, 0xeb, 0x96, 0x0, 0xab, 0x18, 0x85, 0xe7, 0x5d, 0x53, 0x8, 0x61, 0x39, 0x9a, 0xd7, 0x70, 0x93, 0xcd, 0x95, 0xd2, 0x23, 0x58, 0x8, 0x12, 0xef, 0x1b, 0x16, 0xd2, 0x73, 0xb2, 0xe0, 0x3f, 0x23, 0x5b, 0xcb, 0x10, 0x34, 0x64, 0x38, 0xf0, 0x67, 0x9b, 0xcf, 0xa7, 0xae, 0xb, 0x97, 0xcb, 0x8f, 0xd9, 0x16, 0xd6, 0xad, 0x3, 0xb2, 0xe7, 0x6d, 0x84, 0x6d, 0x76, 0x40, 0x75, 0x8e, 0x29, 0xe8, 0xba, 0x0, 0xd, 0x0, 0x16, 0x0, 0x14, 0x6, 0x1, 0x6, 0x3, 0x5, 0x1, 0x5, 0x3, 0x4, 0x1, 0x4, 0x3, 0x3, 0x1, 0x3, 0x3, 0x2, 0x1, 0x2, 0x3, 0x0, 0x5, 0x0, 0x5, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x0, 0x0, 0x75, 0x50, 0x0, 0x0, 0x0, 0xb, 0x0, 0x2, 0x1, 0x0, 0x0, 0xa, 0x0, 0x6, 0x0, 0x4, 0x0, 0x17, 0x0, 0x18,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.otls.ClientEncode(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObfsTLS.ClientEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObfsTLS.ClientEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, MD5([]byte(password)))
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], MD5(d))
	}
	return m[:keyLen]
}

func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// ExampleOtlsClientDecode is test obfs function normal
func ExampleOtlsClientDecode() {

	tc := InitObfs()
	ts := InitObfs()

	args := []byte{0x01, 0x02, 0x03, 0x04}
	result, err := tc.ClientEncode(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, _, err = ts.ServerDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err = ts.ServerEncode([]byte{})
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, err = tc.ClientDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _ = tc.ClientEncode([]byte{})
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, _, err = ts.ServerDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %x\n", result)

	result, err = tc.ClientEncode([]byte{0x01, 0x02})
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, _, err = ts.ServerDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %x\n", result)

	result, err = ts.ServerEncode([]byte{0x01, 0x02})
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, err = tc.ClientDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %x\n", result)

	result, err = tc.ClientEncode(bytes.Repeat([]byte{0x01, 0x02}, 2048))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, _, err = ts.ServerDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %v\n", len(result))

	result, err = ts.ServerEncode(bytes.Repeat([]byte{0x01, 0x02}, 2048))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, _, err = tc.ClientDecode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("result: %v\n", len(result))

	//Output:
	//
}

func ExampleTimea() {
	fmt.Println(time.Now().Unix() & 0xFFFFFFFF)
	//Output:
}

func ExampleSlice() {
	a := "abc"
	fmt.Println(a[:1])
	//Output:
}

func InitObfs() *ObfsTLS {
	otls := &ObfsTLS{
		Plain: &plain{
			ServerInfo: &serverInfo{
				Key:           []byte{0x01, 0x02, 0x03, 0x04},
				ProtocolParam: "",
				Host:          "0.0.0.0",
				Port:          3306,
				HeadLen:       30,
				TCPMss:        1460,
				BufferSize:    65535,
			},
		},
		TLSVersion:   DEFAULT_VERSION,
		ObfsAuthData: NewObfsAuthData(),
	}
	return otls
}