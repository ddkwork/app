#ifndef UNTITLED_LIBRARY_H
#define UNTITLED_LIBRARY_H

__declspec(dllexport) void hello();

typedef  unsigned long long UINT64;

typedef struct _CR3_TYPE
{
    union
    {
        UINT64 Flags;

        struct
        {
            UINT64 Pcid : 12;
            UINT64 PageFrameNumber : 36;
            UINT64 Reserved1 : 12;
            UINT64 Reserved_2 : 3;
            UINT64 PcidInvalidate : 1;
        } Fields;
    };
} CR3_TYPE, *PCR3_TYPE;

#endif //UNTITLED_LIBRARY_H
