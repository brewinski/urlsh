{  
  description = "URL Shortener";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            sqlite
            pkg-config
          ];

          shellHook = ''
            export CGO_ENABLED=1
            export GOPATH=$HOME/go
            export PATH=$PATH:$GOPATH/bin
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "url-shortener";
          version = "0.1.0";
          src = ./.;

          CGO_ENABLED = 1;

          buildInputs = with pkgs; [
            sqlite
          ];

          nativeBuildInputs = with pkgs; [
            pkg-config
          ];

          # Update this hash when your go.mod changes
          # Or set to null and it will be updated automatically
          vendorSha256 = null;
        };
      }
    );
}
