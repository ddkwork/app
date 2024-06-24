package ast

import (
	"strconv"

	"github.com/ddkwork/golibrary/mylog"
)

type IncludedFrom struct {
	File         string        `json:"file"`
	IncludedFrom *IncludedFrom `json:"includedFrom,omitempty"`
}

type Loc struct {
	Offset       int64         `json:"offset,omitempty"`
	File         string        `json:"file,omitempty"`
	Line         int           `json:"line,omitempty"`
	PresumedLine int           `json:"presumedLine,omitempty"`
	Col          int           `json:"col,omitempty"`
	TokLen       int           `json:"tokLen,omitempty"`
	IncludedFrom *IncludedFrom `json:"includedFrom,omitempty"`
}

type Pos struct {
	*Loc
	SpellingLoc         *Loc `json:"spellingLoc,omitempty"`
	ExpansionLoc        *Loc `json:"expansionLoc,omitempty"`
	IsMacroArgExpansion bool `json:"isMacroArgExpansion,omitempty"`
}

type Range struct {
	Begin Pos `json:"begin"`
	End   Pos `json:"end"`
}

// -----------------------------------------------------------------------------

type BareDecl struct {
	ID   ID     `json:"id"`
	Kind Kind   `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
	Type *Type  `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------

type ID string

func (s ID) ToInt() uint64 {
	if s == "" {
		return 0
	}
	num := mylog.Check2(strconv.ParseUint(string(s), 0, 64))

	return num
}

type ValueCategory string

const (
	RValue ValueCategory = "rvalue"
	XValue ValueCategory = "xvalue"
	LValue ValueCategory = "lvalue"
)

type CallConv string

const (
	CDecl             CallConv = "cdecl"
	StdCall           CallConv = "stdcall"
	FastCall          CallConv = "fastcall"
	ThisCall          CallConv = "thiscall"
	Pascal            CallConv = "pascal"
	VectorCall        CallConv = "vectorcall"
	Win64             CallConv = "ms_abi"
	SysV64            CallConv = "sysv_abi"
	RegCall           CallConv = "regcall"
	AAPCS             CallConv = "aapcs"
	AAPCS_VFP         CallConv = "aapcs-vfp"
	AArch64VectorCall CallConv = "aarch64_vector_pcs"
	IntelOclBicc      CallConv = "intel_ocl_bicc"
	SpirFunction      CallConv = "spir_function"
	OpenCLKernel      CallConv = "opencl_kernel"
	Swift             CallConv = "swiftcall"
	PreserveMost      CallConv = "preserve_most"
	PreserveAll       CallConv = "preserve_all"
)

type AssociationKind string

const (
	Case    AssociationKind = "case"
	Default AssociationKind = "default"
)

type StorageClass string

const (
	Static StorageClass = "static"
	Extern StorageClass = "extern"
)

type Type struct {
	// QualType can be:
	//   unsigned int
	//   struct ConstantString
	//   volatile uint32_t
	//   int (*)(void *, int, char **, char **)
	//   int (*)(const char *, ...)
	//   int (*)(void)
	//   const char *restrict
	//   const char [7]
	//   char *
	//   void
	//   ...
	DesugaredQualType string `json:"desugaredQualType,omitempty"`
	QualType          string `json:"qualType"`
	TypeAliasDeclID   ID     `json:"typeAliasDeclId,omitempty"`
}

type CastPath struct {
	Name      string `json:"name"`
	IsVirtual bool   `json:"isVirtual,omitempty"`
}

type DefinitionCtor struct {
	// DefaultCtor
	DefaultedIsConstexpr bool `json:"defaultedIsConstexpr,omitempty"`
	// CopyCtor, MoveCtor, Dtor
	DefaultedIsDeleted bool `json:"defaultedIsDeleted,omitempty"`
	// DefaultCtor, MoveCtor, MoveAssign
	Exists bool `json:"exists,omitempty"`
	// CopyCtor, CopyAssign
	HasConstParam bool `json:"hasConstParam,omitempty"`
	// CopyCtor, CopyAssign
	ImplicitHasConstParam bool `json:"implicitHasConstParam,omitempty"`
	// Dtor
	Irrelevant bool `json:"irrelevant,omitempty"`
	// DefaultCtor
	IsConstexpr bool `json:"isConstexpr,omitempty"`
	// DefaultCtor, CopyCtor, MoveCtor, CopyAssign, MoveAssign, Dtor
	NeedsImplicit bool `json:"needsImplicit,omitempty"`
	// CopyCtor, MoveCtor, CopyAssign, MoveAssign, Dtor
	NeedsOverloadResolution bool `json:"needsOverloadResolution,omitempty"`
	// DefaultCtor, CopyCtor, MoveCtor, CopyAssign, MoveAssign, Dtor
	NonTrivial bool `json:"nonTrivial,omitempty"`
	// CopyCtor, MoveCtor, MoveAssign, Dtor
	Simple bool `json:"simple,omitempty"`
	// DefaultCtor, CopyCtor, MoveCtor, CopyAssign, MoveAssign, Dtor
	Trivial bool `json:"trivial,omitempty"`
	// CopyCtor, MoveCtor, CopyAssign, MoveAssign, Dtor
	UserDeclared bool `json:"userDeclared,omitempty"`
	// DefaultCtor
	UserProvided bool `json:"userProvided,omitempty"`
}

type DefinitionData struct {
	CanConstDefaultInit                bool           `json:"canConstDefaultInit,omitempty"`
	CanPassInRegisters                 bool           `json:"canPassInRegisters,omitempty"`
	CopyAssign                         DefinitionCtor `json:"copyAssign,omitempty"`
	CopyCtor                           DefinitionCtor `json:"copyCtor,omitempty"`
	DefaultCtor                        DefinitionCtor `json:"defaultCtor,omitempty"`
	Dtor                               DefinitionCtor `json:"dtor,omitempty"`
	HasConstexprNonCopyMoveConstructor bool           `json:"hasConstexprNonCopyMoveConstructor,omitempty"`
	HasMutableFields                   bool           `json:"hasMutableFields,omitempty"`
	HasUserDeclaredConstructor         bool           `json:"hasUserDeclaredConstructor,omitempty"`
	HasVariantMembers                  bool           `json:"hasVariantMembers,omitempty"`
	IsAbstract                         bool           `json:"isAbstract,omitempty"`
	IsAggregate                        bool           `json:"isAggregate,omitempty"`
	IsEmpty                            bool           `json:"isEmpty,omitempty"`
	IsGenericLambda                    bool           `json:"isGenericLambda,omitempty"`
	IsLambda                           bool           `json:"isLambda,omitempty"`
	IsLiteral                          bool           `json:"isLiteral,omitempty"`
	IsPOD                              bool           `json:"isPOD,omitempty"`
	IsPolymorphic                      bool           `json:"isPolymorphic,omitempty"`
	IsStandardLayout                   bool           `json:"isStandardLayout,omitempty"`
	IsTrivial                          bool           `json:"isTrivial,omitempty"`
	IsTriviallyCopyable                bool           `json:"isTriviallyCopyable,omitempty"`
	MoveAssign                         DefinitionCtor `json:"moveAssign,omitempty"`
	MoveCtor                           DefinitionCtor `json:"moveCtor,omitempty"`
}

type Access string

const (
	None      Access = "none"
	Private   Access = "private"
	Protected Access = "protected"
	Public    Access = "public"
	Package   Access = "package"
)

type CXXBase struct {
	Type            *Type  `json:"type"`
	Access          Access `json:"access"`
	WrittenAccess   Access `json:"writtenAccess"`
	IsVirtual       bool   `json:"isVirtual,omitempty"`
	IsPackExpansion bool   `json:"isPackExpansion,omitempty"`
}

type HtmlAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
