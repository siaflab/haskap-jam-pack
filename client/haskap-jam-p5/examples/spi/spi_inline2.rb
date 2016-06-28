# Welcome to Sonic Pi v2.10

load "~/github/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"
rp5_inline_sketch({:full_screen => true})

SCALE = 50
COLOR_RANGE = 16581375 # 255 * 255 * 255

attr_reader :grid

def setup
  size(displayWidth, displayHeight, P2D)
  @grid = create_image(width/SCALE, height/SCALE, RGB)
  g.texture_sampling(2)       # 2 = POINT mode sampling
end

def draw
  background 0
  grid.load_pixels
  cols = width/SCALE
  rows = grid.pixels.length / cols
  rows.times do |i|
    cols.times do |j|
      c = rand COLOR_RANGE
      grid.pixels[cols * i + j] = c
    end
  end
  grid.update_pixels
  image(grid, 0, 0, width, height)
end
