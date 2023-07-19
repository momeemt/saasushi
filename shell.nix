{ pkgs ? import ./nix/pkgs.nix }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    go
    gotools
    go-tools
    delve
    gopls
  ];
}
