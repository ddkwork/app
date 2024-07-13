#define MAX_PATH 260
#define WCHAR_MAX   65535
typedef unsigned long long UINT64;
typedef void (*SymbolMapCallback)(unsigned long long Address, char * ModuleName, char * ObjectName, unsigned int ObjectSize);
//typedef BOOLEAN (*LOG_CALLBACK_PREPARE_AND_SEND_MESSAGE_TO_QUEUE)(UINT32       OperationCode,
//                                                                  BOOLEAN      IsImmediateMessage,
//                                                                  BOOLEAN      ShowCurrentSystemTime,
//                                                                  BOOLEAN      Priority,
//                                                                  const char * Fmt,
//                                                                  va_list      ArgList);
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