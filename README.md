# reuse
Simple Go app to test TCP and SSL/TLS session reuse.  This is similar
to [httpstat](https://github.com/reorx/httpstat), except httpstat only
makes one connection.  Reuse makes multiple connections and shows what
parts of the connection process are reused.  Command line options
are inspired by [curl](https://curl.haxx.se/).

# Usage

```
Usage:
  reuse [OPTIONS] URL

Application Options:
  -r, --repetitions= Number of times to repeat connecting (default: 2)
      --max-redirs=  Maximum number of redirects to follow (default: 10)
  -k, --insecure     Skip SSL Verification
  -V, --version      Print Version number and exit
  -d, --data=        Data to send
  -X, --request=     HTTP method (default: GET)
  -w, --wait=        Time to wait between connections (default: 5s)
  -H, --header=      Additional Header

Help Options:
  -h, --help         Show this help message

Arguments:
  URL:               URL to connect to
```

# Output

```
lport Remote Address        Rsp DTS DNS   TCP   SSL   TConn Srv   Reply Total

...

--- statistics ---
         Min    Max    Mean   StdDev  Median  MAD
DNS
TCP
SSL
TConn
Srv
Reply
Total

```

## lport

This is local port number of the outgoing connection.

## Remote Address

This is the IP address and port number of the machine reuse is
connecting to

##  Rsp

This is the HTTP Response code from the connection

## DTS

This is actually 3 items.  D represents if DNS coalescing was used
during the call.  1 for True, 0 for False.  T represents if an
existing TCP connection was reused.  1 for True, 0 for False.  S
represents if SSL parameters were reused from an existing session.  1
for True and 0 for False.

It is worth noting that if the TCP session is re-used, then the
application will **NOT** report that the SSL session was re-used.
This is because it records if during SSL negotiation previously agreed
upon parameters are reused.  If the TCP session is reused there is no
SSL negotiation, the existing parameters are used.

## DNS

This is the time in milliseconds between Starting a DNS query and
receiving an answer.

## TCP

This is the time in milliseconds between when a TCP connection is
established and when the DNS query is returned.

## SSL

This is the time in milliseconds between the time a TCP connection is
established and when the request is written.

## TConn

This is the total time in milliseconds between when Go requests a
connection and when one is returned.  This includes the DNS, TCP and
SSL times.

## Srv
This is time in milliseconds between when the request is written and
the first byte of the response is received.

## Reply
This is the time in milliseconds between when the entire response is
received and when the first byte is received.

## Total

This is the time in milliseconds to request the connection and receive
a response.

## Summary Statistics
### Min

This is the smallest value seen during the run.

### Max

This is the largest value seen during the run.

### Mean

This is the sum of values collected during the run divided by the
number of values.  This can be skewed by outliers.

### StdDev

This is the Standard Deviation of the data set.  This gives a measure
of statistical dispersion and is **not** resilient to outliers in the
data.

### Median

This is the value where half the values are above and half are above.
This value is is not skewed as much by extremely large or small
values.

### MAD

This is the Median Absolute Deviation.  This gives a measure of
statistical dispersion, and is resilient to outliers in the data.



# Example Missing TCP keepalive window
```
reuse -k -r 10 -w 6s 'https://api-example.target.com/v1?q=123456&key=APIKEY'
lport Remote Address        Rsp DTS DNS   TCP   SSL   TConn Srv   Reply Total
52600 10.97.109.188:443     200 000 36.6  11.3  47.3  95.1  17.3  0.4   113.1
52601 10.97.109.188:443     200 001 2.5   12.8  21.9  37.2  19.6  0.2   57.1 
52602 10.97.109.188:443     200 000 392.6 17.3  43.1  453.0 20.9  0.2   474.2
52604 10.97.109.188:443     200 001 2.5   11.7  20.0  34.2  17.3  0.3   51.9 
52605 10.97.109.188:443     200 001 2.4   15.5  14.1  32.0  15.7  0.2   48.0 
52606 10.97.109.188:443     200 000 2.6   10.6  32.0  45.1  24.3  0.2   69.8 
52607 10.97.109.188:443     200 001 412.2 12.7  11.6  436.4 16.4  0.3   453.2
52609 10.97.109.188:443     200 000 2.6   9.1   47.9  59.5  25.1  0.2   85.0 
52611 10.97.109.188:443     200 001 2.5   9.6   10.2  22.4  27.6  0.2   50.3 
52612 10.97.109.188:443     200 000 2.4   15.6  33.7  51.6  22.3  0.2   74.3 
--- statistics ---
         Min    Max    Mean   StdDev  Median  MAD
DNS      2.4    412.2  85.9   158.6   2.5     0.1
TCP      9.1    17.3   12.6   2.6     12.2    2.1
SSL      10.2   47.9   28.2   13.9    27.0    14.1
TConn    22.4   453.0  126.7  160.2   48.4    15.3
Srv      15.7   27.6   20.7   3.9     20.3    3.4
Reply    0.2    0.4    0.3    0.1     0.2     0.0
Total    48.0   474.2  147.7  159.2   72.1    21.0

```

# Example hitting TCP keepalive window

```
reuse -k -r 10 -w 3s 'https://api-example.target.com/v1?q=123456&key=APIKEY'
lport Remote Address        Rsp DTS DNS   TCP   SSL   TConn Srv   Reply Total
52622 10.97.109.188:443     200 000 128.5 15.2  56.9  200.5 18.2  0.3   219.3
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   16.5  0.2   16.8 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   17.7  0.3   18.1 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   26.6  0.2   26.9 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   15.3  0.2   15.6 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   17.0  0.3   17.4 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   14.5  0.2   14.9 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   18.5  0.2   18.9 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   15.9  0.2   16.2 
52622 10.97.109.188:443     200 010 0.0   0.0   0.0   0.0   17.3  0.2   17.7 
--- statistics ---
         Min    Max    Mean   StdDev  Median  MAD
DNS      0.0    128.5  12.9   38.6    0.0     0.0
TCP      0.0    15.2   1.5    4.6     0.0     0.0
SSL      0.0    56.9   5.7    17.1    0.0     0.0
TConn    0.0    200.5  20.1   60.1    0.0     0.0
Srv      14.5   26.6   17.8   3.2     17.1    1.2
Reply    0.2    0.3    0.3    0.0     0.2     0.0
Total    14.9   219.3  38.2   60.5    17.5    1.3

```
