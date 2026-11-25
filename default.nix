# cast-helper for NixOS: nix-build, then result/bin/cast-helper.
# It loads GStreamer at runtime, so the wrapper points it at the store's
# libraries and plugins.
{ pkgs ? import <nixpkgs> { } }:
let
  gst = with pkgs.gst_all_1; [ gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-libav pkgs.pipewire ];
in
pkgs.buildGoModule {
  pname = "cast-helper";
  version = "0.2.0";
  src = pkgs.lib.cleanSource ./.; # leaves out .git and the result link
  vendorHash = pkgs.lib.fakeHash; # placeholder: the first build fails and prints the real hash
  env.CGO_ENABLED = 0;
  doCheck = false; # the tests need a GStreamer setup with VA-API
  nativeBuildInputs = [ pkgs.makeWrapper ];
  postInstall = ''
    wrapProgram $out/bin/cast-helper \
      --prefix LD_LIBRARY_PATH : ${pkgs.lib.makeLibraryPath (gst ++ [ pkgs.glib ])} \
      --prefix GST_PLUGIN_SYSTEM_PATH_1_0 : ${pkgs.lib.makeSearchPathOutput "lib" "lib/gstreamer-1.0" gst}
  '';
}
