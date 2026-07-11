{
  pkgs,
  self,
  stdenv,
}: let
  selfPackages = self.outputs.packages.${stdenv.hostPlatform.system};
  entrypoint = pkgs.writeScript "docker-entrypoint.sh" ''
    #!${pkgs.stdenv.shell}
    set -euo pipefail

    if [ "$#" -gt 0 ]; then
      exec LOGLEVEL=info stock "$@"
      exit 0
    fi

    state_output=$(LOGLEVEL=info stock migrate state 2>&1) || { echo "Failed to check migration state. Exiting."; exit 1; }
    if echo "$state_output" | ${pkgs.gnugrep}/bin/grep -q "Pending"; then
      echo "Database migrations are pending. Running migrations..."
      LOGLEVEL=info stock migrate up || { echo "Database migration failed. Exiting."; exit 1; }
    fi
    exec stock serve
  '';
  envFile = pkgs.writeText "docker-env" (builtins.readFile ./docker.env);
in
  pkgs.dockerTools.streamLayeredImage {
    name = "stock";
    tag = "localdev";

    contents = [
      pkgs.bash
      pkgs.coreutils
      selfPackages.default
    ];

    enableFakechroot = true;
    fakeRootCommands = ''
      mkdir -p /app
      cp ${envFile} /app/.env
    '';

    config = {
      Entrypoint = [entrypoint];
      WorkingDir = "/app";
    };
  }
