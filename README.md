# korgi - k8s organizer


Very early release, use at your own discretion.

Tool to chain templating engine for k8s and execution engions for k8s. Depends on a very opinionated deployment structure:
```
realm
   namespaces
        namespace
            {app_group}
                _app_group.yaml
                app1.yaml
                app2.yaml
```
Supported templating engines:
- helmfile
- kontemplate

Supported execution engines:
- kapp

# Examples

Apply all groups in namespace `default` and env `dev`:

```
korgi apply-namespace -e dev default
```


Apply a group in namespace `default` and env `dev`:

```
korgi apply -e dev -n default monitoring
```

Apply a single app from the `monitoring` group in namespace `default` and env `dev`:

```
korgi apply -e dev -n default -f dummy monitoring
```

Passing extra args to the engines:

```
korgi --template-engine-args "--skip-deps" --exec-engine-args "--color=false" apply-namespace default
```

Delete a single group:

```
korgi delete -e dev -n default monitoring
```
