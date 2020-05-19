#Compilation
This packages uses a static library from Independent JPEG Group (IJG) http://www.ijg.org/, compiled from release 9d.

Describes how to compile libjpeg for 12-bit color depth. Get the source from http://www.ijg.org/files/jpegsrc.v9d.tar.gz

We use linux as out main development platform. All binaries are compiled on linux (cross compilation)
  GOPATH=${GOPATH:=${HOME}/go}
  DESTINATION=${GOPATH}/src/github.com/innosat-mats/rac-extract-payload/internal/decodejpeg
  tar xf jpegsrc.v9d.tar.gz
  cd jpeg-9d

## Linux
  ./configure --prefix=${DESTINATION}/linux
  patch -p1 <12-bitpatch
  make 
  make install

## Windows

To crosscompile for windows  mingw32 

  CC=x86_64-w64-mingw32-gcc ./configure --prefix=${DESTINATION}/windows

## MacOSX
  CC=o64-clang ./configure --prefix=${GOPATH}/src/github.com/innosat-mats/rac-extract-payload/internal/decodejpeg/darwin --host=x86_64-apple-darwin19
  patch -p1 <12-bit_patch
  make 
  make install