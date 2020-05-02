#include "decode.h"

GLOBAL(struct Image)
read_JPEG_file(char *inbuffer, size_t size)
{
  struct Image result;
  struct jpeg_decompress_struct cinfo;
  struct jpeg_error_mgr jerr;

  const int BYTES_PER_SAMPLE = sizeof(JSAMPLE);
  JSAMPARRAY buffer; /* Output row buffer */

  int row_stride; /* physical row width in output buffer */

  cinfo.err = jpeg_std_error(&jerr);

  jpeg_create_decompress(&cinfo);

  jpeg_mem_src(&cinfo, inbuffer, size);

  (void)jpeg_read_header(&cinfo, TRUE);

  (void)jpeg_start_decompress(&cinfo);
  result.pix = (char *)malloc(cinfo.output_height * cinfo.output_width * cinfo.output_components * BYTES_PER_SAMPLE);
  result.width = cinfo.output_width;
  result.height = cinfo.output_height;

  row_stride = cinfo.output_width * cinfo.output_components * BYTES_PER_SAMPLE;
  buffer = (*cinfo.mem->alloc_sarray)((j_common_ptr)&cinfo, JPOOL_IMAGE, row_stride, 1);

  while (cinfo.output_scanline < cinfo.output_height)
  {
    (void)jpeg_read_scanlines(&cinfo, buffer, 1);
    memcpy(&(result.pix)[(cinfo.output_scanline - 1) * row_stride], buffer[0], row_stride * BYTES_PER_SAMPLE);
  }
  (void)jpeg_finish_decompress(&cinfo);

  jpeg_destroy_decompress(&cinfo);

  return result;
}
