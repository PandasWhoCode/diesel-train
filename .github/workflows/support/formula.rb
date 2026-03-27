class Dl < Formula
  desc "A CLI Tool to drive a Diesel Train Engine across your terminal session"
  homepage "https://github.com/PandasWhoCode/diesel-train"
  url "https://github.com/PandasWhoCode/diesel-train/archive/refs/tags/vSUB_VERSION.tar.gz"
  sha256 "SUB_SHA256"
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args, "-o", bin/"dl"
  end

  test do
    assert_match "dl", shell_output("#{bin}/dl --help 2>&1", 2)
  end
end
