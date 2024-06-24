#ifndef __ENCODING_BINARY_GO_H__
#define __ENCODING_BINARY_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    unsigned short (*Uint16)(const unsigned char*);
    unsigned int (*Uint32)(const unsigned char*);
    unsigned long long (*Uint64)(const unsigned char*);
    void (*PutUint16)(unsigned char*, unsigned short);
    void (*PutUint32)(unsigned char*, unsigned long);
    void (*PutUint64)(unsigned char*, unsigned long long);
} GoByteOrder;

extern struct {
    GoByteOrder LittleEndian;
    GoByteOrder BigEndian;
} binary;

#ifdef __cplusplus
}
#endif

#endif
