#ifndef __UNICODE_GO_H__
#define __UNICODE_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    char (*IsUpper)(int);
    char (*IsLower)(int);
    char (*IsTitle)(int);

    char (*ToUpper)(int);
    char (*ToLower)(int);
    char (*ToTitle)(int);
    char (*SimpleFold)(int);

    char (*IsGraphic)(int);
    char (*IsPrint)(int);
    char (*IsControl)(int);
    char (*IsLetter)(int);
    char (*IsMark)(int);
    char (*IsNumber)(int);
    char (*IsPunct)(int);
    char (*IsSpace)(int);
    char (*IsSymbol)(int);

    char (*IsDigit)(int);
} unicode;

#ifdef __cplusplus
}
#endif

#endif
