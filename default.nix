{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
          templ.overlays.default
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication {
  pname = "stock";
  version = "0.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  preBuild = ''
    ${pkgs.templ}/bin/templ generate
    ${pkgs.tailwindcss_4}/bin/tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify
  '';
}
