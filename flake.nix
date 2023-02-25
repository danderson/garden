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
        ];
      };
    };
  in
    flake-utils.lib.eachDefaultSystem (system: flakeForSystem nixpkgs system);
}
# nix-direnv cache busting line: sha256-afcuo/pcLnfFHYTViYi8rPM0ovnUuawuZ26cYhZ1hss= sha256-dhoXBuYV9lE+ssIK4i/TG4cFbzUKSOKnQP47qEEcvsQ=
