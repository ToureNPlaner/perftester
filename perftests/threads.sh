#!/bin/zsh
echo '#Algorithm, Num Threads, avg. Duration in ms, avg. size in kb, throughput, transferrate in Mbit/s'
for threads in `seq 1 16`;
do
	pkill java
	sed "s/##NUM_THREADS##/$threads/g" < tourenplaner-template.conf > tourenplaner.conf
	java  -Xmx8g -Xincgc -jar target/tourenplaner-server-1.0-SNAPSHOT-jar-with-dependencies.jar -f dump -c tourenplaner.conf &>out.log &
	sleep 20	
	perftester -concurrent 16 -format csv -algorithm updowng -requests 20 -server 'http://localhost:8080' &> /dev/null # warm up
	echo -e -n "updowng , $threads, " `perftester -concurrent 16 -format csv -algorithm updowng -requests 20 -server 'http://localhost:8080'` '\n'
done
# Fot Gnuplot datablock seperation
echo -e -n '\n\n'
for threads in `seq 1 16`;
do
	pkill java
	sed "s/##NUM_THREADS##/$threads/g" < tourenplaner-template.conf > tourenplaner.conf
	java  -Xmx8g -Xincgc -jar target/tourenplaner-server-1.0-SNAPSHOT-jar-with-dependencies.jar -f dump -c tourenplaner.conf &>out.log &
	sleep 20	
	perftester -concurrent 16 -format csv -algorithm updowng -constrained maxSearchLevel -intConstrained 40 -requests 20 -server 'http://localhost:8080' &> /dev/null # warm up
	echo -e -n "updowngi-40 , $threads, " `perftester -concurrent 16 -format csv -algorithm updowng -constrained maxSearchLevel -intConstrained 40 -requests 20 -server 'http://localhost:8080'` '\n'
done
# Fot Gnuplot datablock seperation
echo -e -n '\n\n'
for threads in `seq 1 16`;
do
	pkill java
	sed "s/##NUM_THREADS##/$threads/g" < tourenplaner-template.conf > tourenplaner.conf
	java  -Xmx8g -Xincgc -jar target/tourenplaner-server-1.0-SNAPSHOT-jar-with-dependencies.jar -f dump -c tourenplaner.conf &>out.log &
	sleep 20	

	perftester -concurrent 16 -format csv -algorithm sp -requests 20 -server 'http://localhost:8080' &>/dev/null #warm up
	echo -e -n "sp , $threads, " `perftester -concurrent 16 -format csv -algorithm sp -requests 20 -server 'http://localhost:8080' 2> /dev/null` '\n'
done

