#!/bin/bash
cd /home/tbeppe/Documents/tcc/corretoras/Client/Client_200/
for (( c=1; c<=$1; c++ ))
do  
    firefox index.html &
done
