class Aliddns < Formula
  desc "Aliyun ddns service"
  homepage "https://github.com/OpenIoTHub"
  url "https://github.com/OpenIoTHub/aliddns.git",
      tag:      "v0.0.5",
      revision: "bebda34d7f36148bc410706063fb9cec6d9ff1df"
  license "MIT"

  depends_on "go" => :build

  def install
    (etc/"aliddns").mkpath
    system "go", "build", "-mod=vendor", "-ldflags",
             "-s -w -X main.version=#{version} -X main.commit=#{stable.specs[:revision]} -X main.builtBy=homebrew",
             *std_go_args
    etc.install "aliddns.yaml" => "aliddns/aliddns.yaml"
  end

  plist_options manual: "aliddns -c #{HOMEBREW_PREFIX}/etc/aliddns/aliddns.yaml"

  def plist
    <<~EOS
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
        <dict>
          <key>Label</key>
          <string>#{plist_name}</string>
          <key>KeepAlive</key>
          <true/>
          <key>ProgramArguments</key>
          <array>
            <string>#{opt_bin}/aliddns</string>
            <string>-c</string>
            <string>#{etc}/aliddns/aliddns.yaml</string>
          </array>
          <key>StandardErrorPath</key>
          <string>#{var}/log/aliddns.log</string>
          <key>StandardOutPath</key>
          <string>#{var}/log/aliddns.log</string>
        </dict>
      </plist>
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aliddns -v 2>&1")
    assert_match "config created", shell_output("#{bin}/aliddns init --config=aliddns.yml 2>&1")
    assert_predicate testpath/"aliddns.yml", :exist?
  end
end
