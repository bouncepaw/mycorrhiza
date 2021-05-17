{ stdenv, lib, fetchFromGitHub, buildGoModule
, git
}:

buildGoModule rec {
  pname = "mycorrhiza";
  version = "1.1.0";

  src = ./.;

  # That's for a nixpkgs release or something. Mind the hashes.
  # src = fetchFromGitHub {
  #   owner = "bouncepaw";
  #   repo = "mycorrhiza";
  #   rev = "v${version}";
  #   sha256 = "0di4msrl44jcnhck11k0r7974cjnwdyw45b3hf0s3dbwx6ijdkdd";
  #   fetchSubmodules = true;
  # };

  vendorSha256 = "0hxcbfh55avly9gvdysqgjzh66g7rdy2l0wmy9jnlq0skpa6j0jq"; 
  
  subPackages = [ "." ]; 

  propagatedBuildInputs = [ git ];

  meta = with lib; {
    description = "Filesystem and git-based wiki engine written in Go using mycomarkup as its primary markup language";
    homepage = "https://github.com/bouncepaw/mycorrhiza";
    license = licenses.agpl3;
    # maintainers = with maintainers; [ bouncepaw ];
    platforms = platforms.linux;
  };
}
