#ifndef __FMT_GO_H__
#define __FMT_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

extern struct {
    int (*Sscanf)(const char*, const char*, ...);
    int (*Fscanf)(void*, const char*, ...);
    int (*Scanf)(const char*, ...);

    void* (*Appendf)(char*, const char*, ...);
    int (*Fprintf)(void*, const char*, ...);
    int (*Printf)(const char*, ...);
} fmt;

#ifdef __cplusplus
}
#endif

#endif
