# Welcome to Sonic Pi v2.10

load "~/github/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"
set_use_pty true
rp5_inline_sketch({:full_screen => true})

attr_reader :d_font

@rotation = 0

def setup
  size displayWidth, displayHeight, P2D
  no_stroke
  smooth
  @rotation = 0
  @d_font = create_font('Helvetica', 40)
end

def draw
  #background 0
  background 255
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
