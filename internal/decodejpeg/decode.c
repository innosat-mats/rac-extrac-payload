#include <stdio.h>
#include <setjmp.h>

#include "decode.h"

char jpeg_last_error_message[JMSG_LENGTH_MAX];

void jpeg_error_exit (j_common_ptr cinfo)
{
    /* cinfo->err actually points to a jpeg_error_manager struct */
    JpegErrorManager* myerr = (JpegErrorManager*) cinfo->err;

    /* output_message is a method to print an error message */
    ( *(cinfo->err->output_message) ) (cinfo);

    /* Create the message */
    ( *(cinfo->err->format_message) ) (cinfo, jpeg_last_error_message);

    /* Jump to the setjmp point */
    longjmp(myerr->setjmp_buffer, 1);
}

GLOBAL(struct Image)
read_JPEG_file(char* inbuffer, size_t size, char* error)
{
    struct Image result;
    struct jpeg_decompress_struct cinfo;
    struct JpegErrorManager jerr;

    const int BYTES_PER_SAMPLE = sizeof(JSAMPLE);
    JSAMPARRAY buffer; /* Output row buffer */

    int row_stride; /* physical row width in output buffer */

    cinfo.err = jpeg_std_error(&jerr.pub);
    jerr.pub.error_exit = jpeg_error_exit;
    if (setjmp(jerr.setjmp_buffer)) {
        /* If we get here, the JPEG code has signaled an error. */
        jpeg_destroy_decompress(&cinfo);
        strcpy(error, jpeg_last_error_message);
        return result;
    }

    jpeg_create_decompress(&cinfo);

    jpeg_mem_src(&cinfo, inbuffer, size);

    (void)jpeg_read_header(&cinfo, TRUE);

    (void)jpeg_start_decompress(&cinfo);
    row_stride = cinfo.output_width * cinfo.output_components * BYTES_PER_SAMPLE;
    result.pix = (char *)malloc(cinfo.output_height * row_stride);
    result.width = cinfo.output_width;
    result.height = cinfo.output_height;

    buffer = (*cinfo.mem->alloc_sarray)((j_common_ptr)&cinfo, JPOOL_IMAGE, row_stride, 1);

    while (cinfo.output_scanline < cinfo.output_height)
    {
        (void)jpeg_read_scanlines(&cinfo, buffer, 1);
        memcpy(&(result.pix)[(cinfo.output_scanline - 1) * row_stride],
               buffer[0],
               row_stride);
    }
    (void)jpeg_finish_decompress(&cinfo);

    jpeg_destroy_decompress(&cinfo);

    return result;
}
