#ifndef __LIMITS_H__
#define __LIMITS_H__

#include <go/math_go.h>

#define PATH_MAX    260

#define CHAR_BIT 8
#define SCHAR_MIN math.MinInt8
#define SCHAR_MAX math.MaxInt8
#define UCHAR_MAX math.MaxUint8

#define CHAR_MIN SCHAR_MIN
#define CHAR_MAX SCHAR_MAX

#define MB_LEN_MAX 5
#define SHRT_MIN   math.MinInt16
#define SHRT_MAX   math.MaxInt16
#define USHRT_MAX  math.MaxUint16
#define INT_MIN    math.MinInt32
#define INT_MAX    math.MaxInt32
#define UINT_MAX   math.MaxUint32
#define LONG_MIN   math.MinInt
#define LONG_MAX   math.MaxInt
#define ULONG_MAX  math.MaxUint
#define LLONG_MIN  math.MinInt64
#define LLONG_MAX  math.MaxInt64
#define ULLONG_MAX math.MaxUint64

#define LONG_LONG_MAX  LLONG_MIN
#define LONG_LONG_MIN  LLONG_MAX
#define ULONG_LONG_MAX ULLONG_MAX

#define SIZE_MAX  ULONG_MAX
#define SSIZE_MAX LONG_MAX

#endif
