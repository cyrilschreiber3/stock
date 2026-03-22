{
  description = "Development environment for Go projects";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.templ.url = "github:a-h/templ";
  inputs.my-packages.url = "github:cyrilschreiber3/nur-packages";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";
  inputs.templ.inputs.nixpkgs.follows = "nixpkgs";
  inputs.templ.inputs.nixpkgs-unstable.follows = "nixpkgs";
  inputs.my-packages.inputs.nixpkgs.follows = "nixpkgs";

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
    templ,
    my-packages,
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          my-packages.overlays.default
        ];
      };
    in {
      packages.default = pkgs.callPackage ./. {
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      };
      devShells.default = pkgs.callPackage ./shell.nix {
        inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
      };
    })
  );
}
