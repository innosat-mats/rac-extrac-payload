#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "jpeglib.h"

GLOBAL(char *)
read_JPEG_file(char *inbuffer, size_t size)
{
  struct jpeg_decompress_struct cinfo;
  struct jpeg_error_mgr jerr;

  JSAMPARRAY buffer; /* Output row buffer */
  char *image;
  int row_stride; /* physical row width in output buffer */

  cinfo.err = jpeg_std_error(&jerr);

  jpeg_create_decompress(&cinfo);

  jpeg_mem_src(&cinfo, inbuffer, size);

  (void)jpeg_read_header(&cinfo, TRUE);

  (void)jpeg_start_decompress(&cinfo);
  image = (char *)malloc(cinfo.output_height * cinfo.output_width * cinfo.output_components);

  row_stride = cinfo.output_width * cinfo.output_components;
  buffer = (*cinfo.mem->alloc_sarray)((j_common_ptr)&cinfo, JPOOL_IMAGE, row_stride, 1);

  while (cinfo.output_scanline < cinfo.output_height)
  {
    (void)jpeg_read_scanlines(&cinfo, buffer, 1);
    memcpy(&image[(cinfo.output_scanline - 1) * row_stride], buffer[0], row_stride);
  }
  (void)jpeg_finish_decompress(&cinfo);

  jpeg_destroy_decompress(&cinfo);

  return image;
}
