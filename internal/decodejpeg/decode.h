#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <setjmp.h>
#include "jpeglib.h"

typedef struct Image
{
  char* pix;
  JDIMENSION width;
  JDIMENSION height;
} Image;

typedef struct JpegErrorManager {
    /* "public" fields */
    struct jpeg_error_mgr pub;
    /* for return to caller */
    jmp_buf setjmp_buffer;
} JpegErrorManager;

struct Image read_JPEG_file(char*, size_t, char*);
