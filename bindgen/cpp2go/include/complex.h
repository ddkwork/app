#ifndef __COMPLEX_H__
#define __COMPLEX_H__

#include <go/math_cmplx_go.h>

#define complex _Complex
#define _Complex_I ((float _Complex)__I__)
#define I _Complex_I

#define cabs(x)  cmplx.Abs(x)
#define cabsf(x) (float)cmplx.Abs(x)
#define cabsl(x) cmplx.Abs(x)

#define casin(x)  cmplx.Asin(x)
#define casinf(x) cmplx.Asin(x)
#define casinl(x) cmplx.Asin(x)
#define cacos(x)  cmplx.Acos(x)
#define cacosf(x) cmplx.Acos(x)
#define cacosl(x) cmplx.Acos(x)
#define catan(x)  cmplx.Atan(x)
#define catanf(x) cmplx.Atan(x)
#define catanl(x) cmplx.Atan(x)

#define casinh(x)  cmplx.Asinh(x)
#define casinhf(x) cmplx.Asinh(x)
#define casinhl(x) cmplx.Asinh(x)
#define cacosh(x)  cmplx.Acosh(x)
#define cacoshf(x) cmplx.Acosh(x)
#define cacoshl(x) cmplx.Acosh(x)
#define catanh(x)  cmplx.Atanh(x)
#define catanhf(x) cmplx.Atanh(x)
#define catanhl(x) cmplx.Atanh(x)

#define conj(x)  cmplx.Conj(x)
#define conjf(x) cmplx.Conj(x)
#define conjl(x) cmplx.Conj(x)

#define cexp(x)  cmplx.Exp(x)
#define cexpf(x) cmplx.Exp(x)
#define cexpl(x) cmplx.Exp(x)

#define clog(x)  cmplx.Log(x)
#define clogf(x) cmplx.Log(x)
#define clogl(x) cmplx.Log(x)

#define carg(x)  cmplx.Phase(x)
#define cargf(x) (float)cmplx.Phase(x)
#define cargl(x) cmplx.Phase(x)

#define cpow(x, y)  cmplx.Pow(x, y)
#define cpowf(x, y) cmplx.Pow(x, y)
#define cpowl(x, y) cmplx.Pow(x, y)

#define csin(x)  cmplx.Sin(x)
#define csinf(x) cmplx.Sin(x)
#define csinl(x) cmplx.Sin(x)
#define ccos(x)  cmplx.Cos(x)
#define ccosf(x) cmplx.Cos(x)
#define ccosl(x) cmplx.Cos(x)
#define ctan(x)  cmplx.Tan(x)
#define ctanf(x) cmplx.Tan(x)
#define ctanl(x) cmplx.Tan(x)

#define csinh(x)  cmplx.Sinh(x)
#define csinhf(x) cmplx.Sinh(x)
#define csinhl(x) cmplx.Sinh(x)
#define ccosh(x)  cmplx.Cosh(x)
#define ccoshf(x) cmplx.Cosh(x)
#define ccoshl(x) cmplx.Cosh(x)
#define ctanh(x)  cmplx.Tanh(x)
#define ctanhf(x) cmplx.Tanh(x)
#define ctanhl(x) cmplx.Tanh(x)

#define csqrt(x)  cmplx.Sqrt(x)
#define csqrtf(x) cmplx.Sqrt(x)
#define csqrtl(x) cmplx.Sqrt(x)

#define creal  __go_real
#define crealf __go_real
#define creall __go_real

#define cimag  __go_imag
#define cimagf __go_imag
#define cimagl __go_imag

#ifdef __cplusplus
extern "C" {
#endif

double _Complex cproj(double _Complex);
float _Complex cprojf(float _Complex);
long double _Complex cprojl(long double _Complex);

#ifdef __cplusplus
}
#endif

#endif
