{ pkgs ? import <nixpkgs> {} }:

# TODO: add shell support
pkgs.callPackage ./release.nix { }

