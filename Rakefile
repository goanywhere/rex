require 'pathname'
require 'fileutils'

OS=RbConfig::CONFIG["host_os"][0..5]
Package = File.join("#{ENV['GOPATH']}", "pkg", "*")


desc 'Remove previously built packages.'
task :clean do
  Pathname.glob("#{Package}").map {|item|
    if item.basename.to_s.start_with?(OS)
      base = File.join(item.to_s, 'github.com', 'goanywhere')
      # remove compiled rex sub-packages
      if File.exists?(File.join(base, 'rex'))
        FileUtils.rm_r File.join(base, 'rex'), :force => true
      end
      # remove compiled rex package
      if File.exists?(File.join(base, 'rex.a'))
        FileUtils.rm File.join(base, 'rex.a'), :force => true
      end
    end
  }
end

desc 'Start building whole rex packages.'
task :build => :clean do
  sh 'go get -v ./...'
end


desc 'Start testing rex packages...'
task 'test' do
  sh 'go test -v ./...'
end
