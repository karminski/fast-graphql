#include "go_luajit.h"

int amount=0;

void *chmalloc(size_t size) {
    amount++;
    return malloc(size);
}

void increment_allocs() {
    amount++;
}

void chfree(void *pointer) {
    amount--;
    return free(pointer);
}

int outlying_allocs(){
    return amount;
}

void clear_allocs() {
    amount = 0;
}