package debugging

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// print a slice of bytes in various formats
func BytesDump(b []byte) {

	fmt.Printf("Bytes: %d\n", len(b))
	fmt.Printf("Bits: %d\n", len(b)*8)
	fmt.Printf("Base64: %s\n", base64.StdEncoding.EncodeToString(b))
	fmt.Printf("Ascii: %c\n", b)
	fmt.Printf("Binary: %b\n", b)
	fmt.Printf("Hexadecimal: %s\n%s\n", hex.EncodeToString(b), hex.Dump(b))
	//fmt.Printf("Hexadecimal: %x\n", b)
	fmt.Printf("Raw String: %s\n", string(b[:]))
}

// print a struct in human readable form
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
