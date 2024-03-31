#!/usr/bin/env bash
# text2picture in browser render

set -e

src_dir="$(find $HOME/go/pkg/mod/github.com/blampe_/goat\@*/cmd/goat -maxdepth 1 -type d | tail -n 1)"
dst_dir=/usr/share/nginx/html/
tmp_file=$tmp/goat.wasm

cd $HOME/go/pkg/mod/github.com/blampe/goat@v0.0.0-20220815015552-07bb911fe310/cmd/goat/

GOARCH=wasm \
GOOS=js \
go build -o /tmp/goat.wasm goat_wasm.go

if xxd $tmp_file | head -c 68 | grep -sq ': 0061 736d'
then
    echo "ok: $tmp_file have correct wasm header"
else
    echo "fail: $tmp_file have NO correct wasm header:"
    xxd $tmp_file | head -c 6
fi

mv -v $tmp_file $dst_dir
