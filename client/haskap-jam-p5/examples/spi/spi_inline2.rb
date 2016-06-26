# Welcome to Sonic Pi v2.10

load "~/github/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"
#set_use_pty(true)
rp5_inline_sketch({:full_screen => true})

SCALE = 5
COLOR_RANGE = 16581375 # 255 * 255 * 255

attr_reader :grid

def setup
  size(displayWidth, displayHeight, P2D)
  @grid = create_image(width/SCALE, height/SCALE, RGB)
  g.texture_sampling(2)       # 2 = POINT mode sampling
end

def draw
  grid.load_pixels
  cols = width/SCALE
  rows = grid.pixels.length / cols
  rows.times do |i|
    c = rand(COLOR_RANGE)
    ##| c = rand(COLOR_RANGE / 2) + COLOR_RANGE / 2
    ##| c = rand(COLOR_RANGE / 5)
    cols.times do |j|
      ##| c = rand(COLOR_RANGE)
      ##| c = rand(COLOR_RANGE / 2) + COLOR_RANGE / 2
      ##| c = rand(COLOR_RANGE / 5)
      grid.pixels[cols * i + j] = c
    end
  end
  grid.update_pixels
  image(grid, 0, 0, width, height)
end
