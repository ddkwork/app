#ifndef __MATH_CMPLX_GO_H__
#define __MATH_CMPLX_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    float (*Abs)(long double _Complex);

    float _Complex (*Asin)(long double _Complex);
    float _Complex (*Acos)(long double _Complex);
    float _Complex (*Atan)(long double _Complex);

    float _Complex (*Asinh)(long double _Complex);
    float _Complex (*Acosh)(long double _Complex);
    float _Complex (*Atanh)(long double _Complex);

    float _Complex (*Sin)(long double _Complex);
    float _Complex (*Cos)(long double _Complex);
    float _Complex (*Tan)(long double _Complex);

    float _Complex (*Sinh)(long double _Complex);
    float _Complex (*Cosh)(long double _Complex);
    float _Complex (*Tanh)(long double _Complex);

    float _Complex (*Conj)(long double _Complex);

    float _Complex (*Exp)(long double _Complex);

    float _Complex (*Log)(long double _Complex);

    float (*Phase)(long double _Complex);

    float _Complex (*Pow)(long double _Complex, long double _Complex);

    float _Complex (*Sqrt)(long double _Complex);
} cmplx;

float __go_real(long double _Complex);
float __go_imag(long double _Complex);

#ifdef __cplusplus
}
#endif

#endif
