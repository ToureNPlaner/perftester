#!/usr/bin/gnuplot
set terminal pngcairo size 640,480
set datafile separator ','

set output "threads-ud-ntp.png"
set xrange [1:16]
set title "Scaling by Network Throughput"
set xlabel "# Threads"
set ylabel "Network Throughput in Mbit/s"
plot "threads-new.csv" index 0 using 2:6 title "Up-/Downgraph" with linespoints pointtype 5

set output "threads-sp-ntp.png"
plot "threads-new.csv" index 1 using 2:6 title "Shortest Path" with linespoints pointtype 4

set output "threads-both-ntp.png"
plot "threads-new.csv" index 0 using 2:6 title "Up-/Downgraph" with linespoints pointtype 5, "threads-new.csv" index 1 using 2:6 title "Shortest Path" with linespoints pointtype 4

set output "threads-ud-tp.png"
set title "Scaling by Throughput"
set ylabel "Throughput in Reqs/s"
plot "threads-new.csv" index 0 using 2:5 title "Up-/Downgraph" with linespoints pointtype 5

set output "threads-sp-tp.png"
plot "threads-new.csv" index 1 using 2:5 title "Shortest Path" with linespoints pointtype 4

set output "threads-both-tp.png"
plot "threads-new.csv" index 0 using 2:5 title "Up-/Downgraph" with linespoints pointtype 5, "threads-new.csv" index 1 using 2:5 title "Shortest Path" with linespoints pointtype 4


