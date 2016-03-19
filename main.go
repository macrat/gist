package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/macrat/gist/gist"
)

func files2str(x map[string]gist.FileInfo) string {
	var keys []string
	for k, _ := range x {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}

func PrintList(verbose bool, num int, listFunc func() ([]gist.Overview, error)) {
	if parsed, err := listFunc(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	} else {
		for i, x := range parsed {
			if i > num && num > 0 {
				break
			}
			fmt.Printf("%d %s/%s %s\n", i, x.Owner.Login, files2str(x.Files), x.UpdatedAt)
			fmt.Printf(" %s files: %d comments: %d\n", x.ID, len(x.Files), x.Comments)
			if verbose {
				fmt.Printf(" created at: %s updated at: %s\n", x.CreatedAt, x.UpdatedAt)
				fmt.Println("", x.URL)
			}
			fmt.Println("", strings.Replace(x.Description, "\n", "\n ", -1))
			if len(x.Description) > 0 {
				fmt.Println()
			}
		}
	}
}

func parseID(identifier string) (id, fname string, err error) {
	id = identifier
	if xs := strings.SplitN(id, "/", 2); len(xs) == 2 {
		fname = xs[1]
		id = xs[0]
	}

	if idx, err := strconv.Atoi(id); err == nil && idx >= 0 {
		if gists, err := gist.GetList(); err != nil {
			return "", "", err
		} else if idx < len(gists) {
			id = gists[idx].ID
		}
	}

	return
}

func PrintGist(verbose bool, num int, id string) {
	id, fname, err := parseID(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	gist, err := gist.GetGist(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}

	if verbose {
		fmt.Printf("%s\n", gist.ID)
		fmt.Printf("%s\n", gist.URL)
		fmt.Printf("created at: %s updated at: %s\n", gist.CreatedAt, gist.UpdatedAt)
		fmt.Printf("files: %d comments: %d forks: %d\n", len(gist.Files), gist.Comments, len(gist.Forks))
		fmt.Println(gist.Description)
		fmt.Println()
	}

	count := 0
	for name, file := range gist.Files {
		if fname != "" && name != fname {
			continue
		}
		if num > 0 && count >= num {
			break
		}
		count++

		if verbose {
			fmt.Printf("==> %s / %s(%s) <==\n", name, file.Language, file.Type)
			fmt.Println(file.Content)
		} else if fname == "" && num != 1 && len(gist.Files) > 1 {
			fmt.Println("==>", name, "<==")
			fmt.Println(file.Content)
		} else {
			fmt.Print(file.Content)
		}
	}
}

func CreateGist(filename, description string) {
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	result, err := gist.CreateGist(filename, description, string(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result.URL)
}

func main() {
	flag.Usage = func() {
		fmt.Printf("$ %s [options...] [gist ID[/file name]]\n", os.Args[0])
		fmt.Println()
		fmt.Println("optional:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("positional:")
		fmt.Println("  gist ID     If given ID or index, show content.")
		fmt.Println("  file name   Show only that file even if gist has multiple files.")
	}

	num := flag.Int("n", 0, "Number of items. If 0 or less, be show all items.")
	starred := flag.Bool("s", false, "Show starred gists.")
	create := flag.String("c", "", "New gist file name. If given it, create new gist from stdin.")
	description := flag.String("d", "", "Description for gist.")
	verbose := flag.Bool("v", false, "Enable verbose output.")
	help := flag.Bool("h", false, "Show this message.")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if args := flag.Args(); len(args) >= 2 {
		fmt.Fprintln(os.Stderr, "Can't give multiple IDs.")
		fmt.Fprintln(os.Stderr, "Please see help.")
		os.Exit(1)
	} else if len(args) == 1 {
		for x, y := range map[bool]string{
			*starred:              "-s",
			len(*create) > 0:      "-c",
			len(*description) > 0: "-d",
		} {
			if x {
				fmt.Fprintln(os.Stderr, "Can't use", y, "option with show gist.")
				fmt.Fprintln(os.Stderr, "Please see help.")
				os.Exit(1)
			}
		}
		PrintGist(*verbose, *num, args[0])
		os.Exit(0)
	}

	if len(*create) > 0 {
		for x, y := range map[bool]string{
			*starred:  "-s",
			*num != 0: "-n",
			*verbose:  "-v",
		} {
			if x {
				fmt.Fprintln(os.Stderr, "Can't use", y, "option with create gist.")
				fmt.Fprintln(os.Stderr, "Please see help.")
				os.Exit(1)
			}
		}
		CreateGist(*create, *description)
		os.Exit(0)
	}

	if *starred {
		PrintList(*verbose, *num, gist.GetStarredList)
	} else {
		PrintList(*verbose, *num, gist.GetList)
	}
}
