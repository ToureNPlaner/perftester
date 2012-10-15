#!/bin/zsh
echo '#Algorithm, Num Threads, avg. Duration in ms, avg. size in kb, throughput, transferrate in Mbit/s, max Level'

for level in `seq 10 10 120`;
do
	for threads in `seq 1 1 16`;
	do
		pkill java
		sed "s/##NUM_THREADS##/$threads/g" < tourenplaner-template.conf > tourenplaner.conf
		java  -Xmx8g -Xincgc -jar target/tourenplaner-server-1.0-SNAPSHOT-jar-with-dependencies.jar -f dump -c tourenplaner.conf &>out.log &
		sleep 30	
		perftester -concurrent 16 -format csv -algorithm updowng -requests 10 -constrained maxSearchLevel -intConstrained $level -server 'http://localhost:8080' &> /dev/null # warm up
		sleep 1
		echo -e -n "updowng, $threads, " `perftester -concurrent 16 -format csv -algorithm updowng -constrained maxSearchLevel -intConstrained $level -requests 200 -server 'http://localhost:8080'` ", $level\n"
	done
	# For gnuplot datablocks
	echo -e -n '\n\n'
done
