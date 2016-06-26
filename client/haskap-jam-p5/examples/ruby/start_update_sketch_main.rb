#!/usr/bin/ruby
# $ rvm use jruby
# $ cd github/haskap-jam-pack/client/haskap-jam-p5/examples/ruby
# $ ruby start_update_sketch_main.rb

load '../../haskap-jam-p5.rb'

the_code = <<EOC
def setup
  size displayWidth, displayHeight, P2D
  no_stroke
  smooth
  @rotation = 0
end

def draw
  background 0
  #background 200
  fill 0, 20
  rect 0, 0, width, height

  5.times do
    translate rand(width/2), rand(height/2)
    rotate rand(@rotation)

    fill rand(255)
    ellipse rand(60), height/2 -rand(60), rand(100), rand(100)
  end
  @rotation += 1
end
EOC

the_2nd_code = <<EOC
def setup
  size displayWidth, displayHeight, P2D
  no_stroke
  smooth
  @rotation = 0
end

def draw
  #background 0
  background 100
  fill 0, 20
  rect 0, 0, width, height

  5.times do
    translate rand(width/2), rand(height/2)
    rotate rand(@rotation)

    fill rand(255)
    ellipse rand(60), height/2 -rand(60), rand(100), rand(100)
  end
  @rotation += 1
end
EOC

Thread.new do
  rp5_sketch the_code
end
puts '*** sleep while jirb is starting...'
sleep 10
rp5_sketch the_2nd_code
$irb_stdin.close
sleep 1
