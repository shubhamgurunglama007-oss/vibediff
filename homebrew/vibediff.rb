# VibeDiff Homebrew Formula
# Install with: brew install vibediff/vibediff/vibediff

class Vibediff < Formula
  desc "Git-native versioning layer for prompts and AI outputs"
  homepage "https://github.com/shubhamgurunglama007-oss/vibediff"
  url "https://github.com/shubhamgurunglama007-oss/vibediff/archive/refs/tags/v1.0.0.tar.gz"
  sha256 :no_check
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/vibediff"
  end

  test do
    system bin/"vibediff", "version"
  end
end
