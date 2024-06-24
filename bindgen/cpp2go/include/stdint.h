#ifndef __STDINT_H__
#define __STDINT_H__

#include <go/math_go.h>

typedef signed char         int8;
typedef unsigned char       uint8;
typedef short               int16;
typedef unsigned short      uint16;
typedef int                 int32;
typedef unsigned            uint32;
typedef long long           int64;
typedef unsigned long long  uint64;
typedef unsigned            uint;
typedef unsigned long long  uintptr;

#define int8_t      int8
#define uint8_t     uint8
#define int16_t     int16
#define uint16_t    uint16
#define int32_t     int32
#define uint32_t    uint32
#define int64_t     int64
#define uint64_t    uint64

#define int_least8_t    int8
#define uint_least8_t   uint8
#define int_least16_t   int16
#define uint_least16_t  uint16
#define int_least32_t   int32
#define uint_least32_t  uint32
#define int_least64_t   int64
#define uint_least64_t  uint64

#define int_fast8_t     int8
#define uint_fast8_t    uint8
#define int_fast16_t    int16
#define uint_fast16_t   uint16
#define int_fast32_t    int32
#define uint_fast32_t   uint32
#define int_fast64_t    int64
#define uint_fast64_t   uint64

#define intmax_t    int64
#define uintmax_t   uint64

#define intptr_t    uintptr
#define uintptr_t   uintptr

#define INT8_MIN  math.MinInt8
#define INT16_MIN math.MinInt16
#define INT32_MIN math.MinInt32
#define INT64_MIN math.MinInt64

#define INT8_MAX  math.MaxInt8
#define INT16_MAX math.MaxInt16
#define INT32_MAX math.MaxInt32
#define INT64_MAX math.MaxInt64

#define UINT8_MAX  math.MaxUint8
#define UINT16_MAX math.MaxUint16
#define UINT32_MAX math.MaxUint32
#define UINT64_MAX math.MaxUint64

#define INT_LEAST8_MIN  INT8_MIN
#define INT_LEAST16_MIN INT16_MIN
#define INT_LEAST32_MIN INT32_MIN
#define INT_LEAST64_MIN INT64_MIN

#define INT_LEAST8_MAX  INT8_MAX
#define INT_LEAST16_MAX INT16_MAX
#define INT_LEAST32_MAX INT32_MAX
#define INT_LEAST64_MAX INT64_MAX

#define UINT_LEAST8_MAX  UINT8_MAX
#define UINT_LEAST16_MAX UINT16_MAX
#define UINT_LEAST32_MAX UINT32_MAX
#define UINT_LEAST64_MAX UINT64_MAX

#define INT_FAST8_MIN  INT8_MIN
#define INT_FAST16_MIN INT16_MIN
#define INT_FAST32_MIN INT32_MIN
#define INT_FAST64_MIN INT64_MIN

#define INT_FAST8_MAX  INT8_MAX
#define INT_FAST16_MAX INT16_MAX
#define INT_FAST32_MAX INT32_MAX
#define INT_FAST64_MAX INT64_MAX

#define UINT_FAST8_MAX  UINT8_MAX
#define UINT_FAST16_MAX UINT16_MAX
#define UINT_FAST32_MAX UINT32_MAX
#define UINT_FAST64_MAX UINT64_MAX

#define INTPTR_MIN  INT64_MIN
#define INTPTR_MAX  INT64_MAX
#define UINTPTR_MAX UINT64_MAX

#define INTMAX_MIN  math.MinInt
#define INTMAX_MAX  math.MaxInt
#define UINTMAX_MAX math.MaxUint

#define PTRDIFF_MIN INT64_MIN
#define PTRDIFF_MAX INT64_MAX

#define SIG_ATOMIC_MIN INTMAX_MIN
#define SIG_ATOMIC_MAX INTMAX_MAX

#ifndef SIZE_MAX
#define SIZE_MAX UINTMAX_MAX
#endif

#ifndef WCHAR_MIN
#define WCHAR_MIN 0U
#define WCHAR_MAX 0xffffU
#endif

#define WINT_MIN 0U
#define WINT_MAX 0xffffU

#define INT8_C(val) (INT_LEAST8_MAX-INT_LEAST8_MAX+(val))
#define INT16_C(val) (INT_LEAST16_MAX-INT_LEAST16_MAX+(val))
#define INT32_C(val) (INT_LEAST32_MAX-INT_LEAST32_MAX+(val))
#define INT64_C(val) val##LL

#define UINT8_C(val) (val)
#define UINT16_C(val) (val)
#define UINT32_C(val) (val##U)
#define UINT64_C(val) val##ULL

#define INTMAX_C(val) val##LL
#define UINTMAX_C(val) val##ULL

#endif
