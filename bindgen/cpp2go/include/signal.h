#ifndef __SIGNAL_H__
#define __SIGNAL_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef int sig_atomic_t;

#define NSIG 23

#define SIGINT 2
#define SIGILL 4
#define SIGABRT_COMPAT 6
#define SIGFPE 8
#define SIGSEGV 11
#define SIGTERM 15
#define SIGBREAK 21
#define SIGABRT 22
#define SIGABRT2 22

typedef void (*__p_sig_fn_t)(int);

#define SIG_DFL (__p_sig_fn_t)0
#define SIG_IGN (__p_sig_fn_t)1
#define SIG_GET (__p_sig_fn_t)2
#define SIG_SGE (__p_sig_fn_t)3
#define SIG_ACK (__p_sig_fn_t)4
#define SIG_ERR (__p_sig_fn_t)-1

__p_sig_fn_t signal(int, __p_sig_fn_t);
int raise(int);

#ifdef __cplusplus
}
#endif

#endif
