# myapp.rb
require 'sinatra'
require "json"

set :bind, '0.0.0.0'

$foo = 1

# Index html contains dockerstats.js
get '/' do
  send_file File.join(settings.public_folder, 'index.html')
end

get '/containers' do
  ['9078e942cd86','2b58bc125970','fa7d3b7d5d1c'].to_json
end

get '/containers/:cid/stats' do |n|
  if n == 1
    $foo = [[{:cpu => 10}, {:mem => 2}], [{:cpu => 40}, {:mem => 5}]].to_json
  elsif n == 2
    $foo = [[{:cpu => 20}, {:mem => 2}], [{:cpu => 40}, {:mem => 5}]].to_json
  end
  p "foo = #{$foo}"
  [[{:cpu => 50}, {:mem => 2}], [{:cpu => 40}, {:mem => 5}]].to_json
end

error 400..510 do
  'illegal request'
end


