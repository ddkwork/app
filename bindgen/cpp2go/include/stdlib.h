#ifndef __STDLIB_H__
#define __STDLIB_H__

#include <go/os_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef size_t
#define size_t int
#endif

#ifndef NULL
#ifdef __cplusplus
#define NULL 0
#else
#define NULL ((void *)0)
#endif
#endif

#ifndef wchar_t
typedef int rune;
#define wchar_t rune
#endif

#ifndef RAND_MAX
 #define RAND_MAX  32767u
#endif

double atof(const char *);
int atoi(const char *);
long int atol(const char *);
long long atoll(const char *);

double strtod(const char *, char **);
float strtof(const char *, char **);
long double strtold(const char *, char **);
long int strtol(const char *, char **, int);
unsigned long strtoul (char *, char **, unsigned char);
long long strtoll(const char *, char **, int);
unsigned long long strtoull(const char *, char **, int);

int rand(void);
void srand(unsigned int);

void *calloc(size_t, size_t);
void free(void *);
void *malloc(size_t);
void *realloc(void *, size_t);

#ifndef exit
#define exit(code) os.Exit(code)
#endif
#ifndef abort
#define abort() __go_panic("")
#endif
int atexit(void (*)(void));

#define getenv(env) os.Getenv(env)
int  system(const char *);

void *bsearch(const void *, const void *,
    size_t, size_t, int (*)(const void *, const void *));
void qsort(void *, size_t, size_t,
    int (*)(const void *, const void *));

int mblen(const char *, size_t);
int mbtowc(wchar_t *, const char *, size_t);
int wctomb(char *, wchar_t);
size_t mbstowcs(wchar_t *, const char *, size_t);
size_t wcstombs(char *, const wchar_t *, size_t);

#ifdef __cplusplus
}
#endif

#endif
