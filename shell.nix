{ pkgs ? (let
  inherit (builtins) fetchTree fromJSON readFile;
  inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
in import (fetchTree nixpkgs.locked) {
  overlays = [ (import "${fetchTree gomod2nix.locked}/overlay.nix") ];
}), mkGoEnv ? pkgs.mkGoEnv, gomod2nix ? pkgs.gomod2nix }:

let goEnv = mkGoEnv { pwd = ./.; };
in pkgs.mkShell {
  packages = with pkgs; [
    goEnv
    gomod2nix

    # Dev tools
    golangci-lint
    nixfmt
    nodejs_21
  ];
}
