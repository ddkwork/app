#ifndef __TIME_H__
#define __TIME_H__

#ifdef __cplusplus
extern "C" {
#endif

#ifndef size_t
#define size_t int
#endif

#ifndef NULL
#ifdef __cplusplus
#define NULL 0
#else
#define NULL ((void *)0)
#endif
#endif

#define CLK_TCK 100

typedef unsigned int clock_t;
typedef unsigned int time_t;

struct tm {
    int tm_sec;   /* seconds after the minute, 0 to 60
                     (0 - 60 allows for the occasional leap second) */
    int tm_min;   /* minutes after the hour, 0 to 59 */
    int tm_hour;  /* hours since midnight, 0 to 23 */
    int tm_mday;  /* day of the month, 1 to 31 */
    int tm_mon;   /* months since January, 0 to 11 */
    int tm_year;  /* years since 1900 */
    int tm_wday;  /* days since Sunday, 0 to 6 */
    int tm_yday;  /* days since January 1, 0 to 365 */
    int tm_isdst; /* Daylight Savings Time flag */
};

extern clock_t clock(void);
extern double difftime(time_t, time_t);
extern time_t mktime(struct tm *);
extern time_t time(time_t *);
extern char *asctime(const struct tm *);
extern char *asctime_r(const struct tm *, char *);
extern char *ctime(const time_t *);
extern struct tm *gmtime(const time_t *);
extern struct tm *localtime(const time_t *);
extern struct tm *localtime_r(const time_t *, struct tm *);
extern size_t strftime(char *, size_t, const char *, const struct tm *);

#ifdef __cplusplus
}
#endif

#endif
