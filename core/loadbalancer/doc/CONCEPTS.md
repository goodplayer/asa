####model:
1. proxy
2. client side
3. look-aside

####pick:
##### 1. round-robin
* rr
* weight rr
##### 2. random
* random
* weight random
##### 3. least
* connection
* latency
* traffic
##### 4. fixed connection
* hash
* consistency hash
* same connection


| - | non-weight | weight | sort |
| ------ | ------ | ------ | ------ |
| round-robin | √ | √ | ✕ |
| random | √ | √ | ✕ |
| fixed key | √ | √ | √ |
| first | ✕ | ✕ | √ |

