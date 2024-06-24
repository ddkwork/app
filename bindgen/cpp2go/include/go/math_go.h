#ifndef __MATH_GO_H__
#define __MATH_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    int                 MaxInt;
    int                 MinInt;
    signed char         MaxInt8;
    signed char         MinInt8;
    short               MaxInt16;
    short               MinInt16;
    int                 MaxInt32;
    int                 MinInt32;
    long long           MaxInt64;
    long long           MinInt64;
    unsigned            MaxUint;
    unsigned char       MaxUint8;
    unsigned short      MaxUint16;
    unsigned            MaxUint32;
    unsigned long long  MaxUint64;

    double E;
    double Pi;
    double Sqrt2;
    double SqrtPi;
    double Ln2;
    double Log2E;
    double Ln10;
    double Log10E;

    float (*Abs)(long double);

    float (*Asinh)(long double);
    float (*Acosh)(long double);
    float (*Atanh)(long double);

    float (*Asin)(long double);
    float (*Acos)(long double);
    float (*Atan)(long double);
    float (*Atan2)(long double, long double);

    float (*Sin)(long double);
    float (*Cos)(long double);
    float (*Tan)(long double);
    void (*Sincos)(long double, double*, double*);

    float (*Sinh)(long double);
    float (*Cosh)(long double);
    float (*Tanh)(long double);

    float (*NaN)(void);
    int (*IsNaN)(double);
    int (*IsInf)(double, int);
    int (*Signbit)(double);

    float (*Sqrt)(long double);
    float (*Cbrt)(long double);

    float (*Copysign)(long double, long double);

    float (*Dim)(long double, long double);
    float (*Max)(long double, long double);
    float (*Min)(long double, long double);

    float (*Erf)(long double);
    float (*Erfc)(long double);

    float (*Exp)(long double);
    float (*Exp2)(long double);
    float (*Expm1)(long double);

    float (*Floor)(long double);
    float (*Ceil)(long double);
    float (*Trunc)(long double);
    float (*Round)(long double);

    float (*FMA)(long double, long double, long double);

    float (*Frexp)(long double, long*);
    float (*Ldexp)(long double, long);

    float (*Gamma)(long double);

    float (*Hypot)(long double, long double);

    float (*Log)(long double);
    float (*Log1p)(long double);
    float (*Log10)(long double);
    float (*Log2)(long double);

    float (*Logb)(long double);
    int (*Ilogb)(long double);

    float (*Mod)(long double, long double);
    float (*Modf)(long double, double*);

    double (*Nextafter)(long double, long double);
    float (*Nextafter32)(float, float);

    float (*Pow)(long double, long double);
    float (*Pow10)(long);

    float (*Remainder)(long double, long double);
} math;

#ifdef __cplusplus
}
#endif

#endif
