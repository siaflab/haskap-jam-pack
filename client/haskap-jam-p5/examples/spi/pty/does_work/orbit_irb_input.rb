# $ cd <path-to-'haskap-jam-p5'-dir>
# $ LANG=C java -jar ./vendors/ruby-processing-2.6.17/vendors/jruby-complete-1.7.24.jar -e 'load "META-INF/jruby.home/bin/jirb"'
# put the following code to irb input.

require 'psych'
load '~/github/haskap-jam-pack/client/haskap-jam-p5/vendors/ruby-processing-2.6.17/lib/ruby-processing.rb'
load '~/github/haskap-jam-pack/client/haskap-jam-p5/vendors/ruby-processing-2.6.17/lib/ruby-processing/app.rb'

CONFIG_FILE_PATH = '~/github/haskap-jam-pack/client/haskap-jam-p5/rp5rc'
RP_CONFIG = (Psych.load_file(CONFIG_FILE_PATH))

Processing::RP_CONFIG = RP_CONFIG
Processing::App::SKETCH_PATH = defined?(ExerbRuntime) ? ExerbRuntime.filepath : $0

class SonicProcessingLiveSketch < Processing::App
  # Ported from http://nodebox.net/code/index.php/Graphics_State

# This sketch demonstrates how to use the frame rate as orbital state,
# as well as how to use system fonts in Ruby-Processing.
attr_reader :d_font

def setup
  size 450, 450
  frame_rate 30
  smooth
  fill 0
  @d_font = create_font('Helvetica', 40)
end

def draw
  background 255
  translate 225, 225
  text_font d_font
  ellipse 0, 0, 10, 10
  text 'sun', 10, 0
  3.times do |i|
    push_matrix
    rotate frame_count / -180.0 * PI + i * PI / -1.5
    line 0, 0, 120, 0
    translate 120, 0
    ellipse 0, 0, 10, 10
    text_font d_font, 22
    text 'planet', 10, 0
    rotate frame_count / -30.0 * PI
    line 0, 0, 30, 0
    text_font d_font, 15
    text 'moon', 32, 0
    pop_matrix
  end
end

end

SonicProcessingLiveSketch.new()
