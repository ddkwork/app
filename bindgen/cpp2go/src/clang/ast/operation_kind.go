package ast

// CastKind - The kind of operation required for a conversion.
type CastKind string

const (
	// Dependent - A conversion which cannot yet be analyzed because
	// either the expression or target type is dependent.  These are
	// created only for explicit casts; dependent ASTs aren't required
	// to even approximately type-check.
	//   (T*) malloc(sizeof(T))
	//   reinterpret_cast<intptr_t>(A<T>::alloc());
	Dependent CastKind = "Dependent"

	// BitCast - A conversion which causes a bit pattern of one type
	// to be reinterpreted as a bit pattern of another type.  Generally
	// the operands must have equivalent size and unrelated types.
	///
	// The pointer conversion char* -> int* is a bitcast.  A conversion
	// from any pointer type to a C pointer type is a bitcast unless
	// it's actually BaseToDerived or DerivedToBase.  A conversion to a
	// block pointer or ObjC pointer type is a bitcast only if the
	// operand has the same type kind; otherwise, it's one of the
	// specialized casts below.
	///
	// Vector coercions are bitcasts.
	BitCast CastKind = "BitCast"

	// LValueBitCast - A conversion which reinterprets the address of
	// an l-value as an l-value of a different kind.  Used for
	// reinterpret_casts of l-value expressions to reference types.
	//    bool b; reinterpret_cast<char&>(b) = 'a';
	LValueBitCast CastKind = "LValueBitCast"

	// LValueToRValueBitCast - A conversion that causes us to reinterpret the
	// object representation of an lvalue as an rvalue. Created by
	// __builtin_bit_cast.
	LValueToRValueBitCast CastKind = "LValueToRValueBitCast"

	// LValueToRValue - A conversion which causes the extraction of
	// an r-value from the operand gl-value.  The result of an r-value
	// conversion is always unqualified.
	LValueToRValue CastKind = "LValueToRValue"

	// NoOp - A conversion which does not affect the type other than
	// (possibly) adding qualifiers.
	//   int    -> int
	//   char** -> const char * const *
	NoOp CastKind = "NoOp"

	// BaseToDerived - A conversion from a C++ class pointer/reference
	// to a derived class pointer/reference.
	//   B *b = static_cast<B*>(a);
	BaseToDerived CastKind = "BaseToDerived"

	// DerivedToBase - A conversion from a C++ class pointer
	// to a base class pointer.
	//   A *a = new B();
	DerivedToBase CastKind = "DerivedToBase"

	// UncheckedDerivedToBase - A conversion from a C++ class
	// pointer/reference to a base class that can assume that the
	// derived pointer is not null.
	//   const A &a = B();
	//   b->method_from_a();
	UncheckedDerivedToBase CastKind = "UncheckedDerivedToBase"

	// Dynamic - A C++ dynamic_cast.
	Dynamic CastKind = "Dynamic"

	// ToUnion - The GCC cast-to-union extension.
	//   int   -> union { int x; float y; }
	//   float -> union { int x; float y; }
	ToUnion CastKind = "ToUnion"

	// ArrayToPointerDecay - Array to pointer decay.
	//   int[10] -> int*
	//   char[5][6] -> char(*)[6]
	ArrayToPointerDecay CastKind = "ArrayToPointerDecay"

	// FunctionToPointerDecay - Function to pointer decay.
	//   void(int) -> void(*)(int)
	FunctionToPointerDecay CastKind = "FunctionToPointerDecay"

	// NullToPointer - Null pointer constant to pointer, ObjC
	// pointer, or block pointer.
	//   (void*) 0
	//   void (^block)() = 0;
	NullToPointer CastKind = "NullToPointer"

	// NullToMemberPointer - Null pointer constant to member pointer.
	//   int A::*mptr = 0;
	//   int (A::*fptr)(int) = nullptr;
	NullToMemberPointer CastKind = "NullToMemberPointer"

	// BaseToDerivedMemberPointer - Member pointer in base class to
	// member pointer in derived class.
	//   int B::*mptr = &A::member;
	BaseToDerivedMemberPointer CastKind = "BaseToDerivedMemberPointer"

	// DerivedToBaseMemberPointer - Member pointer in derived class to
	// member pointer in base class.
	//   int A::*mptr = static_cast<int A::*>(&B::member);
	DerivedToBaseMemberPointer CastKind = "DerivedToBaseMemberPointer"

	// MemberPointerToBoolean - Member pointer to boolean.  A check
	// against the null member pointer.
	MemberPointerToBoolean CastKind = "MemberPointerToBoolean"

	// ReinterpretMemberPointer - Reinterpret a member pointer as a
	// different kind of member pointer.  C++ forbids this from
	// crossing between function and object types, but otherwise does
	// not restrict it.  However, the only operation that is permitted
	// on a "punned" member pointer is casting it back to the original
	// type, which is required to be a lossless operation (although
	// many ABIs do not guarantee this on all possible intermediate types).
	ReinterpretMemberPointer CastKind = "ReinterpretMemberPointer"

	// UserDefinedConversion - Conversion using a user defined type
	// conversion function.
	//    struct A { operator int(); }; int i = int(A());
	UserDefinedConversion CastKind = "UserDefinedConversion"

	// ConstructorConversion - Conversion by constructor.
	//    struct A { A(int); }; A a = A(10);
	ConstructorConversion CastKind = "ConstructorConversion"

	// IntegralToPointer - Integral to pointer.  A special kind of
	// reinterpreting conversion.  Applies to normal, ObjC, and block
	// pointers.
	//    (char*) 0x1001aab0
	//    reinterpret_cast<int*>(0)
	IntegralToPointer CastKind = "IntegralToPointer"

	// PointerToIntegral - Pointer to integral.  A special kind of
	// reinterpreting conversion.  Applies to normal, ObjC, and block
	// pointers.
	//    (intptr_t) "help!"
	PointerToIntegral CastKind = "PointerToIntegral"

	// PointerToBoolean - Pointer to boolean conversion.  A check
	// against null.  Applies to normal, ObjC, and block pointers.
	PointerToBoolean CastKind = "PointerToBoolean"

	// ToVoid - Cast to void, discarding the computed value.
	//    (void) malloc(2048)
	ToVoid CastKind = "ToVoid"

	// VectorSplat - A conversion from an arithmetic type to a
	// vector of that element type.  Fills all elements ("splats") with
	// the source value.
	//    __attribute__((ext_vector_type(4))) int v = 5;
	VectorSplat CastKind = "VectorSplat"

	// IntegralCast - A cast between integral types (other than to
	// boolean).  Variously a bitcast, a truncation, a sign-extension,
	// or a zero-extension.
	//    long l = 5;
	//    (unsigned) i
	IntegralCast CastKind = "IntegralCast"

	// IntegralToBoolean - Integral to boolean.  A check against zero.
	//    (bool) i
	IntegralToBoolean CastKind = "IntegralToBoolean"

	// IntegralToFloating - Integral to floating point.
	//    float f = i;
	IntegralToFloating CastKind = "IntegralToFloating"

	// FixedPointCast - Fixed point to fixed point.
	//    (_Accum) 0.5r
	FixedPointCast CastKind = "FixedPointCast"

	// FixedPointToIntegral - Fixed point to integral.
	//    (int) 2.0k
	FixedPointToIntegral CastKind = "FixedPointToIntegral"

	// IntegralToFixedPoint - Integral to a fixed point.
	//    (_Accum) 2
	IntegralToFixedPoint CastKind = "IntegralToFixedPoint"

	// FixedPointToBoolean - Fixed point to boolean.
	//    (bool) 0.5r
	FixedPointToBoolean CastKind = "FixedPointToBoolean"

	// FloatingToIntegral - Floating point to integral.  Rounds
	// towards zero, discarding any fractional component.
	//    (int) f
	FloatingToIntegral CastKind = "FloatingToIntegral"

	// FloatingToBoolean - Floating point to boolean.
	//    (bool) f
	FloatingToBoolean CastKind = "FloatingToBoolean"

	// BooleanToSignedIntegral - Convert a boolean to -1 or 0 for true and
	// false, respectively.
	BooleanToSignedIntegral CastKind = "BooleanToSignedIntegral"

	// FloatingCast - Casting between floating types of different size.
	//    (double) f
	//    (float) ld
	FloatingCast CastKind = "FloatingCast"

	// CPointerToObjCPointerCast - Casting a C pointer kind to an
	// Objective-C pointer.
	CPointerToObjCPointerCast CastKind = "CPointerToObjCPointerCast"

	// BlockPointerToObjCPointerCast - Casting a block pointer to an
	// ObjC pointer.
	BlockPointerToObjCPointerCast CastKind = "BlockPointerToObjCPointerCast"

	// AnyPointerToBlockPointerCast - Casting any non-block pointer
	// to a block pointer.  Block-to-block casts are bitcasts.
	AnyPointerToBlockPointerCast CastKind = "AnyPointerToBlockPointerCast"

	// ObjCObjectLValueCast Converting between two Objective-C object types, which
	// can occur when performing reference binding to an Objective-C
	// object.
	ObjCObjectLValueCast CastKind = "ObjCObjectLValueCast"

	// FloatingRealToComplex A conversion of a floating point real to a floating point
	// complex of the original type.  Injects the value as the real
	// component with a zero imaginary component.
	//   float -> _Complex float
	FloatingRealToComplex CastKind = "FloatingRealToComplex"

	// FloatingComplexToReal Converts a floating point complex to floating point real
	// of the source's element type.  Just discards the imaginary
	// component.
	//   _Complex long double -> long double
	FloatingComplexToReal CastKind = "FloatingComplexToReal"

	// FloatingComplexToBoolean Converts a floating point complex to bool by comparing
	// against 0+0i.
	FloatingComplexToBoolean CastKind = "FloatingComplexToBoolean"

	// FloatingComplexCast Converts between different floating point complex types.
	//   _Complex float -> _Complex double
	FloatingComplexCast CastKind = "FloatingComplexCast"

	// FloatingComplexToIntegralComplex Converts from a floating complex to an integral complex.
	//   _Complex float -> _Complex int
	FloatingComplexToIntegralComplex CastKind = "FloatingComplexToIntegralComplex"

	// IntegralRealToComplex Converts from an integral real to an integral complex
	// whose element type matches the source.  Injects the value as
	// the real component with a zero imaginary component.
	//   long -> _Complex long
	IntegralRealToComplex CastKind = "IntegralRealToComplex"

	// IntegralComplexToReal Converts an integral complex to an integral real of the
	// source's element type by discarding the imaginary component.
	//   _Complex short -> short
	IntegralComplexToReal CastKind = "IntegralComplexToReal"

	// IntegralComplexToBoolean Converts an integral complex to bool by comparing against
	// 0+0i.
	IntegralComplexToBoolean CastKind = "IntegralComplexToBoolean"

	// IntegralComplexCast Converts between different integral complex types.
	//   _Complex char -> _Complex long long
	//   _Complex unsigned int -> _Complex signed int
	IntegralComplexCast CastKind = "IntegralComplexCast"

	// IntegralComplexToFloatingComplex Converts from an integral complex to a floating complex.
	//   _Complex unsigned -> _Complex float
	IntegralComplexToFloatingComplex CastKind = "IntegralComplexToFloatingComplex"

	// ARCProduceObject [ARC] Produces a retainable object pointer so that it may
	// be consumed, e.g. by being passed to a consuming parameter.
	// Calls objc_retain.
	ARCProduceObject CastKind = "ARCProduceObject"

	// ARCConsumeObject [ARC] Consumes a retainable object pointer that has just
	// been produced, e.g. as the return value of a retaining call.
	// Enters a cleanup to call objc_release at some indefinite time.
	ARCConsumeObject CastKind = "ARCConsumeObject"

	// ARCReclaimReturnedObject [ARC] Reclaim a retainable object pointer object that may
	// have been produced and autoreleased as part of a function return
	// sequence.
	ARCReclaimReturnedObject CastKind = "ARCReclaimReturnedObject"

	// ARCExtendBlockObject [ARC] Causes a value of block type to be copied to the
	// heap, if it is not already there.  A number of other operations
	// in ARC cause blocks to be copied; this is for cases where that
	// would not otherwise be guaranteed, such as when casting to a
	// non-block pointer type.
	ARCExtendBlockObject CastKind = "ARCExtendBlockObject"

	// AtomicToNonAtomic Converts from _Atomic(T) to T.
	AtomicToNonAtomic CastKind = "AtomicToNonAtomic"
	// NonAtomicToAtomic Converts from T to _Atomic(T).
	NonAtomicToAtomic CastKind = "NonAtomicToAtomic"

	// CopyAndAutoreleaseBlockObject Causes a block literal to by copied to the heap and then
	// autoreleased.
	// This particular cast kind is used for the conversion from a C++11
	// lambda expression to a block pointer.
	CopyAndAutoreleaseBlockObject CastKind = "CopyAndAutoreleaseBlockObject"

	// BuiltinFnToFnPtr Convert a builtin function to a function pointer; only allowed in the
	// callee of a call expression.
	BuiltinFnToFnPtr CastKind = "BuiltinFnToFnPtr"

	// ZeroToOCLOpaqueType Convert a zero value for OpenCL opaque types initialization (event_t,
	// queue_t, etc.)
	ZeroToOCLOpaqueType CastKind = "ZeroToOCLOpaqueType"

	// AddressSpaceConversion Convert a pointer to a different address space.
	AddressSpaceConversion CastKind = "AddressSpaceConversion"

	// IntToOCLSampler Convert an integer initializer to an OpenCL sampler.
	IntToOCLSampler CastKind = "IntToOCLSampler"
)

type OpCode string

// ===- Binary Operations  -------------------------------------------------===//
// Operators listed in order of precedence.
// Note that additions to this should also update the StmtVisitor class and
// BinaryOperator::getOverloadedOperator.
const (
	// [C++ 5.5] Pointer-to-member operators.
	PtrMemD OpCode = ".*"
	PtrMemI OpCode = "->*"
	// [C99 6.5.5] Multiplicative operators.
	Mul OpCode = "*"
	Div OpCode = "/"
	Rem OpCode = "%"
	// [C99 6.5.6] Additive operators.
	Add OpCode = "+"
	Sub OpCode = "-"
	// [C99 6.5.7] Bitwise shift operators.
	Shl OpCode = "<<"
	Shr OpCode = ">>"
	// C++20 [expr.spaceship] Three-way comparison operator.
	Cmp OpCode = "<=>"
	// [C99 6.5.8] Relational operators.
	LT OpCode = "<"
	GT OpCode = ">"
	LE OpCode = "<="
	GE OpCode = ">="
	// [C99 6.5.9] Equality operators.
	EQ OpCode = "=="
	NE OpCode = "!="
	// [C99 6.5.10] Bitwise AND operator.
	And OpCode = "&"
	// [C99 6.5.11] Bitwise XOR operator.
	Xor OpCode = "^"
	// [C99 6.5.12] Bitwise OR operator.
	Or OpCode = "|"
	// [C99 6.5.13] Logical AND operator.
	LAnd OpCode = "&&"
	// [C99 6.5.14] Logical OR operator.
	LOr OpCode = "||"
	// [C99 6.5.16] Assignment operators.
	Assign    OpCode = "="
	MulAssign OpCode = "*="
	DivAssign OpCode = "/="
	RemAssign OpCode = "%="
	AddAssign OpCode = "+="
	SubAssign OpCode = "-="
	ShlAssign OpCode = "<<="
	ShrAssign OpCode = ">>="
	AndAssign OpCode = "&="
	XorAssign OpCode = "^="
	OrAssign  OpCode = "|="
	// [C99 6.5.17] Comma operator.
	Comma OpCode = ","
)

// ===- Unary Operations ---------------------------------------------------===//
// Note that additions to this should also update the StmtVisitor class and
// UnaryOperator::getOverloadedOperator.
const (
	// [C99 6.5.2.4] Postfix increment and decrement
	PostInc OpCode = "++"
	PostDec OpCode = "--"
	// [C99 6.5.3.1] Prefix increment and decrement
	PreInc OpCode = "++"
	PreDec OpCode = "--"
	// [C99 6.5.3.2] Address and indirection
	AddrOf OpCode = "&"
	Deref  OpCode = "*"
	// [C99 6.5.3.3] Unary arithmetic
	Plus  OpCode = "+"
	Minus OpCode = "-"
	Not   OpCode = "~"
	LNot  OpCode = "!"
	// "__real expr"/"__imag expr" Extension.
	Real OpCode = "__real"
	Imag OpCode = "__imag"
	// __extension__ marker.
	Extension OpCode = "__extension__"
	// [C++ Coroutines] co_await operator
	Coawait OpCode = "co_await"
)
