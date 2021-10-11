#include "common.h"
#include "chunk.h"
#include "debug.h"

int main(int argc, const char* argv[]) {
  Chunk chunk;
  init_chunk(&chunk);
  write_chunk(&chunk, OP_RETURN);
  dissassemble_chunk(&chunk, "test");
  free_chunk(&chunk);
  return 0;
}