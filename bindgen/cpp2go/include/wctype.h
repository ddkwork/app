#ifndef __WCTYPE_H__
#define __WCTYPE_H__

#include <go/unicode_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef _WINT_T_DEFINED_
#define _WINT_T_DEFINED_
typedef int wint_t;
#endif

typedef void *wctype_t;
typedef void *wctrans_t;

#undef WEOF
#define WEOF ((wint_t)-1)

int iswalnum(wint_t);
#define iswalpha(c) unicode.IsLetter(c)
int iswblank(wint_t);
#define iswcntrl(c) unicode.IsControl(c)
#define iswgraph(c) unicode.IsGraphic(c)
#define iswprint(c) unicode.IsPrint(c)
#define iswpunct(c) unicode.IsPunct(c)
#define iswspace(c) unicode.IsSpace(c)
#define iswlower(c) unicode.IsLower(c)
#define iswupper(c) unicode.IsUpper(c)
#define iswdigit(c) unicode.IsDigit(c)
int iswxdigit(wint_t);
#define towlower(c) unicode.ToLower(c)
#define towupper(c) unicode.ToUpper(c)

wctype_t wctype(const char *);
int iswctype(wint_t, wctype_t);

wctrans_t wctrans(const char *);
wint_t towctrans(wint_t, wctrans_t);

#ifdef __cplusplus
}
#endif

#endif
