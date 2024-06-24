#ifndef __OS_GO_H__
#define __OS_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    const char* (*Name)();

    void (*Close)();
    int (*Write)(const void*, long);
    int (*Read)(void*, long);
    int (*Seek)(long, int);
    void (*Sync)();

    int (*WriteString)(const char*);

    int (*Chmod)(long);
    int (*Chown)(int, int);
} GoFile;

extern struct {
    GoFile *Stdin;
    GoFile *Stdout;
    GoFile *Stderr;

    int O_RDONLY;
    int O_WRONLY;
    int O_RDWR;
    int O_APPEND;
    int O_CREATE;
    int O_EXCL;
    int O_SYNC;
    int O_TRUNC;

    GoFile* (*OpenFile)(const char*, int, int);
    int (*Chmod)(const char*, long);
    int (*Chown)(const char*, int, int);
    int (*Mkdir)(const char*, long);
    int (*Rename)(const char*, const char*);
    int (*Remove)(const char*);

    void (*Exit)(int);

    char* (*Getenv)(const char*);
    int (*Setenv)(const char*, const char*);
    int (*Unsetenv)(const char*);
} os;

void __go_panic(const char*);

#ifdef __cplusplus
}
#endif

#endif
