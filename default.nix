with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "debateframeEnv";

  src = null;

  buildInputs = [
      go_1_11
      nodejs-11_x
      binaryen # for wasm-opt
  ];

    shellHook = ''
    export GOPATH=$(pwd)/gopath
    export PATH=$GOPATH/bin:$PATH
    go get github.com/shurcooL/go-goon
    go get github.com/shurcooL/goexec
    mkdir dist
    npm i
  '';
}
