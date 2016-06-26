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
module HaskapJam
  module Log
    extend self

    LOGLEVEL_DEBUG = 0 unless defined?(LOGLEVEL_DEBUG)
    LOGLEVEL_INFO = 1 unless defined?(LOGLEVEL_INFO)
    LOGLEVEL_WARN = 2 unless defined?(LOGLEVEL_WARN)
    LOGLEVEL_ERROR = 3 unless defined?(LOGLEVEL_ERROR)
    LOGLEVEL_FATAL = 4 unless defined?(LOGLEVEL_FATAL)
    LOGLEVEL_MAP = { debug: LOGLEVEL_DEBUG, info: LOGLEVEL_INFO,
                     warn: LOGLEVEL_WARN, error: LOGLEVEL_ERROR,
                     fatal: LOGLEVEL_FATAL }.freeze unless defined?(LOGLEVEL_MAP)

    @@log_level = LOGLEVEL_MAP[:debug]

    def set_log_level(level)
      lv = LOGLEVEL_MAP[level]
      @@log_level = lv unless lv.nil?
    end

    def log_debug(msg)
      return if @@log_level > LOGLEVEL_DEBUG
      return unless defined? SonicPi
      log_path = SonicPi::Util.log_path
      File.open("#{log_path}/debug.log", 'a') do |f|
        f.write("[#{Time.now.strftime('%Y-%m-%d %H:%M:%S')}|haskap-jam] #{msg}\n")
      end
    end

    def log_info(msg)
      return if @@log_level > LOGLEVEL_INFO
      puts "[#{Time.now.strftime('%Y-%m-%d %H:%M:%S')}|haskap-jam] #{msg}"
    end

    def log_error(msg)
      return if @@log_level > LOGLEVEL_ERROR
      STDERR.puts "[#{Time.now.strftime('%Y-%m-%d %H:%M:%S')}|haskap-jam] #{msg}"
    end
  end

  module Util
    extend self
    include HaskapJam::Log

    NUMBER_NAMES = %w(zero one two three four five six
                      seven eight nine).freeze unless defined?(NUMBER_NAMES)

    def read_workspace(workspace_id)
      workspace_file_path = workspace_filepath(workspace_id)
      log_debug "workspace_file_path: #{workspace_file_path}"
      code = read_file(workspace_file_path, 5)
      log_debug "code: #{code}"
      if code.nil? || code.empty?
        msg = 'code is nil or empty!'
        log_error msg
        raise msg
      end
      code
    end

    def extract_workspace_id(source_file_name)
      # source_file_name: "Workspace_one"
      matched = source_file_name.match(/([a-zA-Z]|\s)?_([a-zA-Z]*)$/)
      if matched.nil? || matched[2].nil?
        msg = "source_file_name not matched. source_file_name: #{source_file_name}"
        log_error msg
        raise msg
      end
      workspace_id = matched[2] # one
      log_debug("workspace_id: #{workspace_id}")
      workspace_id
    end

    def workspace_filepath(workspace_id)
      file_name = workspace_filename(workspace_id)
      log_debug "file_name: #{file_name}"
      if file_name.nil? || file_name.empty?
        msg = 'file_name is nil or empty!'
        log_error msg
        raise msg
      end

      project_path = SonicPi::Util.project_path
      log_debug "project_path: #{project_path}"
      if project_path.nil? || project_path.empty?
        msg = 'project_path is nil or empty!'
        log_error msg
        raise msg
      end

      project_path + file_name
    end

    def workspace_filename(workspace_id)
      # workspace_id: zero, one, ... nine
      'workspace_' + workspace_id + '.spi'
    end

    def number_name(i)
      if i < 0
        msg = "can not convert to number name: #{i}"
        log_error msg
        raise msg
      end
      name = NUMBER_NAMES.fetch(i, nil)
      if name.nil?
        msg = "can not convert to number name: #{i}"
        log_error msg
        raise msg
      end
      name
    end

    def read_file(file_path, max_retry)
      code = nil
      rep = 1
      while code.nil? || code.empty? || rep < max_retry
        code = File.read(file_path)
        rep += 1
      end
      code
    end
  end
end
