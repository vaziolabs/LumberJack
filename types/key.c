#include <stdio.h>
#include "key.h"
#include "./compiler/memory.h"

Key *createKey(KeyType type, void *value) {
        Key *key = ALLOCATE(Key, 1);
        key->type = type;
        switch (type) {
                case TYPE_INT:
                        key->value.int_v = *((int *) value);
                        break;
                case TYPE_STRING:
                        key->value.string_v = ALLOCATE(char, strlen((char *) value) + 1);
                        strcpy(key->value.string_v, (char *) value);
                        break;
        }
        return key;
}


void printKey(Key key) {
        switch (key.type) {
                case TYPE_INT:
                        printf("%d", key.value.int_v);
                        break;
                case TYPE_STRING:
                        printf("%s", key.value.string_v);
                        break;
        }
}

void freeKey(Key *key) {
        switch (key->type) {
                case TYPE_INT:
                        FREE(key->type, key);
                        break;
                case TYPE_STRING:
                        FREE_ARRAY(key->type, key->value.string_v, strlen(key->value.string_v) + 1);
                        break;
        }
}
