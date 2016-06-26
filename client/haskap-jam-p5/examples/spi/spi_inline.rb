# Welcome to Sonic Pi v2.10

load "~/github/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"
rp5_inline_sketch({:full_screen => true})

def setup
  size displayWidth, displayHeight, P2D
  no_stroke
end

def draw
  ##| background 0
  background 255
  5.times do
    translate rand(width/2), rand(height/2)
    rotate rand(100)
    fill rand(255) 
    ellipse rand(60), height/2 -rand(60), rand(100), rand(100)
  end
end
