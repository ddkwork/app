#include <stdio.h>

// 定义回调函数类型
typedef void (*FileDropCallbackFunc)(char** files, int fileCount);

// 全局变量存储回调函数
FileDropCallbackFunc dragHandler = NULL;

// 设置回调函数
void SetFileDropCallback(FileDropCallbackFunc fn) {
    dragHandler = fn;
}

// 调用回调函数
void TriggerCallback(char** files, int fileCount) {
    if (dragHandler != NULL) {
        dragHandler(files, fileCount);
    }
}
