## normal

Components describe whole set. Granurality is on set, not on individual iterations.

### 12 x 100 - 1.R 2.P

1. 100R - orders 0, 2, 4, 6, 8, 10
2. 100P - orders 1, 3, 5, 7, 9, 11

## compound

Components describe one iteration of distance. Granurality is on iteration.

### 50K 50P ... 400

1. 50K - orders 0, 2, 4, 6
2. 50P - orders 1, 3, 5, 7

### 8 x 50 - 25K 25P

1. 25K - orders 0
2. 25P - orders 1

## super

Components describe each unique iteration. Components have iteration order filled,
in order to partition them to their respective iterations/sets.

### 4 x 200 - 1. 50K 50P 2. 50Z 50M

1. 50K - iteration order 0 - orders 0, 2
2. 50P - iteration order 0 - orders 1, 3

3. 50Z - iteration order 1 - orders 0, 2
4. 50M - iteration order 1 - orders 1, 3

## pyramid

Components describe each step in pyramid.

### start: 50 stop: 150 step 50

1. 50 - orders 0, 4
1. 100 - orders 1, 3
1. 150 - orders 2
