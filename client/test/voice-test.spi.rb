load "~/github/haskap-jam-pack/client/haskap-jam-voice.rb"

voice "ra", [:c4, :a4, :c5]
#voice "ra", chord(:c4, :m13)

#voice "do", :c4
#voice "do", 60

#voice_pattern_timed "ra", chord(:E3, :m7), 0.25
#voice_pattern_timed "ra", scale(:E3, :minor), 0.125, release: 0.1
#voice_pattern_timed "ra", scale(:E3, :minor), [0.125, 0.25, 0.5]

##| loop do
##| voice "ra", choose(chord(:E3, :minor)), release: 0.3
##| sleep 0.25
##| end
