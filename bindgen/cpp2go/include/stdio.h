#ifndef __STDIO_H__
#define __STDIO_H__

#include <stdarg.h>
#include <go/fmt_go.h>
#include <go/io_go.h>
#include <go/os_go.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef NULL
#ifdef __cplusplus
#define NULL 0
#else
#define NULL ((void *)0)
#endif
#endif

#ifndef size_t
#define size_t int
#endif

typedef GoFile FILE;

#define stdin  os.Stdin
#define stdout os.Stdout
#define stderr os.Stderr

#ifndef EOF
#define EOF io.EOF
#endif

#define SEEK_SET io.SeekStart
#define SEEK_CUR io.SeekCurrent
#define SEEK_END io.SeekEnd

#define STDIN_FILENO    0
#define STDOUT_FILENO   1
#define STDERR_FILENO   2

#define FILENAME_MAX 260
#define FOPEN_MAX 20
#define TMP_MAX 32767

#define fclose(f) f->Close()
#define fflush(f) f->Sync()
#define fopen(path, mode) os.OpenFile(path, *(const int*)(mode), 0)

#define fprintf(f, format, arg...) fmt.Fprintf(f, format, ##arg)
#define printf(format, arg...) fmt.Printf(format, ##arg)
#define sprintf(s, format, arg...) fmt.Appendf(s, format, ##arg)
#define snprintf(s, n, format, arg...) fmt.Appendf(s, format, ##arg)

#define vprintf(format, arg) fmt.Printf(format, arg)
#define vfprintf(f, format, arg) fmt.Fprintf(f, format, arg)
#define vsprintf(s, format, arg) fmt.Appendf(s, format, arg)
#define vsnprintf(s, n, format, arg) fmt.Appendf(s, format, arg)

#define fscanf(f, format, arg...) fmt.Fscanf(f, format, ##arg)
#define scanf(format, arg...) fmt.Scanf(format, ##arg)
#define sscanf(s, format, arg...) fmt.Sscanf(s, format, ##arg)

#define vfscanf(f, format, arg) fmt.Fscanf(f, format, arg)
#define vscanf(format, arg) fmt.Scanf(format, arg)
#define vsscanf(s, format, arg) fmt.Sscanf(s, format, arg)

int fgetc(FILE *);
#define fgets(buf, n, f) f->Read(buf, n)
#define fputc(c, f) fmt.Fprintf(f, "%c", c)
#define fputs(s, f) f->WriteString(s)
#define getc(f) fgetc(f)
#define putc(c, f) fputc(c, f)
int ungetc(int, FILE *);

#define getchar() fgetc(os.Stdin)
#define gets(s) fmt.Scanf("%s", s)
#define putchar(c) fmt.Printf("%c", c)
#define puts(s) os.Stdout->WriteString(s)

#define fread(buf, s, n, f) f->Read(buf, (s) * (n))
#define fwrite(buf, s, n, f) f->Write(buf, (s) * (n))
#define fseek(f, o, w) f->Seek(o, w)
#define ftell(f) f->Seek(0, io.SeekCurrent)
#define rewind(f) f->Seek(0, io.SeekStart)

void clearerr(FILE *);
int feof(FILE *);
int ferror(FILE *);

#ifndef exit
#define exit(code) os.Exit(code)
#endif
#define perror(err) __go_panic(err)

#ifdef __cplusplus
}
#endif

#endif
