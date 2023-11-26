{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils}: let
    flakeForSystem = nixpkgs: system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell = pkgs.mkShell {
        packages = with pkgs; [
          curl
          git
          python3
          virtualenv
          go
          flyctl
          python3Packages.python-lsp-server
          sqlite
        ];
      };
    };
  in
    flake-utils.lib.eachDefaultSystem (system: flakeForSystem nixpkgs system);
}
