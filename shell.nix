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
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
  templ ? pkgs.templ,
}: let
  goEnv = mkGoEnv {pwd = ./.;};

  prettier-plugin-go-template-patched = pkgs.prettier-plugin-go-template.overrideAttrs (oldAttrs: {
    postInstall =
      oldAttrs.postInstall
      + ''
        if [[ $nodeModulesPath == *prettier-plugin-go-template/node_modules ]]; then
          echo "Fixing node modules location"
          mv $nodeModulesPath/* $out/lib/node_modules/
        fi
      '';
  });

  prettierPluginsPaths = {
    prettier_go_template_plugin_path = "${prettier-plugin-go-template-patched}/lib/node_modules/prettier-plugin-go-template/lib/index.js";
    prettier_tailwindcss_extra_plus_plugin_path = "${pkgs.mypkgs.prettier-plugin-tailwindcss-extra-plus}/lib/node_modules/prettier-plugin-tailwindcss-extra-plus/dist/main.js";
    prettier_tailwindcss_plugin_path = "${pkgs.mypkgs.prettier-plugin-tailwindcss}/lib/node_modules/prettier-plugin-tailwindcss/dist/index.mjs";
  };
in
  with pkgs;
    mkShell {
      packages = [
        # gomod2nix prerequisites
        goEnv
        gomod2nix

        # Go development
        air
        delve
        go
        golangci-lint
        golangci-lint-langserver
        gomodifytags
        gopls
        gotests
        impl

        # Templ
        templ

        # Database
        sqlc
        goose
        postgresql

        # Web development
        nodePackages.prettier
        prettier-plugin-go-template-patched
        mypkgs.prettier-plugin-tailwindcss
        mypkgs.prettier-plugin-tailwindcss-extra-plus
        tailwindcss_4
      ];

      shellHook = ''
        sed -i '/.*prettier.prettierPath*/c\  "prettier.prettierPath\": "${prettier-plugin-go-template-patched}/lib/node_modules/prettier",' ./.vscode/settings.json

        git update-index --assume-unchanged .prettierc
        echo '${builtins.toJSON prettierPluginsPaths}' | ${jinja2-cli}/bin/jinja2 --format=json prettierrc.j2 > .prettierrc

        echo -e "Welcome to the Go dev environment!\n"

        echo -e "$(${go}/bin/go version)\n"

      '';
    }
