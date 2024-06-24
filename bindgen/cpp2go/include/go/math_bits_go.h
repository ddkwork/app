#ifndef __MATH_BITS_GO_H__
#define __MATH_BITS_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    int (*LeadingZeros)(unsigned long long);
    int (*LeadingZeros8)(unsigned char);
    int (*LeadingZeros16)(unsigned short);
    int (*LeadingZeros32)(unsigned long);
    int (*LeadingZeros64)(unsigned long long);

    int (*TrailingZeros)(unsigned long long);
    int (*TrailingZeros8)(unsigned char);
    int (*TrailingZeros16)(unsigned short);
    int (*TrailingZeros32)(unsigned long);
    int (*TrailingZeros64)(unsigned long long);

    int (*OnesCount)(unsigned long long);
    int (*OnesCount8)(unsigned char);
    int (*OnesCount16)(unsigned short);
    int (*OnesCount32)(unsigned long);
    int (*OnesCount64)(unsigned long long);

    unsigned int (*RotateLeft)(unsigned long long, int);
    unsigned char (*RotateLeft8)(unsigned char, int);
    unsigned short (*RotateLeft16)(unsigned short, int);
    unsigned int (*RotateLeft32)(unsigned long, int);
    unsigned long long (*RotateLeft64)(unsigned long long, int);

    unsigned int (*Reverse)(unsigned long long);
    unsigned char (*Reverse8)(unsigned char);
    unsigned short (*Reverse16)(unsigned short);
    unsigned int (*Reverse32)(unsigned long);
    unsigned long long (*Reverse64)(unsigned long long);

    unsigned int (*ReverseBytes)(unsigned long long);
    unsigned short (*ReverseBytes16)(unsigned short);
    unsigned int (*ReverseBytes32)(unsigned long);
    unsigned long long (*ReverseBytes64)(unsigned long long);

    int (*Len)(unsigned long long);
    int (*Len8)(unsigned char);
    int (*Len16)(unsigned short);
    int (*Len32)(unsigned long);
    int (*Len64)(unsigned long long);
} bits;

#ifdef __cplusplus
}
#endif

#endif
