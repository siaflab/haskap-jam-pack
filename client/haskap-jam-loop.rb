#--
# Copyright (c) 2015, 2016 SIAF LAB.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#++

# require 'socket'
# require_relative 'oscencode'
# require_relative 'udp_client'

load File.expand_path(File.dirname(__FILE__) + '/haskap-jam-util.rb')
load File.expand_path(File.dirname(__FILE__) + '/tcp_client.rb')

include HaskapJam::Log
include HaskapJam::Util
set_log_level :debug

def local_address
  udp = UDPSocket.new
  udp.connect('128.0.0.0', 7)
  local_addr = Socket.unpack_sockaddr_in(udp.getsockname)[1]
  udp.close
  log_debug "local_address: #{local_addr}"
  if local_addr.nil? || local_addr.empty?
    msg = 'local_addr is nil or empty!'
    log_error msg
    fail msg
  end
  local_addr
end

def fetch(hash, key, default)
  (hash.nil? || hash[key].nil?) ? default : hash[key]
end

# entry point
$jam_loop_invoked = false
def jam_loop(_loop_symbol, &proc)
  return if $jam_loop_invoked
  $jam_loop_invoked = true

  # get local address
  local_addr = local_address
  log_debug "local_address: #{local_addr}"
  puts "local_address: #{local_addr}"

  # read config file
  config = eval File.read(File.dirname(__FILE__) + '/haskap-jam-config.rb')
  remote_address = fetch(config, :remote_address, '127.0.0.1')
  log_debug "remote_address: #{remote_address}"
  puts "remote_address: #{remote_address}"
  remote_port = fetch(config, :remote_port, 4557)
  log_debug "remote_port: #{remote_port}"
  puts "remote_port: #{remote_port}"

  # read workspace code
  log_debug "proc.source_location: #{proc.source_location}"
  source_file_name = proc.source_location[0]
  workspace_id = extract_workspace_id(source_file_name)
  code = read_workspace(workspace_id)

  # replace 'jam_loop' with 'live_loop'
  remote_code = code.gsub('jam_loop', 'live_loop')
  # commentout 'load "haskap-jam-loop.rb"'
  remote_code = remote_code.gsub(/^(load.*haskap-jam-loop.rb.*$)/, '#\\1')
  log_debug "remote_code: #{remote_code}"
  puts "remote_code: #{remote_code}"

  # send code to remote sonic pi
  client = HaskapJam::OSC::TCPClient.new(remote_address, remote_port, use_encoder_cache: true)
  client_id = "haskap-client-#{local_addr}"
  log_debug "client_id: #{client_id}"
  buffer_id = "haskap-buffer-#{local_addr}-#{workspace_id}"
  log_debug "buffer_id: #{buffer_id}"
  workspace = "haskap-workspace-#{local_addr}-#{workspace_id}"
  log_debug "workspace: #{workspace}"
  client.send('/save-and-run-buffer', client_id, buffer_id, remote_code, workspace)
  log_debug "done."
end
