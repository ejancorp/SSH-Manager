# SSH Manager
Having a lot of servers can be gruesome to handle, going to ssh on every machine that requires password is hassle.
You can just save it to a file and pick from the list.

![image](image.png)

```
$ git clone git@github.com:ejancorp/ssh-manager.git
$ cd ssh-manager
$ sudo chmod +x ./ssh-manager
$ ssh-manager -f servers.yaml
```