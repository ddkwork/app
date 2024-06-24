#ifndef __FENV_H__
#define __FENV_H__

#ifdef __cplusplus
extern "C" {
#endif

#define FE_INVALID      0x01
#define FE_DENORMAL     0x02
#define FE_DIVBYZERO    0x04
#define FE_OVERFLOW     0x08
#define FE_UNDERFLOW    0x10
#define FE_INEXACT      0x20
#define FE_ALL_EXCEPT   0x3F

#define FE_TONEAREST    0x0000
#define FE_DOWNWARD     0x0400
#define FE_UPWARD       0x0800
#define FE_TOWARDZERO   0x0c00

typedef unsigned short fexcept_t;

typedef struct {
    int _;
} fenv_t;

#define FE_DFL_ENV ((const fenv_t *)0)

int feclearexcept(int);
int fegetexceptflag(fexcept_t*, int);
int feraiseexcept(int);
int fesetexceptflag(const fexcept_t*, int);
int fetestexcept(int);

int fegetround(void);
int fesetround(int);

int fegetenv(fenv_t*);
int fesetenv(const fenv_t*);
int feupdateenv(const fenv_t*);
int feholdexcept(fenv_t*);

#ifdef __cplusplus
}
#endif

#endif
