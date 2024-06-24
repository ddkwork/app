#ifndef __STDDEF_H__
#define __STDDEF_H__

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

typedef long ptrdiff_t;

#define offsetof(t, d) __builtin_offsetof(t, d)

#endif
