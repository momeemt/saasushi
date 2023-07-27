{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
}:
let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.go
    pkgs.gotools
    pkgs.go-tools
    pkgs.delve
    pkgs.gopls
    pkgs.nixpkgs-fmt
    pkgs.actionlint
    goEnv
    pkgs.gomod2nix
    pkgs.niv
  ];
  GOROOT = "${pkgs.go}/share/go";
}
