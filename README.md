## Introduction

I am writing this technical document to outline an issue that I am observing on a <strike>daily</strike> hourly basis with my internet service provider.

### What is Packet Loss?

Packet loss is the loss of data packets during transmission over a computer network. It is a common problem in computer networks, and can be caused by a variety of factors. Packet loss can be caused by a number of factors, including network congestion, hardware failure, and software bugs. Packet loss can also be caused by a number of factors, including network congestion, hardware failure, and software bugs.

### What is latency?

Latency is the time it takes for a packet of data to go from its origin to its destination.  Latency is typically measured in milliseconds.  In competitive gaming, a time greater than 50 milliseconds is considered to be above average.

### Why is Packet Loss a problem for me?

Simple, my family are gamers.  We play highly competitive online games where stability and low latency are key to winning.

### My ISP Plan

I pay for the following package:

    200 Mbps down, 100 Mbps up

I get the following speeds (on average):

    160 Mbps down, 72 Mbps up

I feel like these are respectable numbers, but I am not getting what I am paying for.

Additionally, I typically get a ping time of 15-30 ms, which is also respectable.

### Setup

After multiple house calls to attempt to resolve the issue, I have decided to take matters into my own hands.  I have began the process to monitor connectivity using the following setup:

1.  Nighthawk R7000 router with DD-WRT firmware
2.  Raspberry Pi 3 running Raspbian (connected via ethernet to router)
3.  Developed an open source Golang based application to monitor connectivity and report to a remote server <a href="https://github.com/bdwilliams/packetloss" target="_blank">https://github.com/bdwilliams/packetloss</a>.
4.  Configured application to run on 15 minute intervals pinging the following hosts:
    1.  My Router
    2.  My ISP's Gateway
    3.  My ISP's Provider
    4.  My ISP's Provider's Provider (Atlas Cogentco Tulsa)
    5.  My ISP's Provider's Provider (Atlas Cogentco OKC)
5.  Write every 15 minute report to a remote server hosting MySQL (currently a Digital Ocean Droplet)
6.  Write a simple web application to display the data in a meaningful way (Coming Soon!)

### Results (October 17th, 2022)

I have been running this test for almost 2 weeks now.  I have seen packet loss as high as 12%.  I would best describe the test results as a yo-yo.  When monitoring, I see it stay around 1-2% and increase up to 10-12% in spurts.

#### Cumulative Packet Loss (all time)

I decided one of the queries I wanted to run was a cumulative count on the the total amount of packet loss grouped by each individual host and what I found was interesting.  The issue is ... not my ISP.

mysql> select sum(packet_loss) As total_packet_loss, name from ping_results where recv > 0 group by name order by sum(packet_loss);  
+-------------------+---------------------------------------------------------+  
| total_packet_loss | name                                                    |  
+-------------------+---------------------------------------------------------+  
|                 0 | Home Router                                             |  
|                 0 | DiamondNet                                              |  
|                88 | DiamondNet Provider (rcr21.tul01.atlas.cogentco.com)    |  
|               264 | DiamondNet Provider (rcr21.okc01.atlas.cogentco.com)    |  
|               380 | DiamondNet Provider (static137-block97.intelleqcom.net) |  
+-------------------+---------------------------------------------------------+  
5 rows in set (0.07 sec)  

As you can see, the total amount of packet loss is not coming from my ISP, but from the providers of my ISP.

#### Cumulative Packet Loss (all time) - averaging packet loss

mysql> select count(*) AS total_reports, sum(packet_loss) As total_packet_loss, sum(packet_loss) / count(*) AS average_packet_loss, name from ping_results where recv > 0 group by name order by sum(packet_loss);  
+---------------+-------------------+---------------------+---------------------------------------------------------+  
| total_reports | total_packet_loss | average_packet_loss | name                                                    |  
+---------------+-------------------+---------------------+---------------------------------------------------------+  
|           332 |                 0 |                   0 | Home Router                                             |  
|           330 |                 0 |                   0 | DiamondNet                                              |  
|           330 |                88 | 0.26666666666666666 | DiamondNet Provider (rcr21.tul01.atlas.cogentco.com)    |  
|           330 |               264 |                 0.8 | DiamondNet Provider (rcr21.okc01.atlas.cogentco.com)    |  
|           330 |               380 |  1.1515151515151516 | DiamondNet Provider (static137-block97.intelleqcom.net) |  
+---------------+-------------------+---------------------+---------------------------------------------------------+  
5 rows in set (0.07 sec)  

