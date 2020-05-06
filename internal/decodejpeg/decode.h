#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "jpeglib.h"

typedef struct Image
{
  char *pix;
  JDIMENSION width;
  JDIMENSION height;
} Image;

struct Image read_JPEG_file(char *, size_t);
