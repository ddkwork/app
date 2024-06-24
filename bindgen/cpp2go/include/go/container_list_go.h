#ifndef __CONTAINER_LIST_GO_H__
#define __CONTAINER_LIST_GO_H__

#ifdef __cplusplus
extern "C" {
#endif

typedef struct GoListElement {
    void *Value;
    struct GoListElement* (*Next)();
    struct GoListElement* (*Prev)();
} GoListElement;

typedef struct GoList {
    struct GoList* (*Init)();
    int (*Len)();
    GoListElement* (*Front)();
    GoListElement* (*Back)();
    void* (*Remove)(GoListElement*);
    GoListElement* (*PushFront)(void*);
    GoListElement* (*PushBack)(void*);
    GoListElement* (*InsertBefore)(void*, GoListElement*);
    GoListElement* (*InsertAfter)(void*, GoListElement*);
    void (*MoveToFront)(GoListElement*);
    void (*MoveToBack)(GoListElement*);
    void (*MoveBefore)(GoListElement*, GoListElement*);
    void (*MoveAfter)(GoListElement*, GoListElement*);
    void (*PushBackList)(struct GoList*);
    void (*PushFrontList)(struct GoList*);
} GoList;

extern struct {
    GoList* (*New)();
} list;

#ifdef __cplusplus
}
#endif

#endif
