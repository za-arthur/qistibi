## Deploy `qistibi`

- Copy files:

```
cd qistibi/ansible
cp ../qistibi roles/qistibi/files/
cp -r ../static roles/qistibi/files
cp -r ../templates roles/qistibi/files
```

- Run Ansible:

```
ansible-playbook -K -i hosts qistibi.yml
```
