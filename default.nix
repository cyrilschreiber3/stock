{
  pkgs,
  buildGoApplication,
}: let
  assets = builtins.fromJSON (builtins.readFile ./assets.json);

  assetDeps =
    builtins.mapAttrs (
      name: asset:
        pkgs.fetchurl {
          url = asset.url;
          sha256 = asset.sha256;
        }
    )
    assets;

  copyCmds = builtins.concatStringsSep "\n    " (
    builtins.map (
      name: let
        asset = assets.${name};
        dep = assetDeps.${name};
      in "cp ${dep} ${asset.destination}"
    ) (builtins.attrNames assets)
  );
in
  buildGoApplication {
    pname = "stock";
    version = "0.1";
    pwd = ./.;
    src = ./.;
    modules = ./gomod2nix.toml;

    preBuild = ''
      mkdir -p ./styles ./static
      ${copyCmds}

      ${pkgs.tailwindcss_4}/bin/tailwindcss -i ./styles/tailwind.css -o ./static/styles.css --minify
      ${pkgs.templ}/bin/templ generate
      ${pkgs.sqlc}/bin/sqlc generate
    '';
  }
