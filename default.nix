{ pkgs ? import ./nix/pkgs.nix }:
pkgs.buildGoModule {
  pname = "sushi-as-a-service";
  version = "0.1.0";

  src = ./.;

  vendorHash = null;
}
