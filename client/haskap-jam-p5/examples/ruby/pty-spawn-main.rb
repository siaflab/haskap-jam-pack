#!/usr/bin/ruby
# $ rvm use 2.3
# $ cd github/haskap-jam-pack/client/haskap-jam-p5/examples/ruby/
# $ ruby pty-spawn-main.rb
require 'pty'
require 'expect'

commands = ["puts 'hi'", "a = 1", "puts a"]
cmd = "LANG=C java -jar ../../vendors/ruby-processing-2.6.17/vendors/jruby-complete-1.7.24.jar -e 'load \"META-INF/jruby.home/bin/jirb\"'"

PTY.spawn(cmd) do |r, w, pid|
  # r is node's stdout/stderr and w is stdin
  r.expect(/(.*)>/m)
  commands.each do |cmd|
    puts "Command: " + cmd
    w.puts cmd
    r.expect(/(.*?)\r\n(.*)\r\n(.*)>/m) { |res|
      puts "Output: " + res[2]
    }
  end
end
