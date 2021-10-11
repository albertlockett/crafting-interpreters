#include<stdio.h>

#include "debug.h"

void print_value(Value value) {
  printf("%g", value);
}

static int constant_instruction(const char* name, Chunk* chunk, int offset) {
  uint8_t constant = chunk->code[offset+1];
  printf("%-16s %4d '", name, constant);
  print_value(chunk->constants.values[constant]);
  printf("'\n");
  return offset + 2;
}

int simple_instruction(const char* name, int offset) {
  printf("%s\n", name);
  return offset + 1;
}

void dissassemble_chunk(Chunk* chunk, const char* name) {
  printf("== %s ==\n", name);
  for (int offset = 0; offset < chunk->count;) {
    offset = dissassemble_instruction(chunk, offset);
  }
}

int dissassemble_instruction(Chunk* chunk, int offset) {
  printf("%04d ", offset);
  
  uint8_t instruction = chunk->code[offset];
  switch (instruction) {
    case OP_CONSTANT:
      return constant_instruction("OP_CONSTANT", chunk, offset);
    case OP_RETURN:
      return simple_instruction("OP_RETURN", offset);
    default:
      printf("Unkown instruction %d\n", instruction);
      return offset + 1;
  }
}