/**
 * This golang jit code generator are port from https://github.com/nelhage/gojit
 * 
 */
package jit

import (
	"fmt"
	"unsafe"
)

type ABI int

const (
	GoABI  ABI = iota
	CgoABI 
)

/**
 * This comment from Intel® 64 and IA-32 Architectures Software Developer’s Manual Vol. 2A. 2-1
 * 2.1 INSTRUCTION FORMAT FOR PROTECTED MODE, REAL-ADDRESS MODE, AND VIRTUAL-8086 MODE
 * The Intel 64 and IA-32 architectures instruction encodings are subsets of the format shown in Figure 2-1. Instructions consist of optional instruction prefixes (in any order), primary opcode bytes (up to three bytes), an
 * addressing-form specifier (if required) consisting of the ModR/M byte and sometimes the SIB (Scale-Index-Base)
 * byte, a displacement (if required), and an immediate data field (if required).
 * 
 * | Instruction Prefixs | Opcode | ModR/M | SIB | Displacement | Immediate |
 * 
 * - Instruction Prefixs: 
 *     Prefixes of 1 byte each (optional) <sup>1</sup>, <sup>2</sup>.
 * - Opcode:
 * 	   1-, 2-, or 3-byte opcode.
 * - ModR/M:
 * 	   1 byte (if required)
 *     
 *     7    6 5            3 2      0
 * 	   | Mod | Reg / Opcode | R / M |
 * 
 * - SIB:
 * 	   1 byte (if required)
 * 
 * 	   7      6 5     3 2     0
 * 	   | Scale | Index | base |
 * 
 * - Displacement:
 * 	   Address displacement of 1, 2, or 4 bytes or none <sup>3</sup>.
 * 
 * - Immediate:
 * 	   Immediate data of 1, 2, or 4 bytes or none <sup>3</sup>.
 * 
 * 
 * 1. The REX prefix is optional, but if used must be immediately before the opcode; see Section
 * 2.2.1, "REX Prefixes" for additional information.
 * 2. For VEX encoding information, see Section 2.3, "Intel® Advanced Vector Extensions (Intel®
 * AVX)".
 * 3. Some rare instructions can take an 8B immediate or 8B displacement.
 *
 */


/**
 * Assembler implements a simple amd64 assembler. All methods on
 * Assembler will emit code to Buf[Off:] and advances Off. Buf will
 * never be reallocated, and attempts to assemble off the end of Buf
 * will panic.
 */
type Assembler struct {
	addr *Assembler
	Buf []byte
	Off int
	ABI ABI
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (a *Assembler) copyCheck() {
	if a.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		a.addr = (*Assembler)(noescape(unsafe.Pointer(a)))
	} else if a.addr != a {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
}

func New(size int) (*Assembler, error) {
	buf, e := Alloc(size)
	if e != nil {
		return nil, e
	}
	return &Assembler{Buf: buf}, nil
}

func NewGoABI(size int) (*Assembler, error) {
	buf, e := Alloc(size)
	if e != nil {
		return nil, e
	}
	return &Assembler{Buf: buf, ABI: GoABI}, nil
}

func (a *Assembler) Dump() {
	for i := 0; i < len(a.Buf); i += 16 {
		fmt.Printf("\t")
		j := i
		for ; j < i+16 && j < len(a.Buf); j++ {
			fmt.Printf("0x%02x, ", a.Buf[j])
		}
		fmt.Printf("\n")
	}
}

func (a *Assembler) Release() {
	Release(a.Buf)
	a.addr = nil
}

func (a *Assembler) BuildTo(out interface{}) {
	switch a.ABI {
	case CgoABI:
		BuildToCgo(a.Buf, out)
	case GoABI:
		BuildTo(a.Buf, out)
	default:
		panic("bad ABI")
	}
}

func (a *Assembler) byte(b byte) {
	a.copyCheck()
	a.Buf[a.Off] = b
	a.Off++
}

func (a *Assembler) int16(i uint16) {
	a.copyCheck()
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Off += 2
}

func (a *Assembler) int32(i uint32) {
	a.copyCheck()
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Buf[a.Off+2] = byte(i >> 16)
	a.Buf[a.Off+3] = byte(i >> 24)
	a.Off += 4
}

func (a *Assembler) int64(i uint64) {
	a.copyCheck()
	a.Buf[a.Off] = byte(i & 0xFF)
	a.Buf[a.Off+1] = byte(i >> 8)
	a.Buf[a.Off+2] = byte(i >> 16)
	a.Buf[a.Off+3] = byte(i >> 24)
	a.Buf[a.Off+4] = byte(i >> 32)
	a.Buf[a.Off+5] = byte(i >> 40)
	a.Buf[a.Off+6] = byte(i >> 48)
	a.Buf[a.Off+7] = byte(i >> 56)
	a.Off += 8
}

func (a *Assembler) rel32(addr uintptr) {
	a.copyCheck()
	off := uintptr(addr) - Addr(a.Buf[a.Off:]) - 4
	if uintptr(int32(off)) != off {
		panic("call rel: target out of range")
	}
	a.int32(uint32(off))
}

/**
 * This comment from Intel® 64 and IA-32 Architectures Software Developer’s Manual Vol. 2A. 2-8
 * 2.2.1 REX Prefixes
 * REX prefixes are instruction-prefix bytes used in 64-bit mode. They do the following:
 * - Specify GPRs and SSE registers.
 * - Specify 64-bit operand size.
 * - Specify extended control registers.
 * Not all instructions require a REX prefix in 64-bit mode. A prefix is necessary only if an instruction references one
 * of the extended registers or uses a 64-bit operand. If a REX prefix is used when it has no meaning, it is ignored.
 * Only one REX prefix is allowed per instruction. If used, the REX prefix byte must immediately precede the opcode
 * byte or the escape opcode byte (0FH). When a REX prefix is used in conjunction with an instruction containing a
 * mandatory prefix, the mandatory prefix must come before the REX so the REX prefix can be immediately preceding
 * the opcode or the escape byte. For example, CVTDQ2PD with a REX prefix should have REX placed between F3 and
 * 0F E6. Other placements are ignored. The instruction-size limit of 15 bytes still applies to instructions with a REX
 * prefix. See Figure 2-3.
 */

func (a *Assembler) rex(w, r, x, b bool) {
	a.copyCheck()
	var bits byte
	if w {
		bits |= REXW
	}
	if r {
		bits |= REXR
	}
	if x {
		bits |= REXX
	}
	if b {
		bits |= REXB
	}
	if bits != 0 {
		a.byte(PFX_REX | bits)
	}
}

func (a *Assembler) rexBits(lsize, rsize byte, r, x, b bool) {
	a.copyCheck()
	if lsize != 0 && rsize != 0 && lsize != rsize {
		panic("mismatched instruction sizes")
	}
	lsize = lsize | rsize
	if lsize == 0 {
		lsize = 64
	}
	a.rex(lsize == 64, r, x, b)
}

/**
 * This comment from Intel® 64 and IA-32 Architectures Software Developer’s Manual Vol. 2A. 2-3
 * ModR/M and SIB Bytes
 * Many instructions that refer to an operand in memory have an addressing-form specifier byte (called the ModR/M
 * byte) following the primary opcode. The ModR/M byte contains three fields of information:
 * 
 * - The mod field combines with the r/m field to form 32 possible values: eight registers and 24 addressing modes.
 * - The reg/opcode field specifies either a register number or three more bits of opcode information. The purpose
 * of the reg/opcode field is specified in the primary opcode.
 * - The r/m field can specify a register as an operand or it can be combined with the mod field to encode an
 * addressing mode. Sometimes, certain combinations of the mod field and the r/m field are used to express
 * opcode information for some instructions.
 * 
 * Certain encodings of the ModR/M byte require a second addressing byte (the SIB byte). The base-plus-index and
 * scale-plus-index forms of 32-bit addressing require the SIB byte. The SIB byte includes the following fields:
 * 
 * - The scale field specifies the scale factor.
 * - The index field specifies the register number of the index register.
 * - The base field specifies the register number of the base register.
 * 
 * See Section 2.1.5 for the encodings of the ModR/M and SIB bytes.
 * 
 * - ModR/M:
 * 	1 byte (if required)
 *     
 *     7    6 5            3 2      0
 * 	| Mod | Reg / Opcode | R / M |
 */

// ModR/M bytes 
func (a *Assembler) modrm(mod, reg, rm byte) {
	a.copyCheck()
	a.byte((mod << 6) | (reg << 3) | rm)
}

// SIB, Scale-Index-Base byte
func (a *Assembler) sib(s, i, b byte) {
	a.copyCheck()
	a.byte((s << 6) | (i << 3) | b)
}
