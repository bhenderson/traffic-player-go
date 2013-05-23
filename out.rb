require 'timeout'

srand 1
pre = ('a'..'z').to_a

$stdout.sync = true

begin
  Timeout.timeout(1) do
    while true do
      # sleep 0.01
      puts "#{pre.sample*4}:This is my message"
    end
  end
rescue TimeoutError
  pre.pop
  retry unless pre.empty?
end

=begin

// 10% every 1sec
sec | count | out
0   | 15    | 1
1   | 15    | 1
2   | 15    | 1
3   | 15    | 1
4   | 15    | 1
5   | 15    | 1

// 10%
sec | count | out
0   | 15    | 1
1   | 15    | 2
2   | 15    | 1
3   | 15    | 2
4   | 15    | 1
5   | 15    | 2

=end
