# korgi - k8s organizer

<p align="center">
  <img src="https://emojis.slackmojis.com/emojis/images/1488330086/1793/party-corgi.gif?1488330086">
   </br>
   Fetches your templated manifests and delivers it to Kapp like a good boi
</p>

---

### WARNING: Very early release, use at your own discretion.

Tool to chain templating engines for k8s and execution engines for k8s. Depends on a very opinionated deployment structure:
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

A reference implementation of the assumed deployment structure can be found [here](https://github.com/DataReply/korganizor-reference).

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
korgi apply -e dev -n default -a dummy monitoring
```

Passing extra args to the engines:

```
korgi --helmfile-args "--skip-deps" --kapp-args "--color=false" apply-namespace default
```

Delete a single group:

```
korgi delete -e dev -n default monitoring
```
