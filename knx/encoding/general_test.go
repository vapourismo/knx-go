package encoding

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"
)

func TestReadAll(t *testing.T) {
	data := []byte{1, 3, 3, 7}
	reader := bytes.NewReader(data)

	n, result := ReadAll(reader)

	if n != int64(len(data)) {
		t.Error("Length mismatch:", n)
	}

	if !bytes.Equal(data, result) {
		t.Errorf("Result mismatch: %v != %v", data, result)
	}

	shData := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	shResult := (*reflect.SliceHeader)(unsafe.Pointer(&result))

	if shData.Data == shResult.Data {
		t.Error("Underlying storage is the same")
	}
}
