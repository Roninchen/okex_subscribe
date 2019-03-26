## contribution

### Begin

* fork my repository Roninchen/okex_subscribe.git
* `git clone https://github.com/Roninchen/okex_subscribe.git`


```
make addupstream
```

#### develop begin： this branch needs to be set up by yourself

```
make branch b=mydevbranchname
```

#### develop done: push 

```
make push b=mydevbranchname m="one message to commit"
```


If m is not set, the GIT commit command will not be executed

### modify others pull requset

like i want to modify name=chauncy branch chauncy-fix 的pr

#### step1: pull

```
make pull name=chauncy b=chauncy-fix
```

then modify, done ,commit

#### step2: push modify done 

```
make pullpush name=chauncy b=chauncy-fix
```

