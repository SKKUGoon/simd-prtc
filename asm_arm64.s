// +build arm64,noasm

// MOVD is used for moving 64bit values
// R0, R1 are general purpose registers

TEXT ·Add(SB), $0       // `$0` is the frame size.
    MOVD x+0(FP), R0    // Load first argument into R0
    MOVD y+8(FP), R1    // Load second argument into R1
    ADD R0, R1, R0      // Add R0 and R1, result in R0
    MOVD R0, ret+16(FP) // Store result in return value
    RET

TEXT ·SliceOp(SB), $0
    MOVD a+8(FP), R0
    MOVD R0, ret+24(FP)
    RET
