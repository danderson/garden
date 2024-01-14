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
	      sqlite-interactive
        flyctl
        go
        gopls
        (buildGoModule {
          pname = "sqlc";
          version = "1.25.0";

          src = fetchFromGitHub {
            owner = "sqlc-dev";
            repo = "sqlc";
            rev = "v1.25.0";
            hash = "sha256-VrR/oSGyKtbKHfQaiLQ9oKyWC1Y7lTZO1aUSS5bCkKY=";
          };

          vendorHash = "sha256-nzqRP4U/VpfABCFFSV9KivGPuTL7u3StW09YO/QMD1Q=";
          CGO_ENABLED = 0;

          subPackages = ["cmd/sqlc"];
        })
        templ.packages.x86_64-linux.templ
        tailwindcss
        wgo
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
