#ifndef __CTYPE_H__
#define __CTYPE_H__

#include <go/unicode_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#define isalpha(c) unicode.IsLetter(c)
#define isupper(c) unicode.IsUpper(c)
#define islower(c) unicode.IsLower(c)
#define isdigit(c) unicode.IsDigit(c)
int isxdigit(int);
#define isspace(c) unicode.IsSpace(c)
#define ispunct(c) unicode.IsPunct(c)
int isalnum(int);
#define isprint(c) unicode.IsPrint(c)
#define isgraph(c) unicode.IsGraphic(c)
#define iscntrl(c) unicode.IsControl(c)
#define toupper(c) unicode.ToUpper(c)
#define tolower(c) unicode.ToLower(c)
int isblank(int);

int isascii(int);
int toascii(int);
int iscsymf(int);
int iscsym(int);

#ifdef __cplusplus
}
#endif

#endif
