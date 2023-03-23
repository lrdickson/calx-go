{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs; with xorg; [
      buildPackages.go
      libXcursor
      libXinerama
      libXi
      libXrandr
      mesa
    ];
}
