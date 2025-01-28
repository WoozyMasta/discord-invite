#!/usr/bin/env bash
# require upx
set -eu

: "${WORK_DIR:=${1:-cli}}"
: "${BIN_NAME:=discord-invite}"

build() {
  local GOOS="${1:-linux}" GOARCH="${2:-amd64}" bin

  bin=$BIN_NAME-$GOOS-$GOARCH
  [ "$GOOS" = windows ] && bin+=.exe

  printf 'Build:\t%-10s%-7s' "$GOOS" "$GOARCH"

  CGO_ENABLED=0 GOARCH="$GOARCH" GOOS="$GOOS" \
  GOFLAGS="-buildvcs=false -trimpath" \
    go build -ldflags="-s -w -X '$version' -X '$commit' -X '$date' -X '$url'" \
      -o "./build/$bin" -tags=forceposix "$WORK_DIR/"

  [ "$GOOS" = "windows" ] && GOARCH="$GOARCH" go-winres patch \
    --no-backup --in "winres/winres.json" "./build/$bin"

  cyclonedx-gomod bin -json -output "./build/$bin.sbom.json" "./build/$bin"

  if [ "$GOOS" == "linux" ] && command -v upx &>/dev/null; then
    upx --lzma --best "./build/$bin" > /dev/null
    upx -t "./build/$bin" > /dev/null
  fi

  echo "./build/$bin"
}

mod="$(grep -Po 'module \K.*$' go.mod)"
version="$mod/internal/vars.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 0.0.0)"
commit="$mod/internal/vars.Commit=$(git rev-parse HEAD 2>/dev/null || :)"
date="$mod/internal/vars.BuildTime=$(date -uIs)"
url="$mod/internal/vars.URL=https://$mod"

mkdir -p ./build
go mod tidy

if [ -z "${3:-}" ]; then
  build darwin amd64
  build darwin arm64
  build linux 386
  build linux amd64
  build linux arm
  build linux arm64
  build windows 386
  build windows amd64
  build windows arm64

  cyclonedx-gomod app -json -packages -files -licenses \
    -output "./build/$BIN_NAME.sbom.json" -main "$WORK_DIR"
else
  build "${@:3}"
fi
