{
  description = "knx-go";

  inputs = {
    flake-utils.url = github:numtide/flake-utils;
    nixpkgs.url = github:NixOS/nixpkgs;
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell = pkgs.mkShell {
        packages =
          # go tooling
          (with pkgs; [
            go
            gopls
          ])
          ++
          # Nix tooling
          (with pkgs; [
            nil
            alejandra
          ]);
      };
    });
}
