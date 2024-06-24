#ifndef __LOCALE_H__
#define __LOCALE_H__

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

#define LC_ALL 0
#define LC_COLLATE 1
#define LC_CTYPE 2
#define LC_MONETARY 3
#define LC_NUMERIC 4
#define LC_TIME 5

#define LC_MIN LC_ALL
#define LC_MAX LC_TIME

struct lconv {
    char *_;
};

char* setlocale(int, const char*);
struct lconv* localeconv(void);

#ifdef __cplusplus
}
#endif

#endif
