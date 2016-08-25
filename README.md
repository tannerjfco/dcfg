# drudconfig
A cli for configuration management...better than Brad's

## Use

Get the binary. Either build the source or head to the releases section.

Create a Config Set (drud.yaml).

```
- name: install
  tasks:
  - cmd: echo "testing this thing"
    wait: 1s
    repeat: 3
  - name: create directory
    cmd: mkdir testing
    ignore: yes
  - name: make file in new directory
    workdir: testing
    cmd: touch newfile.txt
  - name: look
    cmd: ls -l testing
  - name: echo "woot"

- name: uninstall
  tasks:
  - cmd: echo "testing this thing"
    name: say something cool
  - cmd: rm newfile.txt
    name: remove new file
    workdir: testing
  - name: rm directory
    cmd: rm -rf testing
  - cmd: echo "woot"
```

Now you can run all Config groups in sequential order or you can run jus tone group at a time.

```
dcfg run all
```
or
```
dcfg run install
dcfg run uninstall
```
or
```
dcfg run install uninstall
```