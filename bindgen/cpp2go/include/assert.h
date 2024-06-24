#ifndef __ASSERT_H__
#define __ASSERT_H__

#include <go/os_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#undef static_assert
#define static_assert _Static_assert

#undef assert
#ifndef NDEBUG
#define assert(expr) if (!(expr)) __go_panic((#expr));
#else
#define assert(expr)
#endif

#ifndef exit
#define exit(code) os.Exit(code)
#endif
#ifndef abort
#define abort() __go_panic("")
#endif

#ifdef __cplusplus
}
#endif

#endif
