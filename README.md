# cmr
cmr command tool

https://www.digitalocean.com/community/tutorials/how-to-use-a-private-go-module-in-your-own-project
https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account


once the machines ssh key has been added to the account
add this to ~/.gitconfig 
```
[url "ssh://git@github.com/"]
	insteadOf = https://github.com/
```

then to get a repo using ssh
```
git clone git@github.com:stalwartgiraffe/cmr.git
```
