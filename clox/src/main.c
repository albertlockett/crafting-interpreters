#include <stdio.h>

#include "common.h"
#include "chunk.h"
#include "debug.h"
#include "vm.h"

int main(int argc, const char* argv[]) {
  initVM();

  Chunk chunk;
  init_chunk(&chunk);

  int line = 123;

  int constant = add_constant(&chunk, 1.2);
  write_chunk(&chunk, OP_CONSTANT, line);
  write_chunk(&chunk, constant, line);

  constant = add_constant(&chunk, 3.4);
  write_chunk(&chunk, OP_CONSTANT, line);
  write_chunk(&chunk, constant, line);

  write_chunk(&chunk, OP_ADD, line);

  constant = add_constant(&chunk, 5.6);
  write_chunk(&chunk, OP_CONSTANT, line);
  write_chunk(&chunk, constant, line);
  
  write_chunk(&chunk, OP_DIVIDE, line);

  write_chunk(&chunk, OP_NEGATE, line);
  write_chunk(&chunk, OP_RETURN, line);

  dissassemble_chunk(&chunk, "test");
  printf("\n\n");
  interpret(&chunk);
  freeVM();
  free_chunk(&chunk);
  return 0;
}
