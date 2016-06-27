# ruby-processing (experimental)

"ruby-processing" feature enables you to control ruby-processing visuals from Sonic Pi.

## Requirements

* Sonic Pi v2.10
* Processing v2.2.1

## Configuration
Set Processing2.2.1.app installed path in rp5rc file.
```
PROCESSING_ROOT: "/Applications/Processing2.2.1.app/Contents/Java"  # Path for Mac
```

## Usage
### Inline Mode
#### Run
With "Inline Mode", you can directly run "ruby-processing" code on Sonic Pi.

* Open Sonic Pi.
* Include `load "[path to haskap-jam-p5.rb]/haskap-jam-p5.rb"` in the code.
* Include `rp5_inline_sketch` in the code.
* Put "ruby-processing" code.
* Run the code.

```ruby
load "~/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"
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
```

#### Stop
Switch to "ruby-processing" sketch window (probably named 'java') and quit (cmd + q).

### Manual Mode
You can manually run "ruby-processing" code with rp5_sketch method.

```ruby
load "~/haskap-jam-pack/client/haskap-jam-p5/haskap-jam-p5.rb"

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

the_2nd_code = <<EOC
def setup
  size displayWidth, displayHeight, P2D
  no_stroke
  @rotation = 0
end

def draw
  background 255
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

live_loop :sketch do
  rp5_sketch the_code, {:full_screen => true}
  sleep 5
  rp5_sketch the_2nd_code, {:full_screen => true}
  sleep 5
end
```

#### Stop
Add `stop_rp5_sketch` and remove `rp5_sketch`.

```ruby
live_loop :sketch do
  stop_rp5_sketch
  #  rp5_sketch the_code, {:full_screen => true}
  sleep 5
  #  rp5_sketch the_2nd_code, {:full_screen => true}
  sleep 5
end


```
