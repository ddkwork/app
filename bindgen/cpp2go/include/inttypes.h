#ifndef _INTTYPES_H_
#define _INTTYPES_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    intmax_t quot;
    intmax_t rem;
} imaxdiv_t;

#define PRId8 "d"
#define PRId16 "d"
#define PRId32 "d"
#define PRId64 "I64d"

#define PRIdLEAST8 "d"
#define PRIdLEAST16 "d"
#define PRIdLEAST32 "d"
#define PRIdLEAST64 "I64d"

#define PRIdFAST8 "d"
#define PRIdFAST16 "d"
#define PRIdFAST32 "d"
#define PRIdFAST64 "I64d"
#define PRIdMAX "I64d"

#define PRIi8 "i"
#define PRIi16 "i"
#define PRIi32 "i"
#define PRIi64 "I64i"

#define PRIiLEAST8 "i"
#define PRIiLEAST16 "i"
#define PRIiLEAST32 "i"
#define PRIiLEAST64 "I64i"

#define PRIiFAST8 "i"
#define PRIiFAST16 "i"
#define PRIiFAST32 "i"
#define PRIiFAST64 "I64i"
#define PRIiMAX "I64i"

#define PRIo8 "o"
#define PRIo16 "o"
#define PRIo32 "o"
#define PRIo64 "I64o"

#define PRIoLEAST8 "o"
#define PRIoLEAST16 "o"
#define PRIoLEAST32 "o"
#define PRIoLEAST64 "I64o"

#define PRIoFAST8 "o"
#define PRIoFAST16 "o"
#define PRIoFAST32 "o"
#define PRIoFAST64 "I64o"
#define PRIoMAX "I64o"

#define PRIu8 "u"
#define PRIu16 "u"
#define PRIu32 "u"
#define PRIu64 "I64u"

#define PRIuLEAST8 "u"
#define PRIuLEAST16 "u"
#define PRIuLEAST32 "u"
#define PRIuLEAST64 "I64u"

#define PRIuFAST8 "u"
#define PRIuFAST16 "u"
#define PRIuFAST32 "u"
#define PRIuFAST64 "I64u"
#define PRIuMAX "I64u"

#define PRIx8 "x"
#define PRIx16 "x"
#define PRIx32 "x"
#define PRIx64 "I64x"

#define PRIxLEAST8 "x"
#define PRIxLEAST16 "x"
#define PRIxLEAST32 "x"
#define PRIxLEAST64 "I64x"

#define PRIxFAST8 "x"
#define PRIxFAST16 "x"
#define PRIxFAST32 "x"
#define PRIxFAST64 "I64x"
#define PRIxMAX "I64x"

#define PRIX8 "X"
#define PRIX16 "X"
#define PRIX32 "X"
#define PRIX64 "I64X"

#define PRIXLEAST8 "X"
#define PRIXLEAST16 "X"
#define PRIXLEAST32 "X"
#define PRIXLEAST64 "I64X"

#define PRIXFAST8 "X"
#define PRIXFAST16 "X"
#define PRIXFAST32 "X"
#define PRIXFAST64 "I64X"
#define PRIXMAX "I64X"

#define SCNd16 "hd"
#define SCNd32 "d"
#define SCNd64 "I64d"

#define SCNdLEAST16 "hd"
#define SCNdLEAST32 "d"
#define SCNdLEAST64 "I64d"

#define SCNdFAST16 "hd"
#define SCNdFAST32 "d"
#define SCNdFAST64 "I64d"
#define SCNdMAX "I64d"

#define SCNi16 "hi"
#define SCNi32 "i"
#define SCNi64 "I64i"

#define SCNiLEAST16 "hi"
#define SCNiLEAST32 "i"
#define SCNiLEAST64 "I64i"

#define SCNiFAST16 "hi"
#define SCNiFAST32 "i"
#define SCNiFAST64 "I64i"
#define SCNiMAX "I64i"

#define SCNo16 "ho"
#define SCNo32 "o"
#define SCNo64 "I64o"

#define SCNoLEAST16 "ho"
#define SCNoLEAST32 "o"
#define SCNoLEAST64 "I64o"

#define SCNoFAST16 "ho"
#define SCNoFAST32 "o"
#define SCNoFAST64 "I64o"
#define SCNoMAX "I64o"

#define SCNx16 "hx"
#define SCNx32 "x"
#define SCNx64 "I64x"

#define SCNxLEAST16 "hx"
#define SCNxLEAST32 "x"
#define SCNxLEAST64 "I64x"

#define SCNxFAST16 "hx"
#define SCNxFAST32 "x"
#define SCNxFAST64 "I64x"
#define SCNxMAX "I64x"

#define SCNu16 "hu"
#define SCNu32 "u"
#define SCNu64 "I64u"

#define SCNuLEAST16 "hu"
#define SCNuLEAST32 "u"
#define SCNuLEAST64 "I64u"

#define SCNuFAST16 "hu"
#define SCNuFAST32 "u"
#define SCNuFAST64 "I64u"
#define SCNuMAX "I64u"

#define PRIdPTR "d"
#define PRIiPTR "i"
#define PRIoPTR "o"
#define PRIuPTR "u"
#define PRIxPTR "x"
#define PRIXPTR "X"
#define SCNdPTR "d"
#define SCNiPTR "i"
#define SCNoPTR "o"
#define SCNxPTR "x"
#define SCNuPTR "u"

#define SCNd8 "hhd"
#define SCNdLEAST8 "hhd"
#define SCNdFAST8 "hhd"

#define SCNi8 "hhi"
#define SCNiLEAST8 "hhi"
#define SCNiFAST8 "hhi"

#define SCNo8 "hho"
#define SCNoLEAST8 "hho"
#define SCNoFAST8 "hho"

#define SCNx8 "hhx"
#define SCNxLEAST8 "hhx"
#define SCNxFAST8 "hhx"

#define SCNu8 "hhu"
#define SCNuLEAST8 "hhu"
#define SCNuFAST8 "hhu"

intmax_t imaxabs(intmax_t);
imaxdiv_t imaxdiv(intmax_t, intmax_t);

intmax_t strtoimax(const char*, char**, int);
uintmax_t strtoumax(const char*, char**, int);

intmax_t wcstoimax(const wchar_t*, wchar_t**, int);
uintmax_t wcstoumax(const wchar_t*, wchar_t**, int);

#ifdef __cplusplus
}
#endif

#endif
