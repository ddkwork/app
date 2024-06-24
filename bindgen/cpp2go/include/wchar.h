#ifndef __WCHAR_H__
#define __WCHAR_H__

#include <stdarg.h>
#include <go/fmt_go.h>
#include <go/os_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef NULL
#ifdef __cplusplus
#define NULL 0
#else
#define NULL ((void *)0)
#endif
#endif

#ifndef size_t
#define size_t int
#endif

#ifndef wchar_t
typedef int rune;
#define wchar_t rune
#endif

#ifndef _WINT_T_DEFINED_
 #define _WINT_T_DEFINED_
 typedef int wint_t;
#endif

#undef WCHAR_MIN
#undef WCHAR_MAX
#define WCHAR_MIN   0
#define WCHAR_MAX   65535

#undef WEOF
#define WEOF ((wint_t)-1)

typedef struct {
    int _;
} mbstate_t;

struct tm;
extern size_t wcsftime(wchar_t *, size_t, const wchar_t *, const struct tm *);

#define swprintf(s, n, format, arg...) fmt.Appendf(s, format, ##arg)
#define vswprintf(s, n, format, arg) fmt.Appendf(s, format, arg)
#define swscanf(s, format, arg...) fmt.Sscanf(s, format, ##arg)
#define vswscanf(s, format, arg) fmt.Sscanf(s, format, arg)

#define fwprintf(f, format, arg...) fmt.Fprintf(f, format, ##arg)
#define vfwprintf(f, format, arg) fmt.Fprintf(f, format, arg)
#define fwscanf(f, format, arg...) fmt.Fscanf(f, format, ##arg)
#define vfwscanf(f, format, arg) fmt.Fscanf(f, format, arg)

#define wprintf(format, arg...) fmt.Printf(format, ##arg)
#define vwprintf(format, arg) fmt.Printf(format, arg)
#define wscanf(format, arg...) fmt.Scanf(format, ##arg)
#define vwscanf(format, arg) fmt.Scanf(format, arg)

wint_t fgetwc(void*);
#define fgetws(buf, n, f) f->Read(buf, n)
#define fputwc(c, f) fmt.Fprintf(f, "%c", c)
#define fputws(s, f) f->WriteString(s)
int fwide(void*, int);
#define getwc(f) fgetwc(f)
#define getwchar() fgetwc(os.Stdin)
#define putwc(c, f) fputwc(c, f)
#define putwchar(c) fmt.Printf("%c", c)
wint_t ungetwc(wint_t, void*);

wint_t btowc(int);
wint_t wctob(int);
int mbsinit(const mbstate_t *);
size_t mbrlen(const char *, size_t, mbstate_t *);
size_t mbrtowc(wchar_t *, const char *, size_t, mbstate_t *);
size_t wcrtomb(char *, wchar_t, mbstate_t *);
size_t mbsrtowcs(wchar_t *, const char **, size_t, mbstate_t *);
size_t wcsrtombs(char *, const wchar_t **, size_t, mbstate_t *);
size_t mbsnrtowcs(wchar_t *, const char **, size_t, size_t, mbstate_t *);
size_t wcsnrtombs(char *, const wchar_t **, size_t, size_t, mbstate_t *);

double wcstod(const wchar_t *, wchar_t **);
float wcstof(const wchar_t *, wchar_t **);
long double wcstold(const wchar_t *, wchar_t **);
long int wcstol(const wchar_t *, wchar_t **, int);
unsigned long int wcstoul(const wchar_t *, wchar_t **, int);
long long int wcstoll(const wchar_t *, wchar_t **, int);
unsigned long long int wcstoull(const wchar_t *, wchar_t **, int);

wchar_t *wcscpy(wchar_t *, const wchar_t *);
wchar_t *wcsncpy(wchar_t *, const wchar_t *, size_t);
wchar_t *wmemcpy(wchar_t *, const wchar_t *, size_t);
wchar_t *wmemmove(wchar_t *, const wchar_t *, size_t);
wchar_t *wcscat(wchar_t *, const wchar_t *);
wchar_t *wcsncat(wchar_t *, const wchar_t *, size_t);
int wcscmp(const wchar_t *, const wchar_t *);
int wcsncmp(const wchar_t *, const wchar_t *, size_t);
int wcscasecmp(const wchar_t *, const wchar_t *);
int wcsncasecmp(const wchar_t *, const wchar_t *, size_t);
int wcscoll(const wchar_t *, const wchar_t *);
size_t wcsxfrm(wchar_t *, const wchar_t *, size_t);
int wmemcmp(const wchar_t *, const wchar_t *, size_t);
size_t wcscspn(const wchar_t *, const wchar_t *);
size_t wcsspn(const wchar_t *, const wchar_t *);
wchar_t *wcstok(wchar_t *, const wchar_t *, wchar_t **);
size_t wcslen(const wchar_t *);
wchar_t *wmemset(wchar_t *, wchar_t, size_t);

wchar_t *wcschr(const wchar_t *, wchar_t);
wchar_t *wcspbrk(const wchar_t *, const wchar_t *);
wchar_t *wcsrchr(const wchar_t *, wchar_t);
wchar_t *wcsstr(const wchar_t *, const wchar_t *);
wchar_t *wmemchr(const wchar_t *, wchar_t, size_t);

#ifdef __cplusplus
}
#endif

#endif
