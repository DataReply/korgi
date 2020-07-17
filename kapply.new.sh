#!/bin/bash

# NOTES:
# Don't try to pass around complex structs etc, have to modify state (could be global, or embedded in a function)
# Otherwise check for function return signals, but that's about it.
# ISSUES:
# * Cmmplex logic control flows don't work
# * Data structs don't exist
# * Probably best to rewrite this in Python or Ruby before making the control flow logic even more complex.
# GUIDE ON STYLE: https://google.github.io/styleguide/shellguide.html

set -o errexit
set -o pipefail
set -o nounset

ROOT="$PWD"
LOGGER="$ROOT/tools/color_logger.sh"
GENSUB="$ROOT/tools/gensub.py"
source "$LOGGER"
PARENT_DIR="/tmp/kapp/$(date '+%m-%d-%Y_%H-%M-%S')"
PIPELINE=
HELMFILE=
DELETE=
APPLY=
TEMPLATE=
LINT=

# These will be populated from the arguments being passed to the script.
FUNCTION=
NAMESPACE=
ENVIRONMENT=
GROUP=

function die() {
  echo "$*" >&2
  exit 2
} # complain to STDERR and exit with error

function needs_arg() {
  if [[ -z $OPTARG ]]; then
    die "No arg for --$OPT option"
  fi
}

while getopts htladpf:g:e:n:b:-: OPT; do
  # support long options: https://stackoverflow.com/a/28466267/519360
  if [ "$OPT" = "-" ]; then # long option: reformulate OPT and OPTARG
    OPT="${OPTARG%%=*}"     # extract long option name
    OPTARG="${OPTARG#$OPT}" # extract long option argument (may be empty)
    OPTARG="${OPTARG#=}"    # if long option argument, remove assigning `=`
  fi
  case "$OPT" in
    g | group)
      needs_arg
      GROUP="$OPTARG"
      ;;
    e | environment)
      needs_arg
      ENVIRONMENT="$OPTARG"
      ;;
    n | namespace)
      needs_arg
      NAMESPACE="$OPTARG"
      ;;
    f | function)
      needs_arg
      FUNCTION="$OPTARG"
      ;;
    b | bases)
      PARENT_DIR="$OPTARG"
     ;;
    l | lint)  LINT=true ;;
    t | template) TEMPLATE=true ;;
    p | pipeline) PIPELINE=true ;;
    h | helmfile) HELMFILE=true ;;
    d | delete) DELETE=true ;;
    a | apply) APPLY=true ;;
    ??*) die "Illegal option --$OPT" ;; # bad long optionk
    \?) exit 2 ;;                       # bad short option (error reported via getopts)
  esac
done

shift $((OPTIND - 1)) # remove parsed options and args from $@ list

function check_binaries() {
  info "Checking for binaries..."
  local -a arr=("pass" "helm" "kapp" "kubectl")
  for i in ${arr[@]}; do
    success "$i: $($i version)" || error "$i not found."
  done
  success "sops: $(sops --version)" || error "sops not found."
  success "gpg2: $(gpg2 --version)" || error "gpg2 not found."
  success "helm secrets: $(helm secrets)" || error "Helm secrets plugin not found."
  success "Binary check completed!"
}


function dir_exists?() {
  function _cd() {
    cd $1
    info "PWD: $PWD"
  }
  if [[ -d "$1" ]]; then
    _cd "$1"
    return 0
  fi
  error "$1 is not found."
  exit 1
}

function _mutate_state() {
  local namespace="$1"
  local app_group="$2"
  shift 2
  local -a arr=$(ls | grep -v -e '_')
  local app_group_dir="$PARENT_DIR/${namespace}/${app_group}"
  info "App Group Tree: $app_group_dir"
  mkdir -p "$app_group_dir"
  local _app=
  for app in ${arr[@]}; do
    # regular expression captures ./1_grafana.yaml or ./abc-Something-123.yml
    # the app will be named 1_grafana, or abc-Something-123 respectively.
    _app=$(python "$GENSUB" "$app")
    mkdir -p "$app_group_dir/$_app"
    if [[ $LINT ]]; then
      info "Linting app: $_app"
      helmfile --environment "$ENVIRONMENT" --file "$app" --state-values-set app="$_app" lint
    fi
    if [[ $TEMPLATE ]]; then
      if [[ $HELMFILE ]]; then
        info "Templating app using Helmfile: $_app"
        helmfile --environment "$ENVIRONMENT" --file "$app" --state-values-set app="$_app" template --output-dir "$app_group_dir/$_app"
        #info "Validating $_app..."
        #time kubeval --ignore-missing-schemas -i kapp.config.yaml -d "$app_group_dir/$_app" --exit-on-error -v 1.12.10
        #info "Validating $_app done!"
        cp $ROOT/tools/config/kapp.config.yaml $app_group_dir/$_app/kapp.config.yaml

      else
        info "Templating app using Kontemplate: $_app"
        kontemplate template "$app" -o $app_group_dir/$_app
      fi
    fi
  done

  success "Templating complete for all components! Manifests are generated at: $app_group_dir"
  info "Starting the deployment process via kapp!"
  if [[ $DELETE ]]; then
    kapp app-group delete \
      -n $namespace -g $app_group \
      "$@"
    info "Deletion done!"
  elif [[ $APPLY ]]; then
    kapp app-group deploy \
      -d $app_group_dir -n $namespace -g $app_group \
      "$@"
    info "Deployment done!"
  else
    warn "Deployment/deletion skipped for $app_group!"
  fi
}

# helmfile --environment client-k8s.dlp-ingest-cons.aws.de.pri.o2.com  --file grafana.yaml template

function run() {
  mkdir -p $PARENT_DIR
  info "All templates are here: $PARENT_DIR"
  echo "ARGS PASSED TO RUN: $@"
  local namespace="$1"
  local group="$2"
  local environment="$3"
  shift 3
  local args=()
  set +o nounset
  if [[ $PIPELINE ]]; then
    args+=("--kubeconfig-context")
    args+=("${environment}/${namespace}")
  fi
  args+=("$@")
  info "Additional args: ${args[@]}"
  _mutate_state $namespace $group ${args[@]}
  set -o nounset
}

# CAN BE PASSED INTO THE FUNCTION FLAG
function iter_on_app_group() {
  if [[ $HELMFILE ]]; then
    dir_exists? realm
    dir_exists? namespaces
    dir_exists? $NAMESPACE
    dir_exists? $GROUP
  else
    dir_exists? environments
    dir_exists? $ENVIRONMENT
    dir_exists? $NAMESPACE
    dir_exists? $GROUP
  fi
  run $NAMESPACE $GROUP $ENVIRONMENT "$@"
}

function bootstrap_pipeline() {
  local ENDPOINT="https://api."$ENVIRONMENT
  local TOKEN=$( cat /etc/gitlab-cd-runner-keys/admin-user-$NAMESPACE-token )
  kubectl config set-cluster $ENVIRONMENT --server=$ENDPOINT --certificate-authority=/etc/gitlab-cd-runner-keys/admin-user-$NAMESPACE-cert
  kubectl config set-credentials admin-user-$NAMESPACE --token=$TOKEN --certificate-authority=/etc/gitlab-cd-runner-keys/admin-user-$NAMESPACE-cert
  kubectl config set-context $ENVIRONMENT/$NAMESPACE --cluster=$ENVIRONMENT --user=admin-user-$NAMESPACE --namespace=$NAMESPACE
  kubectl config use-context $ENVIRONMENT/$NAMESPACE
  #kubectl config use-context $ENVIRONMENT
  warn "################## KUBECTL CONFIG VIEW ################"
  kubectl config view
  info "Initializing pass..."
  gpg2 --import /etc/gitlab-cd-runner-keys/pubkey.asc /etc/gitlab-cd-runner-keys/privkey.asc
  local FINGERPRINT=$( cat /etc/gitlab-cd-runner-keys/fingerprint )
  pass init $FINGERPRINT
  success "Pass initialized!"
}

function iter_on_ns() {
  if [[ $HELMFILE ]]; then
    dir_exists? realm
    dir_exists? namespaces
    dir_exists? $NAMESPACE
  else
    dir_exists? environments
    dir_exists? $ENVIRONMENT
    dir_exists? $NAMESPACE
  fi
  local -a groups=$(ls | grep -v -e '_' -e 'archive')
  info "App Groups are: ${groups[@]}"
  for group in ${groups[@]}; do
    dir_exists? $group
    run $NAMESPACE $group $ENVIRONMENT "$@"
    cd ..
  done
}

## MAIN
function main() {
  check_binaries
  if [[ $PIPELINE ]]; then
    bootstrap_pipeline
  fi
  if [[ $FUNCTION == "iter-on-namespace" ]]; then
    iter_on_ns "$@"
  elif [[ $FUNCTION == "iter-on-app-group" ]]; then
    if [[ -z $GROUP ]]; then
      error "\$App Group is not given, please provide an app group name."
    else
      iter_on_app_group "$@"
    fi
  else
    error "Wrong function: <$FUNCTION>, please check again."
    error "Possible functions: iter-on-namespace, iter-on-app-group"
  fi
}
main "$@"
