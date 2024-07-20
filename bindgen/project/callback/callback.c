#include <stdio.h>

// define callback function pointer
typedef void (*FileDropCallbackFunc)(char** files, int fileCount);

// export callback function pointer
FileDropCallbackFunc dragHandler = NULL;

// export function to set callback function
void SetFileDropCallback(FileDropCallbackFunc fn) {
    dragHandler = fn;
}

// set log buffer in cmd function
void TriggerCallback(char** files, int fileCount) {
    if (dragHandler != NULL) {
        dragHandler(files, fileCount);
    }
}
