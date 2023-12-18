{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }: {
    devShell.x86_64-linux = with nixpkgs.legacyPackages.x86_64-linux; mkShell {
      buildInputs = [
        elixir
        elixir_ls
        inotify-tools
        nodejs
        erlang
	sqlite
      ];
    };
  };
}
