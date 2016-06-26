# $ cd <path-to-'haskap-jam-p5'-dir>
# $ LANG=C java -jar ./vendors/ruby-processing-2.6.17/vendors/jruby-complete-1.7.24.jar -e 'load "META-INF/jruby.home/bin/jirb"'
# put the following code to irb input.

require 'psych'
CONFIG_FILE_PATH = './rp5rc'
RP_CONFIG = (Psych.load_file(CONFIG_FILE_PATH))

load './vendors/ruby-processing-2.6.17/lib/ruby-processing.rb'
load './vendors/ruby-processing-2.6.17/lib/ruby-processing/app.rb'

Processing::RP_CONFIG = RP_CONFIG
Processing::App::SKETCH_PATH = defined?(ExerbRuntime) ? ExerbRuntime.filepath : $0

class SonicProcessingLiveSketch < Processing::App
  def setup
    size 200, 200
    background 0
    no_stroke
    smooth
    @rotation = 0
  end

  def draw
    #background 200
    fill 0, 20
    rect 0, 0, width, height

    translate width/2, height/2
    rotate @rotation

    fill 255
    ellipse 0, -60, 20, 20

    @rotation += 0.1
  end
end

SonicProcessingLiveSketch.new(x: 10, y: 30)
