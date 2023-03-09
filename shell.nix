{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = with pkgs; with xorg; [
      buildPackages.go
      libX11
      libXcursor
      libXinerama
      libXrandr
      xorgproto libX11 libXext libXi libXinerama libXrandr
      mesa
      libdrm
      driversi686Linux.mesa
      libGL
      pkg-config
    ];
}
