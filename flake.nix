{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    templ = {
      url = github:a-h/templ;
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, templ }: {
    devShell.x86_64-linux = with nixpkgs.legacyPackages.x86_64-linux; mkShell {
      buildInputs = [
        elixir
        elixir_ls
        inotify-tools
        nodejs
        erlang
	      sqlite
        flyctl
        go
        gopls
        sqlc
        templ.packages.x86_64-linux.templ
      ];
    };

    packages.x86_64-linux.default = with nixpkgs.legacyPackages.x86_64-linux; let
      mix = "${pkgs.elixir}/bin/mix";
    in stdenv.mkDerivation {
      name = "garden";
      src = ./.;

      nativeBuildInputs = [
        elixir
        erlang
        nodejs
        cacert
      ];
      __noChroot = true;

      buildPhase = ''
        export HOME=`pwd`
        export MIX_ENV=prod
        ${mix} deps.get --only prod
        ${mix} compile
        ${mix} assets.deploy
        ${mix} phx.gen.release
        ${mix} release --path $out
      '';
    };
  };
}
