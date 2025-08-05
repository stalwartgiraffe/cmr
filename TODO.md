make a linter that auto applies go imports ie


	if !isFirst{


create branch from jira url

insert comments based on branch_name

autorebase
squash tail commits
pull main
rebase main
v2 pull main proactively or otherwise indicate new merges to main

create an reservation for experiment machines

create an experiment tag

create an experiment deployment

create an experiment rollback

on push, deploy to gauntlet and scan for log messages

worktree update
 cd to work dir
 switch to main
 switch to main
 pull
 switch to work
 rebase

remove branches that are not mine

opx -
   open current page
   create it if it does not exist


make branch aware chore/feat etc comment-er

make ft  maker
ticket
builder
pusher

make ticker pusher that knows the rules, will autmake ft

make a ft inserter
names its well
adds it to the ft file  
makes the branch
makes the commit

slack client

when getting a review, automate getting the repo and the branch.

in vim diff the current file to is pair in other workspaces.

eazy diff alt workspace:
fzf ... e -d  <alt workspace>...file file
maybe  gno / mr aware

add a spell check plugin to slack 


on local docker test with FT in both states on deply

fix the dang thing where gpnew uses the branch name in the MR title rather then the comment
seems to only happen in the features repo

nvim:
  disable spellunker pluging spelling on files above 1mb
https://www.reddit.com/r/neovim/comments/1fcidod/cspell_for_neovim/
https://github.com/davidmh/cspell.nvim
https://github.com/streetsidesoftware/cspell-cli
https://github.com/nvimtools/none-ls.nvim
https://github.com/streetsidesoftware/cspell
https://github.com/streetsidesoftware/cspell/tree/main/packages/cspell
a service that pushes the number of differs behind tip into the branch prompt
needs some polling and caching


add manual TODO to list of review requests found in jira
ie md txt and url 

# remove local branches from reviews -  not my branches
gb | rg -v Karl
git branch -d branch_name


# automate gauntlet full 
gauntlet vm fill up disk with stale docker images and containers over time
see 
df -f
can manually clean with 
docker images prune
docker container prune
then will need to 
git pull
gauntlet reset


figure out ssh keepalive to gauntlet



something like zsh <C-d> list choices of command but with zsh style search and preview and explanations
