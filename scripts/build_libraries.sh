set -x
GOPATH=${GOPATH:=${HOME}/go}
REPO=${GOPATH}/src/github.com/innosat-mats/rac-extract-payload
HOST_WINDOWS=x86_64-w64-mingw32
HOST_DARWIN=x86_64-apple-darwin19
CC_WINDOWS=gcc
CC_DARWIN=cc
SUPPORTED=(linux darwin windows)

function build() {
    tmpdir=$(mktemp -d)
    pushd ${tmpdir}
    tar xf $REPO/third-party/jpegsrc.v9d.tar.gz
    cd jpeg-9d
    case $1 in
    linux)
        ./configure --prefix=${REPO}/third-party/linux
        ;;
    darwin)
        CC=${HOST_DARWIN}-${CC_DARWIN} ./configure --prefix=${REPO}/third-party/darwin --host=${HOST_DARWIN}
        ;;
    windows)
        CC=${HOST_WINDOWS}-${CC_WINDOWS} ./configure --prefix=${REPO}/third-party/windows --host=${HOST_WINDOWS}
        ;;
    esac
    patch -p1 <${REPO}/scripts/12-bit_patch
    make
    make install
    popd
    rm -rf ${tmpdir}
}

if [ $# -eq 0 ]; then
    set -- ${SUPPORTED[@]}
fi

for platform in $@; do
    if [[ "${SUPPORTED[@]}" =~ "${platform}" ]]; then
        build $platform
    fi
done
