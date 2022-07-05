package GoLimit

type GoLimit struct {
	n int
	c chan struct{}
}

func NewGoLimit(n int) *GoLimit {
	return &GoLimit{
		n: n,
		c: make(chan struct{}, n),
	}
}
func (g *GoLimit) Run(f func()) {
	//g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
