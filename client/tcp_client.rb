#--
# Copyright (c) 2016 SIAF LAB.
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

# created from
#  https://github.com/samaaron/sonic-pi/blob/v2.10.0/app/server/sonicpi/lib/sonicpi/osc/tcp_client.rb
#  and https://github.com/samaaron/sonic-pi/blob/v2.10.0/app/server/sonicpi/lib/sonicpi/osc/udp_client.rb

require 'socket'

module HaskapJam
  module OSC
    class TCPClient

      def initialize(host, port, opts={})
        @host = host
        @port = port
        @opts = opts
        use_encoder_cache = opts[:use_encoder_cache]
        encoder_cache_size = opts[:encoder_cache_size] || 1000
        @encoder = SonicPi::OSC::OscEncode.new(use_encoder_cache, encoder_cache_size)
      end

      def send(pattern, *args)
        so = TCPSocket.new(@host, @port)
        msg = @encoder.encode_single_message(pattern, args)
        so.send(msg, 0)
        so.close
      end

      def send_ts(ts, pattern, *args)
        so = TCPSocket.new(@host, @port)
        msg = @encoder.encode_single_bundle(ts, pattern, args)
        so.send(msg, 0)
        so.close
      end

      def to_s
        "#<HaskapJam::OSC::TCPClient host: #{@host}, port: #{@port}, opts: #{@opts.inspect}>"
      end

      def inspect
        to_s
      end

    end
  end
end
