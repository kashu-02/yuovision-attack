{
  description = "Yuovision attack";

  inputs = {
    nixpkgs = { url = "github:nixos/nixpkgs/nixos-unstable"; };
    flake-utils = { url = "github:numtide/flake-utils"; };
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        inherit (nixpkgs.lib) optional;
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs;[
          	go
            ];

            shellHook = ''
              export GOPATH=$HOME/go
              export GOROOT=${pkgs.go.outPath}/share/go/
              export GOPROXY=https://proxy.golang.org,direct
              export GOSUMDB=sum.golang.org
              export GOTOOLCHAIN=local
            '';
        };
      });
}
