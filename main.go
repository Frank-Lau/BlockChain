package main

func main() {
	bc := NewBlockChain("班长")
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}
