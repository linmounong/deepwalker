# deepwalker

Walker for deepwalk, generating sentences from a graph for word2vec training.

# install

```
go get github.com/linmounong/deepwalker

$ deepwalker -h
Usage of ./deepwalker:
  -alpha float
        probability of restarts
  -input string
        input graph file in adjlist format (default "-")
  -number-walks int
        number of random walks to start at each node (default 10)
  -output string
        output representation file (default "-")
  -save-vocab string
        file to save vocab, format conforming to the c implementation
  -walk-length int
        length of the random walk started at each node (default 40)
  -workers int
        number of parallel go routines (default 1)
```
