#include <stdlib.h>
#include "object.h"
#include "memory.h"
#include <stdio.h>

void* reallocate (void* ptr, size_t old_size, size_t new_size) {
    if (new_size == 0) {
        free(ptr);
        return NULL;
    }

    void* result = realloc(ptr, new_size);

    if (result == NULL) { 
        fprintf(stderr, "Memory reallocation failed!\n");
        exit(1); 
    }

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

void freeObjects (Obj** objects) {
    for (int i = 0; objects[i] != NULL; i++) {
		freeObject(objects[i]);
	}

    return;
}

