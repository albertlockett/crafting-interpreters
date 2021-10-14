#ifndef clox_debug_h
#define clox_debug_h

#include "chunk.h"

void print_value(Value value);
void dissassemble_chunk(Chunk* chunk, const char* name);
int dissassemble_instruction(Chunk* chunk, int offset);

#endif