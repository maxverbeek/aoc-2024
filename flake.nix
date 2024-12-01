{
  inputs = {
    nixpkgs.url = "nixpkgs";
    utils.url = "github:numtide/flake-utils";
    gocount.url = "github:maxverbeek/gocounter";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
      gocount,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gocount.overlays.${system}.default ];
        };

        inputgetter = pkgs.writeScriptBin "getinput" ''
          outpath=''${1:-input.txt}
          daynum=$(expr $(basename $PWD) + 0)

          url="https://adventofcode.com/2024/day/''${daynum}/input"

          echo $url

          exec ${pkgs.curl}/bin/curl -o $outpath -H "cookie: $(cat ../.cookies.txt)" $url
        '';
      in
      {
        devShell = pkgs.mkShell {
          name = "devshell";
          buildInputs = with pkgs; [
            inputgetter
            go
            pkgs.gocount
          ];
        };
      }
    );
}
