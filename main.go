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
			if i >= num && num > 0 {
				break
			}
			fmt.Printf("%d %s/%s %s\n", i, x.Owner.Login, files2str(x.Files), x.UpdatedAt)
			fmt.Printf(" %s files: %d comments: %d\n", x.ID, len(x.Files), x.Comments)
			if verbose {
				fmt.Printf(" created at: %s updated at: %s\n", x.CreatedAt, x.UpdatedAt)
				fmt.Println("", x.HTMLURL)
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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	if verbose {
		fmt.Printf("%s\n", gist.HTMLURL)
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
		os.Exit(-1)
	}

	fmt.Println(result.HTMLURL)
}

func EditGist(id, description string) {
	id, fname, err := parseID(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	if fname == "" || description == "" {
		gist, err := gist.GetGist(id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(-1)
		}

		if fname == "" {
			if len(gist.Files) != 1 {
				fmt.Fprintln(os.Stderr, "This gist has multiple files.")
				os.Exit(1)
			}

			for k, _ := range gist.Files {
				fname = k
			}
		}

		if description == "" {
			description = gist.Description
		}
	}

	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result, err := gist.UpdateGist(id, fname, description, string(content))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Println(result.HTMLURL)
}

func DeleteGist(id string) {
	id, fname, err := parseID(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	if fname != "" {
		fmt.Fprintln(os.Stderr, "Can't specify file name when deleting gist.")
		os.Exit(1)
	}

	if err := gist.DeleteGist(id); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
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
	update := flag.Bool("u", false, "Update gist by input from stdin. You must specify ID.")
	delete_ := flag.Bool("delete", false, "Delete gist. You must specify ID.")
	description := flag.String("d", "", "Description for gist.")
	verbose := flag.Bool("v", false, "Enable verbose output.")
	help := flag.Bool("h", false, "Show this message.")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	check := func(operation string, flags map[bool]string) {
		for x, y := range flags {
			if x {
				fmt.Fprintf(os.Stderr, "Can't use %s option with %s.", y, operation)
				fmt.Fprintln(os.Stderr, "Please see help.")
				os.Exit(1)
			}
		}
	}

	if args := flag.Args(); len(args) >= 2 {
		fmt.Fprintln(os.Stderr, "Can't give multiple IDs.")
		fmt.Fprintln(os.Stderr, "Please see help.")
		os.Exit(1)
	} else if len(args) == 1 {
		if *update {
			check("update gist", map[bool]string{
				*delete_:         "-delete",
				*num != 0:        "-n",
				*starred:         "-s",
				*verbose:         "-v",
				len(*create) > 0: "-c",
			})
			EditGist(args[0], *description)
		} else if *delete_ {
			check("delete gist", map[bool]string{
				*num != 0:        "-n",
				*starred:         "-s",
				*verbose:         "-v",
				len(*create) > 0: "-c",
			})
			DeleteGist(args[0])
		} else {
			check("show gist", map[bool]string{
				*starred:              "-s",
				len(*create) > 0:      "-c",
				len(*description) > 0: "-d",
			})
			PrintGist(*verbose, *num, args[0])
		}
		os.Exit(0)
	}

	if len(*create) > 0 {
		check("create gist", map[bool]string{
			*delete_:  "-delete",
			*num != 0: "-n",
			*starred:  "-s",
			*update:   "-u",
			*verbose:  "-v",
		})
		CreateGist(*create, *description)
		os.Exit(0)
	}

	if *starred {
		PrintList(*verbose, *num, gist.GetStarredList)
	} else {
		PrintList(*verbose, *num, gist.GetList)
	}
}
