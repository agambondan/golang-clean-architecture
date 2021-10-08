go build -o build/main main.go

./build/main

#build() {
#  go build -o build/main main.go
#}
#
#run() {
#  build &
#  grep golang
#  wait $! || {
#    exit_code=$?
#    echo '--> cleanup'
#    return $exit_code
#  }
#}
#
#run
#set -e
#
#outer() {
#  echo '--> outer'
#  inner &
#  wait $! || {
#    exit_code=$?
#    echo '--> cleanup'
#    return $exit_code
#  }
#  echo '<-- outer'
#}
#
#inner() {
#  set -e
#  echo '--> inner'
#  some_failed_command
#  echo '<-- inner'
#}
#
#outer
