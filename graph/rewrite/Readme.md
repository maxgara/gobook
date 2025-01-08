## graph
graph is a CLI tool to create HTML documents containing plots of data.

data format is inferred unless specified by the user.


### input
graph reads data from files or stdin, and can be invoked with flags to control rendering.
supported flags initially will be limited to:  

`--multiseries[=0|1]`  
`-s (multiseries=1)  `  
> when set to 1, this flag indicates that successive data series should render as separate plots on the same graph.
the flag is included to allow almost-backwards-compatibility with previous versions, tracked in graph/minimal,
in which this was the default operating scheme.



`grid=$count`
>- this flag indicates that successive graphs render in a grid, ie. a 3 by 10 grid can be created with —gridrows=3, when supplying data for 30 graph elements.

data is read line by line, with additional flag options allowed inside of the data file:


### Input Formatting Flags:  
`-d1`  
>each column of input data in the file is a separate data series; each line of data contains y values for each series corresponding to the same x (line index)
  
`-d2`   
>each consecutive pair of columns is a data series; the first column is an x-coord and the second is a y-coord.

`-d3`  
>similar to d1 and d2; for three dimensional data plots

`-dx`  
>each data series shares the x value at the beginning of each data line, unless overridden by dt.

`-dt`  
>time series. data plot is animated, lines of data begin with 1 additional coordinate for time. similar to dx option, each series shares this value. 
>if dx is also specified, then t is the first line entry and x is the second. data series grouping is still controlled by options d{1,2,3} after the initial time index on each line.

`--stack=true`  
>each series is stacked on the previous ones. ie. points [(3,3), (3,1), (3,2)] are displayed like [(3,3), (3,4), (3,6)]. Areas between curves are shaded.  





