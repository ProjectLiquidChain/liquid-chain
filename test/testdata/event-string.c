#include <stdint.h>
#include <string.h>

typedef Event;
extern Event Say(uint8_t message[]);
extern void chain_arg_size_set(void *, size_t);

int say(int i) {
  char *s = "Checking";
  chain_arg_size_set(s, strlen(s));
  Say("Checking");
  return i;
}
