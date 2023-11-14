package main

/*
#include <arm_neon.h>

float32x4_t neon_add(float32x4_t *a, float32x4_t *b, float32x4_t *result) {
    return vaddq_f32(*a, *b);
}

void vector_add_neon(float *a, float *b, float *result, size_t n) {
	// Process in chunks of 4 floats, since each NEON register can hold 4 floats
	for (size_t i = 0; i < n; i += 4) {
		float32x4_t va = vld1q_f32(a + i);
		float32x4_t vb = vld1q_f32(b + i);
		float32x4_t vres = vaddq_f32(va, vb);
		vst1q_f32(result + i, vres); // Store the result
	}
}
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	_ "net/http/pprof"
)

func noNeon(s1, s2, result [4]float32) [4]float32 {
	for i := 0; i < 4; i++ {
		result[i] = s1[i] + s2[i]
	}
	return result
}

func noNeonBig(s1, s2, result []float32, length int) {
	for i := 0; i < length; i++ {
		result[i] = s1[i] + s2[i]
	}
}

func noNeonBigGoRoutine(s1, s2, result []float32, length int) {

}

func benchmarkBigSize() {
	const n = 10000000
	a := make([]float32, n)
	b := make([]float32, n)
	result := make([]float32, n)

	// Initialize a and b with some values
	for i := range a {
		a[i] = float32(i)
		b[i] = float32(i)
	}

	bigNeon := time.Now()
	C.vector_add_neon(
		(*C.float)(unsafe.Pointer(&a[0])),
		(*C.float)(unsafe.Pointer(&b[0])),
		(*C.float)(unsafe.Pointer(&result[0])),
		C.size_t(n),
	)
	bigNeonEnd := time.Now()

	bigNoNeon := time.Now()
	noNeonBig(a, b, result, len(a))
	bigNoNeonEnd := time.Now()

	fmt.Println("With SIMD", bigNeonEnd.Sub(bigNeon), "Without SIMD", bigNoNeonEnd.Sub(bigNoNeon))
}

func benchmarkSmallSize() {
	var a, b, result [4]float32
	a = [4]float32{1, 2, 3, 4}
	b = [4]float32{2, 4, 6, 8}

	singleNeon := time.Now()
	c := C.neon_add((*C.float32x4_t)(unsafe.Pointer(&a[0])),
		(*C.float32x4_t)(unsafe.Pointer(&b[0])),
		(*C.float32x4_t)(unsafe.Pointer(&result[0])))
	fmt.Println("neon result", c)
	singleNeonEnd := time.Now()

	singleNoNeon := time.Now()
	noNeon(a, b, result)
	singleNoNeonEnd := time.Now()

	fmt.Println("With SIMD", singleNeonEnd.Sub(singleNeon), "Without SIMD", singleNoNeonEnd.Sub(singleNoNeon))
}

func main() {
	// Benchmark NEON small size (1 work with 4 elements vectors)
	fmt.Println("Small")
	benchmarkSmallSize()

	// Benchmark NEON big size (1 work with 10000 elements vectors
	fmt.Println("Big")
	benchmarkBigSize()
}
