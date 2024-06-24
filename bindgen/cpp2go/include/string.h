#ifndef __STRING_H__
#define __STRING_H__

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

void *memcpy(void *, const void *, size_t);
void *memmove(void *, const void *, size_t);
char *strcpy(char *, const char *);
char *strncpy(char *, const char *, size_t);
char *strcat(char *, const char *);
char *strncat(char *, const char *, size_t);

int memcmp(const void *, const void *, size_t);
int strcmp(const char *, const char *);
int strncmp(const char *, const char *, size_t);
int strcasecmp(const char *, const char *);
int strncasecmp(const char *, const char *, size_t);

int strcoll(const char *, const char *);
size_t strxfrm(char *, const char *, size_t);

void *memchr(const void *, int, size_t);
char *strchr(const char *, int);

size_t strcspn(const char *, const char *);
char *strpbrk(const char *, const char *);
char *strrchr(const char *, int);
size_t strspn(const char *, const char *);

char *strstr(const char *, const char *);
char *strtok(char *, const char *);

void *memset(void *, int, size_t);
char *strerror(int);
size_t strlen(const char *);

size_t strlcpy(char *, const char *, size_t);
size_t strlcat(char *, const char *, size_t);

#ifdef __cplusplus
}
#endif

#endif
