package gosamplerate

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestGetConverterName(t *testing.T) {
	name, err := GetName(SRC_LINEAR)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Linear Interpolator" {
		t.Fatal("Unexpected string")
	}
}

func TestGetConverterNameError(t *testing.T) {
	_, err := GetName(5)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "unknown samplerate converter" {
		t.Fatal("unexpected string")
	}
}

func TestGetConverterDescription(t *testing.T) {
	desc, err := GetDescription(SRC_LINEAR)
	if err != nil {
		t.Fatal(err)
	}
	if desc != "Linear interpolator, very fast, poor quality." {
		t.Fatal("Unexpected string")
	}
}

func TestGetConverterDescriptionError(t *testing.T) {
	_, err := GetDescription(5)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "unknown samplerate converter" {
		t.Fatal("unexpected string")
	}
}

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if !strings.Contains(version, "libsamplerate-") {
		t.Fatal("Unexpected string")
	}
}

func TestInitAndDestroy(t *testing.T) {
	channels := 2
	src, err := Make(SRC_SINC_FASTEST, channels, 100)
	if err != nil {
		t.Fatal(err)
	}

	chs, err := src.GetChannels()
	if err != nil {
		t.Fatal(err)
	}
	if chs != channels {
		t.Fatal("unexpected amount of channels")
	}

	err = src.Reset()
	if err != nil {
		t.Fatal(err)
	}

	err = Delete(src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidSrcObject(t *testing.T) {
	_, err := Make(5, 2, 100)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "Could not initialize samplerate converter object" {
		t.Log("unexpected Error string")
	}
}

func TestSimple(t *testing.T) {
	input := []float32{0.1, -0.5, 0.3, 0.4, 0.1}
	expectedOutput := []float32{0.1, 0.1, -0.10000001, -0.5, 0.033333343, 0.33333334, 0.4, 0.2}

	output, err := Simple(input, 1.5, 1, SRC_LINEAR)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Log("input", input)
		t.Log("output", output)
		t.Fatal("unexpected output")
	}
}

func TestSimpleLessThanOne(t *testing.T) {
	var input []float32
	for i := 0; i < 10; i++ {
		input = append(input, 0.1, -0.5, 0.3, 0.4, 0.1)
	}
	expectedOutput := []float32{0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3}

	output, err := Simple(input, 0.5, 1, SRC_LINEAR)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Log("input", input)
		t.Log("output", output)
		t.Fatal("unexpected output")
	}
}

func TestSimpleError(t *testing.T) {

	input := []float32{0.1, 0.9}
	var invalidRatio float64 = -5.3

	_, err := Simple(input, invalidRatio, 1, SRC_LINEAR)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "Error code: 6; SRC ratio outside [1/256, 256] range." {
		t.Log(err.Error())
		t.Fatal("unexpected string")
	}
}

func TestProcess(t *testing.T) {
	src, err := Make(SRC_LINEAR, 2, 100)
	if err != nil {
		t.Fatal(err)
	}

	input := []float32{0.1, -0.5, 0.2, -0.3}
	output, err := src.Process(input, 2.0, false)
	if err != nil {
		t.Fatal(err)
	}
	expOutput := []float32{0.1, -0.5, 0.1, -0.5, 0.1, -0.5, 0.15, -0.4}

	if !reflect.DeepEqual(output, expOutput) {
		t.Log("input:", input)
		t.Log("output:", output)
		t.Fatal("unexpected output")
	}

	err = Delete(src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessWithEndOfInputFlagSet(t *testing.T) {
	src, err := Make(SRC_SINC_FASTEST, 2, 100)
	if err != nil {
		t.Fatal(err)
	}

	input := []float32{0.1, -0.5, 0.2, -0.3}
	output, err := src.Process(input, 2.0, true)
	if err != nil {
		t.Fatal(err)
	}
	expOutput := []float32{0.11488709,
		-0.46334597, 0.18373828, -0.48996875, 0.1821644,
		-0.32879135, 0.10804618, -0.11150829}

	if !reflect.DeepEqual(output, expOutput) {
		t.Log("input:", input)
		t.Log("output:", output)
		t.Fatal("unexpected output")
	}

	err = Delete(src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessDataSliceBiggerThanInputBuffer(t *testing.T) {
	src, err := Make(SRC_LINEAR, 1, 100)
	if err != nil {
		t.Fatal(err)
	}

	input := make([]float32, 150)
	_, err = src.Process(input, 150.0, true)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "data slice is larger than buffer" {
		t.Log("unexpected Error string")
	}
}

func TestProcessErrorWithInvalidRatio(t *testing.T) {
	src, err := Make(SRC_LINEAR, 1, 100)
	if err != nil {
		t.Fatal(err)
	}

	input := make([]float32, 100)
	_, err = src.Process(input, -5, true)
	if err == nil {
		t.Fatal("expected Error")
	}
	if err.Error() != "Error code: 6; SRC ratio outside [1/256, 256] range." {
		t.Log(err.Error())
		t.Log("unexpected Error string")
	}
}

func TestGetChannels(t *testing.T) {
	channels := 2
	src, err := Make(SRC_SINC_FASTEST, channels, 100)
	if err != nil {
		t.Fatal(err)
	}
	chLength, err := src.GetChannels()
	if err != nil {
		t.Fatal(err)
	} else if chLength != channels {
		t.Fatal("unexpected channel length")
	}
}

func TestSetRatio(t *testing.T) {
	src, err := Make(SRC_LINEAR, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if err = src.SetRatio(25.0); err != nil {
		t.Fatal("unexpected result; should be valid conversion rate")
	}
}

func TestSetRatioInvalid(t *testing.T) {
	src, err := Make(SRC_LINEAR, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	err = src.SetRatio(-5)
	if err == nil {
		t.Fatal("expected Error")
	}
}

func TestIsValidRatio(t *testing.T) {
	if !IsValidRatio(5) {
		t.Fatal("unexpected result; should be valid")
	}

	if IsValidRatio(-1) {
		t.Fatal("unexpected result; should be invalid")
	}

	if IsValidRatio(257) {
		t.Fatal("unexpected result; should be invalid")
	}
}

func TestErrors(t *testing.T) {
	channels := 2
	src, err := Make(SRC_SINC_FASTEST, channels, 100)
	if err != nil {
		t.Fatal(err)
	}

	errNo := src.ErrorNo()
	if errNo != 0 {
		t.Fatal("unexpected error number")
	}

	errString := Error(0)
	if errString != "No error." {
		t.Fatal("unexpected Error string")
	}

	err = Delete(src)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFloatToInt16Array(t *testing.T) {
	input := []float32{-1.0, 1.0}
	output := make([]int16, 2)
	if err := FloatToInt16Array(input, output); err != nil {
		t.Fatalf("did not expect failure: %v", err)
	}
	if output[0] != -32768 {
		t.Fatalf("did not expect -1.0 to map to %v", output[0])
	}
	if output[1] != 32767 {
		t.Fatalf("did not expect  1.0 to map to %v", output[1])
	}
}

func TestInt16ToFloatArray(t *testing.T) {
	input := []int16{-32768, 32767}
	output := make([]float32, 2)
	if err := Int16ToFloatArray(input, output); err != nil {
		t.Fatalf("did not expect failure: %v", err)
	}
	if math.Abs(float64(output[0] - -1.0)) > 0.0001 {
		t.Fatalf("did not expect -32768 to map to %v", output[0])
	}
	if math.Abs(float64(output[1]-1.0)) > 0.0001 {
		t.Fatalf("did not expect  32767 to map to %v", output[1])
	}
}

func TestInt16ByteToFloatArray(t *testing.T) {
	input := []byte{0x00, 0x80, 0xFF, 0x7F}
	output := make([]float32, 2)
	if err := Int16ByteToFloatArray(input, output); err != nil {
		t.Fatalf("did not expect failure: %v", err)
	}
	if math.Abs(float64(output[0] - -1.0)) > 0.0001 {
		t.Fatalf("did not expect -32768 to map to %v", output[0])
	}
	if math.Abs(float64(output[1]-1.0)) > 0.0001 {
		t.Fatalf("did not expect  32767 to map to %v", output[1])
	}
}
