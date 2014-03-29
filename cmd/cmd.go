package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type CmdFunc func([]string) error

type Node struct {
	Name        string
	Description string
	Usage       string
	Fn          CmdFunc
	Children    []*Node
}

var (
	ErrInvalidArgs = errors.New("invalid args")
)

func Root(name string) *Node {
	return &Node{
		Name: name,
	}
}

func (n *Node) Parent(name, description string) *Node {
	child := &Node{
		Name:        name,
		Description: description,
	}
	n.Children = append(n.Children, child)
	return child
}

func (n *Node) Command(name, description, usage string, fn CmdFunc) *Node {
	child := &Node{
		Name:        name,
		Description: description,
		Usage:       usage,
		Fn:          fn,
	}
	n.Children = append(n.Children, child)
	return child
}

func (n *Node) Dispatch(args []string, index int) error {
	if n.Fn == nil {
		// Parent.
		if len(args[index:]) == 0 {
			printParentHelp(strings.Join(args[:index], " "), n)
			return nil
		}
		name := args[index]
		for _, c := range n.Children {
			if c.Name == name {
				return c.Dispatch(args, index+1)
			}
		}
		fmt.Printf("invalid command: %s\n\n", name)
		printParentHelp(strings.Join(args[:index], " "), n)
		return nil
	} else {
		// Command.
		err := n.Fn(args[index:])
		if err == ErrInvalidArgs {
			printCommandHelp(strings.Join(args[:index], " "), n)
			return nil
		}
		return err
	}
}

func printParentHelp(matchedCmd string, n *Node) {
	fmt.Printf("usage: %s <command> [<args>]\n\n", matchedCmd)
	fmt.Println("commands:\n")
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, child := range n.Children {
		fmt.Fprintln(w, " ", child.Name, "\t", child.Description)
	}
	w.Flush()
	fmt.Println()
}

func printCommandHelp(matchedCmd string, n *Node) {
	fmt.Printf("usage: %s %s\n", matchedCmd, n.Usage)
}
