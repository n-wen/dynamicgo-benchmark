#!/bin/bash
for ((i=1;i<256;i++))
do
    sudo ifconfig lo0 alias 127.1.1.$i up
done
