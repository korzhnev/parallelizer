Introduction
------------

`parallelizer` is a package for the Golang that allows you to parallelize the flow of actions.

We assume that the order of actions is important. But it is important for related actions in fact.
Thus we can split actions into independent groups and process the actions sequentially for each group.

Each action contains an arbitrary number of objects. If different actions contain the same objects, they are considered 
related (transitively!). If actions have do not contain any object, they will be assigned to a separate group.

Groups are processed in separate goroutines. 

For effective separation of actions into groups, optimized Union-Find data structure is used 
under the hood (see package `unionfind`).

## Install
`go get github.com/korzhnev/parallelizer`

## Usage

You should implement the interface parallelizer.Action and run parallelizer.Run().
See `parallelizer_test.go` file for more details. 
  
  


