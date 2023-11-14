# Go Assembly and SIMD on M2

Little project to learn `Go Assembly` and activating `SIMD` operation in Go. 

## What is SIMD? Why do you need SIMD?

SIMD stands for Single Instruction, Multiple Data. It is a class of parallel computing used in processors to enhance performance by executing the same operation on multiple data points simultaneously. SIMD is a key feature in many modern CPU architectures, including x86 (with SSE and AVX extensions) and ARM (with NEON).

Traditionally, single instruction operates on the single data point. SIMD allows a single instruction to be applied simultaneously to an array of data. 

To put it simply, adding two arrays of numbers could be done in one operation, rather than looping through the arrays and adding individual pairs of numbers.

SIMD can significantly boost performance for certain types of computations, especially those involving large arrays of data that undergo the same operation.

## Does M2 chip (Apple silicon) allow SIMD?

Yes. All Apple Silicons beginning with M1 and A7, supports Advance SIMD instruction.
Since it's classified as `arm64` it uses `NEON`. Although you don't have to check for it
you can see that it's enabled by using these two commands 

```bash
systemctl -a | grep SIMD
# hw.optional.AdvSIMD: 1
# hw.optional.AdvSIMD_HPFPCvt: 1

systemctl -a | grep neon
# hw.optional.neon: 1
# hw.optional.neon_hpfp: 1
# hw.optional.neon_fp16: 1
```

## How to enable SIMD operation in go? - CGO

Using CGO gives you practical approach to use SIMD operation. In this case I'm using MacOS 
so, I'll be trying to use NEON operation.

To use a cgo, enclose your C code inside multi-line comment and write import "C".
Small but important thing is that you should put no lines between your C code comment and `import "C"`.

Example NEON Code is

```cgo
#include <arm_neon.h>

float32x4_t neon_add(float32x4_t *a, float32x4_t *b, float32x4_t *result) {
    return vaddq_f32(*a, *b);
}
import "C"
```

In Go code, the C code will be loaded like so:
```go
_ = C.neon_add((*C.float32x4_t)(unsafe.Pointer(&a[0])),
    (*C.float32x4_t)(unsafe.Pointer(&b[0])),
    (*C.float32x4_t)(unsafe.Pointer(&result[0])))
```
Some key points:

* Regarding `unsafe.Pointer(&a[0])` - It's a way to pass the entire array to the C function not just the first element
  * In C, when you pass an array to a function, what you're actually passing is a pointer to the first element of the array
  * If you have an array `arr`, then `arr` and `&arr[0]` effectively represent the same memory location: the start of the array.
  * More precisely a `slice`


## How to enable SIMD operation in go? - Assembly

SIMD operation in go typically requires Go Assembly. `asm_arm64.s` file gives brief introductions of go's assembly usage.

### file: asm_arm64.s
Function that adds two 64 bit integers. The function content is as follows.
Notice that it's slightly different from vanilla Assembly. 
```plan9_x86
    MOVD x+0(FP), R0    // Load first argument into R0
    MOVD y+8(FP), R1    // Load second argument into R1
    ADD R0, R1, R0      // Add R0 and R1, result in R0
    MOVD R0, ret+16(FP) // Store result in return value
    RET
```

Let's break down `x+0(FP)`, and `y+8(FP)` - the storing part. 
* `FP` - Frame Pointer. It is virtual register that points to the start of the function's stack frame. Stack frame is a section of the stack memory allocated for a function call, which includes space for the function's arguments and local variables.  
* Offset from `FP` is passed as `+0(FP)` or `+8(FP)`. These offsets are used to access specific arguments passed to the functions.
* `x+0(FP)` - The memory location at an offset of 0 bytes from te fram pointer `FP`. Where the first argument is stored. 
  * `x` is just for readability. 
* `y+8(FP)`. Refers the memory location at an offset of 8 bytes from `FP`. This is where the second argument to the function is stored. 
* `+8` assumes that the first argument is 8 bytes in size. 64-bit data types like `int64` or pointers on a 64-bit architecture. 

Rest of the operation is pretty straight forward. 

### How to use Go assembly functions in Go code?
First one must create function with empty function body. 
```go
//go:noescape
func Add(x, y int64) int64

func SliceOp(a []int64) int64
```
Make sure that your assembly has all the code necessary for `Add` and `SliceOp` like so.
Also, always remember to put a newline or else it will not assemble properly. 
```plan9_x86
TEXT ·Add(SB), $0
    MOVD x+0(FP), R0 // Load first argument into R0
    MOVD y+8(FP), R1 // Load second argument into R1
    ADD R0, R1, R0      // Add R0 and R1, result in R0
    MOVD R0, ret+16(FP) // Store result in return value
    RET
    
TEXT ·SliceOp(SB), $0
    MOVD a+8(FP), R0
    MOVD R0, ret+24(FP)
    RET

```
Then your code will run at your request
```go
func main() {
  // Add function in go assembly
  fmt.Println(Add(1, 2))
  
  c := []int64{1, 17, 3, 4}
  fmt.Println(SliceOp(c))
}
```

## Go assembly for ARM64

* Remember to Always end the file with newline


<b>Reference</b>
* [Apple Developer](https://developer.apple.com/documentation/kernel/1387446-sysctlbyname/determining_instruction_set_characteristics)
* [Golang ASM_arm64](https://go.dev/doc/asm#arm64)
* [Golang arm64 package](https://pkg.go.dev/cmd/internal/obj/arm64#section-sourcefiles)