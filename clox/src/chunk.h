#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include "chunk.h"

typedef enum {
  OP_RETURN
} OpCode;

typedef struct {
  int count;
  int capacity;
  uint8_t* code;
} Chunk;

void free_chunk(Chunk* chunk);
void init_chunk(Chunk* chunk);
void write_chunk(Chunk* chunk, uint8_t byte);

#endif