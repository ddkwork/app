package ast

type CastExpr_St struct {
	CastKind             CastKind   `json:"castKind,omitempty"`
	ConversionFunc       *BareDecl  `json:"conversionFunc,omitempty"`
	IsPartOfExplicitCast bool       `json:"isPartOfExplicitCast,omitempty"`
	Path                 []CastPath `json:"path,omitempty"`
}

type CXXConstructExpr_St struct {
	// CXXConstructExpr
	ConstructionKind string `json:"constructionKind,omitempty"`
	// CXXConstructExpr
	CtorType *Type `json:"ctorType,omitempty"`
	// CXXConstructExpr
	Elidable bool `json:"elidable,omitempty"`
	// CXXConstructExpr
	HadMultipleCandidates bool `json:"hadMultipleCandidates,omitempty"`
	// CXXConstructExpr
	InitList bool `json:"initializer_list,omitempty"`
	// CXXUnresolvedConstructExpr, CXXConstructExpr
	List bool `json:"list,omitempty"`
	// CXXUnresolvedConstructExpr
	TypeAsWritten *Type `json:"typeAsWritten,omitempty"`
	// CXXConstructExpr
	Zeroing bool `json:"zeroing,omitempty"`
}

type CXXNewDeleteExpr_St struct {
	// CXXNewExpr
	InitStyle string `json:"initStyle,omitempty"`
	// CXXNewExpr, CXXDeleteExpr
	IsArray bool `json:"isArray,omitempty"`
	// CXXDeleteExpr
	IsArrayAsWritten bool `json:"isArrayAsWritten,omitempty"`
	// CXXNewExpr, CXXDeleteExpr
	IsGlobal bool `json:"isGlobal,omitempty"`
	// CXXNewExpr
	IsPlacement bool `json:"isPlacement,omitempty"`
	// CXXNewExpr, CXXDeleteExpr
	OperatorDeleteDecl *BareDecl `json:"operatorDeleteDecl,omitempty"`
	// CXXNewExpr
	OperatorNewDecl *BareDecl `json:"operatorNewDecl,omitempty"`
}

type IfSwitchWhileStmt_St struct {
	// IfStmt
	HasElse bool `json:"hasElse,omitempty"`
	// IfStmt, SwitchStmt
	HasInit bool `json:"hasInit,omitempty"`
	// IfStmt, SwitchStmt, WhileStmt
	HasVar bool `json:"hasVar,omitempty"`
	// IfStmt
	IsConstexpr bool `json:"isConstexpr,omitempty"`
}

type ParamCommandComment_St struct {
	// ParamCommandComment
	Direction string `json:"direction,omitempty"`
	// ParamCommandComment
	Explicit bool `json:"explicit,omitempty"`
	// ParamCommandComment, TParamCommandComment
	Param string `json:"param,omitempty"`
	// ParamCommandComment
	ParamIdx uint `json:"paramIdx,omitempty"`
	// TParamCommandComment
	Positions []int `json:"positions,omitempty"`
}

type Node struct {
	// Attr, Stmt, Type, QualType, Decl, Comment
	ID ID `json:"id,omitempty"`
	// Attr, Stmt, Type, QualType, Decl, Comment, TemplateArgument, CXXCtorInitializer, Capture
	Kind Kind `json:"kind,omitempty"`
	// Decl, Comment
	Loc *Loc `json:"loc,omitempty"`
	// Attr, Stmt, Decl, Comment, TemplateArgument
	Range *Range `json:"range,omitempty"`
	// Decl
	IsImplicit bool `json:"isImplicit,omitempty"` // is this type implicit defined
	// Decl
	IsReferenced bool `json:"isReferenced,omitempty"` // is this type refered or not
	// Decl
	IsUsed bool `json:"isUsed,omitempty"` // is this variable used or not
	// Decl
	ParentDeclContextID ID `json:"parentDeclContextId,omitempty"`
	// Decl
	PreviousDecl ID `json:"previousDecl,omitempty"`
	// NamedDecl, TypedefDecl, TypeAliasDecl, NamespaceDecl, NamespaceAliasDecl, UsingDecl, VarDecl, FieldDecl
	// FunctionDecl, EnumDecl, EnumConstantDecl, RecordDecl, CXXRecordDecl, TemplateTypeParmDecl
	// NonTypeTemplateParmDecl, TemplateTemplateParmDecl, ObjCIvarDecl, PredefinedExpr, MemberExpr
	// UnaryExprOrTypeTraitExpr, SizeOfPackExpr, UnresolvedLookupExpr, AddrLabelExpr, LabelStmt
	// InlineCommandComment, HTMLStartTagComment, HTMLEndTagComment, BlockCommandComment, VerbatimBlockComment
	Name        string `json:"name,omitempty"`
	MangledName string `json:"mangledName,omitempty"`
	// Stmt, Type, QualType, TypedefDecl, TypeAliasDecl, VarDecl, FieldDecl, FunctionDecl, EnumConstantDecl
	// NonTypeTemplateParmDecl, FriendDecl, ObjCIvarDecl, TemplateArgument
	Type *Type `json:"type,omitempty"`
	// RecordDecl, CXXRecordDecl, TemplateTypeParmDecl
	TagUsed string `json:"tagUsed,omitempty"` // struct | union
	// VarDecl, FunctionDecl
	Inline bool `json:"inline,omitempty"`
	// Stmt
	ValueCategory ValueCategory `json:"valueCategory,omitempty"`

	// AccessSpecDecl, ObjCIvarDecl
	Access Access `json:"access,omitempty"`
	// NamespaceAliasDecl
	AliasedNamespace *BareDecl `json:"aliasedNamespace,omitempty"`
	// CXXCtorInitializer
	AnyInit *BareDecl `json:"anyInit,omitempty"`
	// UnaryExprOrTypeTraitExpr
	ArgType *Type `json:"argType,omitempty"`
	// InlineCommandComment, BlockCommandComment
	Args []string `json:"args,omitempty"`
	// HTMLStartTagComment
	Attrs []HtmlAttr `json:"attrs,omitempty"`
	// CXXCtorInitializer
	BaseInit *Type `json:"baseInit,omitempty"`
	// CXXRecordDecl
	Bases []CXXBase `json:"bases,omitempty"`
	// Capture
	Byref bool `json:"byref,omitempty"`
	// BlockDecl
	CapturesThis bool `json:"capturesThis,omitempty"`
	// RecordDecl, CXXRecordDecl
	CompleteDefinition bool `json:"completeDefinition,omitempty"`
	// FunctionProtoType
	ConditionEvaluatesTo bool `json:"conditionEvaluatesTo,omitempty"`
	// FunctionProtoType
	Const bool `json:"const,omitempty"`
	// VarDecl, FunctionDecl
	Constexpr bool `json:"constexpr,omitempty"`
	// Type
	ContainsUnexpandedPack bool `json:"containsUnexpandedPack,omitempty"`
	// TypedefType, UnresolvedUsingType, TagType, TemplateTypeParmType, InjectedClassNameType
	// ObjCInterfaceType, TemplateArgument
	Decl *BareDecl `json:"decl,omitempty"`
	// LabelStmt
	DeclId ID `json:"declId,omitempty"`
	// TemplateTypeParmDecl, NonTypeTemplateParmDecl, TemplateTemplateParmDecl
	DefaultArg *Node `json:"defaultArg,omitempty"`
	// CXXRecordDecl
	DefinitionData *DefinitionData `json:"definitionData,omitempty"`
	// CXXCtorInitializer
	DelegatingInit *Type `json:"delegatingInit,omitempty"`
	// TemplateTypeParmType, TemplateTypeParmDecl, NonTypeTemplateParmDecl, TemplateTemplateParmDecl
	Depth uint `json:"depth,omitempty"`
	// FunctionProtoType
	ExceptionSpec string `json:"exceptionSpec,omitempty"`
	// FunctionProtoType
	ExceptionTypes string `json:"exceptionTypes,omitempty"`
	// CXXDependentScopeMemberExpr
	ExplicitTemplateArgs []Node `json:"explicitTemplateArgs,omitempty"`
	// FunctionDecl
	ExplicitlyDefaulted string `json:"explicitlyDefaulted,omitempty"`
	// FunctionDecl
	ExplicitlyDeleted bool `json:"explicitlyDeleted,omitempty"`
	// EnumDecl
	FixedUnderlyingType *Type `json:"fixedUnderlyingType,omitempty"`
	// DeclRefExpr
	FoundReferencedDecl *BareDecl `json:"foundReferencedDecl,omitempty"`
	// TemplateArgument
	FromDecl *BareDecl `json:"fromDecl,omitempty"`

	// CXXDependentScopeMemberExpr
	HasExplicitTemplateArgs bool `json:"hasExplicitTemplateArgs,omitempty"`
	// FieldDecl
	HasInClassInitializer bool `json:"hasInClassInitializer,omitempty"`
	// CXXDependentScopeMemberExpr
	HasTemplateKeyword bool `json:"hasTemplateKeyword,omitempty"`
	// Attr, CXXThisExpr
	Implicit bool `json:"implicit,omitempty"`
	// TemplateTypeParmType, TemplateTypeParmDecl, NonTypeTemplateParmDecl, TemplateTemplateParmDecl
	Index uint `json:"index,omitempty"`
	// Attr
	Inherited bool `json:"inherited,omitempty"`
	// VarDecl
	Init string `json:"init,omitempty"`
	// MemberExpr, CXXDependentScopeMemberExpr
	IsArrow bool `json:"isArrow,omitempty"` // is ptr->member not obj.member
	// VarDecl, FieldDecl
	IsBitfield bool `json:"isBitfield,omitempty"`
	// Type
	IsDependent bool `json:"isDependent,omitempty"`
	// TemplateArgument
	IsExpr bool `json:"isExpr,omitempty"`
	// Decl
	IsHidden bool `json:"isHidden,omitempty"`
	// Type
	IsImported bool `json:"isImported,omitempty"`
	// NamespaceDecl
	IsInline bool `json:"isInline,omitempty"`
	// Type
	IsInstantiationDependent bool `json:"isInstantiationDependent,omitempty"`
	// Decl
	IsInvalid bool `json:"isInvalid,omitempty"`
	// TemplateArgument
	IsNull bool `json:"isNull,omitempty"`
	// TemplateArgument
	IsNullptr bool `json:"isNullptr,omitempty"`
	// TemplateTypeParmType, TemplateArgument
	IsPack bool `json:"isPack,omitempty"`
	// VarDecl, TemplateTypeParmDecl, NonTypeTemplateParmDecl, TemplateTemplateParmDecl
	IsParameterPack bool `json:"isParameterPack,omitempty"`
	// Type
	IsVariablyModified bool `json:"isVariablyModified,omitempty"`
	// AddrLabelExpr
	LabelDeclId ID `json:"labelDeclId,omitempty"`
	// UnresolvedLookupExpr
	Lookups []BareDecl `json:"lookups,omitempty"`
	// HTMLStartTagComment
	Malformed bool `json:"malformed,omitempty"`
	// CXXDependentScopeMemberExpr
	Member string `json:"member,omitempty"`
	// VarDecl, FieldDecl
	ModulePrivate bool `json:"modulePrivate,omitempty"`
	// FieldDecl
	Mutable bool `json:"mutable,omitempty"`
	// Capture
	Nested bool `json:"nested,omitempty"`
	// DeclRefExpr, MemberExpr
	NonOdrUseReason string `json:"nonOdrUseReason,omitempty"`
	// VarDecl
	Nrvo bool `json:"nrvo,omitempty"`

	// NamespaceDecl
	OriginalNamespace *BareDecl `json:"originalNamespace,omitempty"`
	// FunctionDecl
	Pure bool `json:"pure,omitempty"`
	// QualType
	Qualifiers string `json:"qualifiers,omitempty"`
	// FunctionProtoType
	RefQualifier string `json:"refQualifier,omitempty"`
	// DeclRefExpr
	ReferencedDecl *BareDecl `json:"referencedDecl,omitempty"`
	// MemberExpr
	ReferencedMemberDecl ID `json:"referencedMemberDecl,omitempty"`
	// InlineCommandComment
	RenderKind string `json:"renderKind,omitempty"`
	// FunctionProtoType
	Restrict bool `json:"restrict,omitempty"`
	// EnumDecl
	ScopedEnumTag string `json:"scopedEnumTag,omitempty"`
	// HTMLStartTagComment
	SelfClosing bool `json:"selfClosing,omitempty"`
	// VarDecl, FunctionDecl
	StorageClass StorageClass `json:"storageClass,omitempty"`
	// ObjCIvarDecl
	Synthesized bool `json:"synthesized,omitempty"`
	// FunctionProtoType
	ThrowsAny bool `json:"throwsAny,omitempty"`
	// VarDecl
	TLS string `json:"tls,omitempty"`
	// FunctionProtoType
	TrailingReturn bool `json:"trailingReturn,omitempty"`

	// UnresolvedLookupExpr
	UsesADL bool `json:"usesADL,omitempty"`
	// ConstantExpr, IntegerLiteral, CharacterLiteral, FixedPointLiteral, FloatingLiteral, StringLiteral
	// CXXBoolLiteralExpr, TemplateArgument
	Value interface{} `json:"value,omitempty"`
	// Capture
	Var *BareDecl `json:"var,omitempty"`
	// FunctionProtoType, FunctionDecl, BlockDecl
	Variadic bool `json:"variadic,omitempty"`
	// FunctionDecl
	Virtual bool `json:"virtual,omitempty"`
	// FunctionProtoType
	Volatile bool `json:"volatile,omitempty"`

	// CastExpr
	*CastExpr_St
	// CXXUnresolvedConstructExpr, CXXConstructExpr
	*CXXConstructExpr_St
	// CXXNewExpr, CXXDeleteExpr
	*CXXNewDeleteExpr_St
	// ParamCommandComment, TParamCommandComment
	*ParamCommandComment_St
	// IfStmt, SwitchStmt, WhileStmt
	*IfSwitchWhileStmt_St

	*SimpleNode
	Inner       []*Node `json:"inner,omitempty"`
	ArrayFiller []*Node `json:"array_filler,omitempty"`
}

type SimpleNode struct {
	// CallExpr
	Adl bool `json:"adl,omitempty"`
	// DependentSizedExtVectorType
	AttrLoc *Loc `json:"attrLoc,omitempty"`
	// VerbatimBlockComment
	CloseName string `json:"closeName,omitempty"`
	// InitListExpr
	Field *Node `json:"field,omitempty"`
	// ObjCAtCatchStmt
	IsCatchAll bool `json:"isCatchAll,omitempty"`
	// CaseStmt
	IsGNURange bool `json:"isGNURange,omitempty"`
	// MacroQualifiedType
	MacroName string `json:"macroName,omitempty"`
	// UsingDirectiveDecl
	NominatedNamespace *BareDecl `json:"nominatedNamespace,omitempty"`
	// PackExpansionType
	NumExpansions uint `json:"numExpansions,omitempty"`
	// GenericSelectionExpr
	ResultDependent bool `json:"resultDependent,omitempty"`
	// ReferenceType
	SpelledAsLValue bool `json:"spelledAsLValue,omitempty"`
	// UsingShadowDecl
	Target *BareDecl `json:"target,omitempty"`
	// GotoStmt
	TargetLabelDeclId ID `json:"targetLabelDeclId,omitempty"`
	// TextComment, VerbatimBlockLineComment, VerbatimLineComment
	Text string `json:"text,omitempty"`
	// UnaryTransformType
	TransformKind string `json:"transformKind,omitempty"`

	// UnaryOperator
	IsPostfix bool `json:"isPostfix,omitempty"`
	// UnaryOperator, BinaryOperator
	OpCode OpCode `json:"opcode,omitempty"`

	// ElaboratedType
	OwnedTagDecl *BareDecl `json:"ownedTagDecl,omitempty"`
	Qualifier    string    `json:"qualifier,omitempty"`
	// MaterializeTemporaryExpr
	BoundToLValueRef bool      `json:"boundToLValueRef,omitempty"`
	ExtendingDecl    *BareDecl `json:"extendingDecl,omitempty"`
	StorageDuration  string    `json:"storageDuration,omitempty"`
	// ExprWithCleanups
	Cleanups                []BareDecl `json:"cleanups,omitempty"`
	CleanupsHaveSideEffects bool       `json:"cleanupsHaveSideEffects,omitempty"`
	// ConstAssociation
	AssociationKind AssociationKind `json:"associationKind,omitempty"`
	Selected        bool            `json:"selected,omitempty"`
	// CXXBindTemporaryExpr
	Dtor *BareDecl `json:"dtor,omitempty"`
	Temp ID        `json:"temp,omitempty"`
	// FunctionType
	CC             CallConv `json:"cc,omitempty"`
	NoReturn       bool     `json:"noreturn,omitempty"`
	ProducesResult bool     `json:"producesResult,omitempty"`
	RegParm        uint     `json:"regParm,omitempty"`
	// CXXTypeidExpr
	AdjustedTypeArg *Type `json:"adjustedTypeArg,omitempty"`
	TypeArg         *Type `json:"typeArg,omitempty"`
	// ArrayType
	IndexTypeQualifiers string `json:"indexTypeQualifiers,omitempty"`
	Size                int    `json:"size,omitempty"` // array size
	SizeModifier        string `json:"sizeModifier,omitempty"`
	// VectorType
	NumElements uint   `json:"numElements,omitempty"`
	VectorKind  string `json:"vectorKind,omitempty"`
	// AutoType
	TypeKeyword string `json:"typeKeyword,omitempty"`
	Undeduced   bool   `json:"undeduced,omitempty"`
	// TemplateSpecializationType
	IsAlias      bool   `json:"isAlias,omitempty"`
	TemplateName string `json:"templateName,omitempty"`
	// MemberPointerType
	IsData     bool `json:"isData,omitempty"`
	IsFunction bool `json:"isFunction,omitempty"`
	// LinkageSpecDecl
	Language  string `json:"language,omitempty"`
	HasBraces bool   `json:"hasBraces,omitempty"`
	// CompoundAssignOperator
	ComputeLHSType    *Type `json:"computeLHSType,omitempty"`
	ComputeResultType *Type `json:"computeResultType,omitempty"`
}
