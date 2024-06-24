#ifndef __IO_GO_H__
#define __IO_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    int SeekStart;
    int SeekCurrent;
    int SeekEnd;

    int EOF;
} io;

#ifdef __cplusplus
}
#endif

#endif
