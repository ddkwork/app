#ifndef __SETJMP_H__
#define __SETJMP_H__

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

typedef int jmp_buf[1];

int setjmp(jmp_buf);
void longjmp(jmp_buf, int);

#ifdef __cplusplus
}
#endif

#endif
