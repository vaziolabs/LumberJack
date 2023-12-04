#include <stdlib.h>
#include "./compiler/memory.h"

void* reallocate (void* ptr, size_t old_size, size_t new_size) {
    if (new_size == 0) {
        free(ptr);
        return NULL;
    }

    void* result = realloc(ptr, new_size);

    if (result == NULL) { exit(1); }

    return result;
}

static void freeObject (Obj* object) {
    switch (object->type) {
        case O_STRING: {
            OString* string = (OString*)object;
            FREE_ARRAY(char, string->chars, string->length + 1);
            FREE(OString, object);
            break;
        }
        case O_OPERATION: {
            OOperation* operation = (OOperation*)object;
            //freeValues(&operation->body);
            FREE(OOperation, object);
            break;
        }
        default:
            break;
    }
}

void freeObjects () {
    // TODO: createt objects to free
    return 0;
}

