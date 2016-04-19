#--
# Copyright (c) 2015, 2016 SIAF LAB.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#++

def voice_pattern(syllable, notes, *args)
  notes.each { |note| voice(syllable, note, *args); sleep 1 }
end

def voice_pattern_timed(syllable, notes, times, *args)
  if times.is_a?(Array) || times.is_a?(SonicPi::Core::SPVector)
    t = times.ring
    notes.each_with_index { |note, idx| voice(syllable, note, *args); sleep(t[idx]) }
  else
    notes.each { |note| voice(syllable, note, *args); sleep times }
  end
end

def voice_chord(syllable, notes, *args)
  args_h = resolve_synth_opts_hash_or_array(args)
  shifted_notes = notes.map { |n| normalise_transpose_and_tune_note_from_args(n, args_h) }
  use_thread = shifted_notes.size > 3
  if use_thread
    shifted_notes.each { |n| in_thread { voice(syllable, n, *args) } }
  else
    shifted_notes.each { |n| voice(syllable, n, *args) }
  end
end

def voice(syllable, n, *args)
  if n.is_a?(Array) || n.is_a?(SonicPi::Core::RingVector)
    return voice_chord(syllable, n, *args)
  end

  # voice_ba.wav c4 base
  c4_base_note = note(:c4)
  target_note = note(n)
  target_ratio = pitch_to_ratio(target_note - c4_base_note)

  file_base = File.dirname(__FILE__) + '/voice/voice_'
  wav_file_path = "#{file_base}#{syllable}.wav"
  load_sample wav_file_path

  args_h = resolve_synth_opts_hash_or_array(args)
  if args_h.empty?
    # empty
    args_h = { rate: target_ratio } # add rate:
  elsif args_h.key?(:rate)
    # exist rate:
    args_h.store(:rate, target_ratio * args_h[:rate])
  else
    # nil rate:
    args_h.store(:rate, target_ratio)
  end

  ensure_good_timing!
  trigger_sampler wav_file_path, args_h
end
