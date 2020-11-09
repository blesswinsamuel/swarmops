package main

import (
	"context"
	"docker_swarm_gitops/internal/config"
	"docker_swarm_gitops/internal/docker"
	"docker_swarm_gitops/internal/git"
	"docker_swarm_gitops/internal/server"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("You must pass a subcommand serve or deploy")
		os.Exit(1)
	}
	cmds := []Runner{
		NewServeCommand(),
		NewDeployCommand(),
		NewUndeployCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			err := cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}

	log.Fatalf("Unknown subcommand: %s", subcommand)
}

type ServeCommand struct {
	fs *flag.FlagSet

	repoDir        string
	privateKeyFile string
	gitRepo        string
	gitBranch      string
	stackFile      string
	syncInterval   time.Duration
	port           string
}

func NewServeCommand() *ServeCommand {
	sc := &ServeCommand{
		fs: flag.NewFlagSet("serve", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.repoDir, "repo-dir", "", "Repo clone directory")
	sc.fs.StringVar(&sc.privateKeyFile, "private-key-file", "", "Key persistense directory")
	sc.fs.StringVar(&sc.gitRepo, "git-repo", "", "SSH git clone repo URL")
	sc.fs.StringVar(&sc.gitBranch, "git-branch", "master", "SSH git branch")
	sc.fs.StringVar(&sc.stackFile, "stack-file", "stack.yaml", "Stack file")
	sc.fs.DurationVar(&sc.syncInterval, "sync-interval", 5*time.Minute, "Sync Interval")
	sc.fs.StringVar(&sc.port, "port", "8080", "Server port")

	return sc
}

func (c *ServeCommand) Name() string {
	return c.fs.Name()
}

func (c *ServeCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *ServeCommand) Run() error {
	var ()
	flag.Parse()

	if c.gitRepo == "" {
		return errors.New("--git-repo should not be empty")
	}
	if c.gitBranch == "" {
		return errors.New("--git-branch should not be empty")
	}

	git, err := git.NewGit(c.gitRepo, c.gitBranch, c.repoDir, c.privateKeyFile)
	if err != nil {
		log.Fatalf("NewGit: %v", err)
	}

	server := server.NewServer(git, c.repoDir, c.stackFile)

	log.Infof("sync interval: %v", c.syncInterval)
	quit := make(chan struct{})
	if c.syncInterval > 0 {
		go server.BackgroundSync(c.syncInterval, quit)
	}
	h := server.HttpHandler()
	s := &http.Server{Addr: ":" + c.port, Handler: h}

	go func() {
		log.Infof("Server started at port %s", c.port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Infof("Server ListenAndServe failed: %v", err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	close(quit)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Infof("Server Shutdown failed: %v", err)
	}
	return nil
}

type DeployCommand struct {
	fs *flag.FlagSet

	stackFile string
	baseDir   string
}

func NewDeployCommand() *DeployCommand {
	sc := &DeployCommand{
		fs: flag.NewFlagSet("deploy", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.stackFile, "stack-file", "stack.yaml", "Stack file")
	sc.fs.StringVar(&sc.baseDir, "base-dir", "./", "Base directory")

	return sc
}

func (c *DeployCommand) Name() string {
	return c.fs.Name()
}

func (c *DeployCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *DeployCommand) Run() error {
	cfg, err := config.ParseConfig(c.baseDir, c.stackFile)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %s", err)
	}
	d := docker.NewDockerStackCmd()
	err = d.Deploy(cfg)
	if err != nil {
		return fmt.Errorf("Failed to run deploy: %s", err)
	}
	return nil
}

type UndeployCommand struct {
	fs *flag.FlagSet

	stackFile string
	baseDir   string
}

func NewUndeployCommand() *UndeployCommand {
	sc := &UndeployCommand{
		fs: flag.NewFlagSet("undeploy", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.stackFile, "stack-file", "stack.yaml", "Stack file")
	sc.fs.StringVar(&sc.baseDir, "base-dir", "./", "Base directory")

	return sc
}

func (c *UndeployCommand) Name() string {
	return c.fs.Name()
}

func (c *UndeployCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *UndeployCommand) Run() error {
	cfg, err := config.ParseConfig(c.baseDir, c.stackFile)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %s", err)
	}
	d := docker.NewDockerStackCmd()
	err = d.Remove(cfg)
	if err != nil {
		return fmt.Errorf("Failed to run undeploy: %s", err)
	}
	return nil
}
