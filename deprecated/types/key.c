#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include "key.h"
#include "memory.h"


void initKeys(Keys *keys) {
		keys->count = 0;
		keys->capacity = 0;
		keys->key = NULL;
}

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

void addKey(Keys *keys, Key *key) {
		if (keys->count + 1 > keys->capacity) {
				int old_capacity = keys->capacity;
				keys->capacity = EXPAND_CAPACITY(old_capacity);
				keys->key = EXPAND(Key, keys->key, old_capacity, keys->capacity);
		}
		
        keys->key[keys->count] = key;
        keys->count++;
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

void destroyKey(Key *key) {
        switch (key->type) {
                case TYPE_INT:
                        FREE(key->type, key);
                        break;
                case TYPE_STRING:
                        FREE_ARRAY(key->type, key->value.string_v, strlen(key->value.string_v) + 1);
                        break;
        }
}

void destroyKeys(Keys *keys) {
		for (int i = 0; i < keys->count; i++) {
				destroyKey(keys->key[i]);
		}

		FREE_ARRAY(Key, keys->key, keys->capacity);
		initKeys(keys);
}