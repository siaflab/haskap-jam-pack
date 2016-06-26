# Welcome to Sonic Pi v2.10

load "~/github/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"

the_code = <<EOC
def setup
  size displayWidth, displayHeight, P2D
  no_stroke
  @rotation = 0
end

def draw
  background 0
  #background 255
  5.times do
    translate rand(width/2), rand(height/2)
    rotate rand(@rotation)
    fill rand(255)
    ellipse rand(60), height/2 -rand(60), rand(100), rand(100)
  end
  @rotation += 1
end
EOC

#start_rp5_sketch the_code, {:full_screen => true}
#stop_rp5_sketch
rp5_sketch the_code, {:full_screen => true}
