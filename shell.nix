{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs; with xorg; [
      buildPackages.go
      libGL
      libXcursor
      libXinerama
      libXi
      libXrandr
      mesa
      pkg-config
    ];
}
