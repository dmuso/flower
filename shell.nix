{ pkgs ? import <nixpkgs> {
  config = {
    allowUnfree = true;
  };
} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    air
    go-migrate
    golangci-lint

    bun

    # Overmind runs the API + frontend hot reload processes from Procfile
    overmind

    git
    ripgrep
    gnumake
    scc
  ];

  shellHook = ''
    node() {
      echo "There is no Node [node] here, only Bun. Use [bun] to run your command."
      return 254
    }
    npm() {
      echo "There is no Node [npm] here, only Bun. Use [bun] to run your command."
      return 254
    }
    npx() {
      echo "There is no Node [npx] here, only Bun. Use [bunx] to run your command."
      return 254
    }
    vite() {
      echo "There is no Vite [vite] here, only Bun. Use [bun] to run your command."
      return 254
    }
  '';
}
