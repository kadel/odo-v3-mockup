# odo v3 mockup

This is empty mockup to demo odo v3 CLI.
It actually doesn't perform any action, only when it says "Downloading devfile" it will create empty `devfile.yaml` file.


You can download latest binaries at https://github.com/kadel/odo-v3-mockup/releases/tag/latest

Run it as `odov3`
```
$ odov3 -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application.

Usage:
  odo-v3 [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete
  deploy      Deploy application to Kubernetes cluster.
  dev         Run application in a developer mode.
  help        Help about any command
  init        Bootstrap new application based on Devfile
  list        List existing resources
  login       Log in to your server and save login for subsequent use.
  logout      Log out of the active session out by clearing saved tokens.

Flags:
  -h, --help   help for odo-v3

Use "odo-v3 [command] --help" for more information about a command.
```


## Install
### Quick install on Linux
```
curl -L https://github.com/kadel/odo-v3-mockup/releases/download/latest/linux-amd64-odo3.gz | gzip -d > odov3
chmod +x odov3
sudo mv odov3 /usr/local/bin/odov3
```

### Quick install on Mac
```
curl -L https://github.com/kadel/odo-v3-mockup/releases/download/latest/darwin-amd64-odo3.gz | gzip -d > odov3
chmod +x odov3
sudo mv odov3 /usr/local/bin/odov3
```


### Quick install on Mac M1
```
curl -L https://github.com/kadel/odo-v3-mockup/releases/download/latest/darwin-arm64-odo3.gz | gzip -d > odov3
chmod +x odov3
sudo mv odov3 /usr/local/bin/odov3
```

Quick install on Windows
https://github.com/kadel/odo-v3-mockup/releases/download/latest/windows-amd64-odo3.exe.zip

## Uninstall

### Quick uninstall Linux/Mac/Mac M1
```
sudo rm /usr/local/bin/odov3
```

### Quick uninstall Windows