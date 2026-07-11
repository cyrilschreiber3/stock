{
  pkgs,
  self,
  stdenv,
}: let
  selfPackages = self.outputs.packages.${stdenv.hostPlatform.system};
  entrypoint = pkgs.writeScript "docker-entrypoint.sh" ''
    #!${pkgs.stdenv.shell}
    set -e

    if [ "$#" -gt 0 ]; then
      exec stock "$@"
      exit 0
    fi

    if stock migrate state | ${pkgs.gnugrep}/bin/grep -q "Pending"; then
      echo "Database migrations are pending. Running migrations..."
      stock migrate up || { echo "Database migration failed. Exiting."; exit 1; }
    fi
    exec stock serve
  '';
in
  pkgs.dockerTools.streamLayeredImage {
    name = "stock";
    tag = "localdev";

    contents = [
      pkgs.bash
      pkgs.coreutils
      selfPackages.default
    ];

    config = {
      Entrypoint = [entrypoint];
    };
  }
