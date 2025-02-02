set -o errexit

usage="$(basename "$0 {NAME} {NAMESPACE}") [-h]

This script installs demonstration scernario for Kyma's GitHub Connector.

where:
	NAME 		- release name, under which scenario will be deployed
	NAMESPACE 	- namespace in which scenario will be deployed

options:
	-h  		shows this help prompt"

if getopts 'h' option; then
	case $option in
		h) echo "$usage"
	   		exit 0
	  ;;
   	*) echo "Please try '$0 -h' for help." >&2
	   		exit 1
	  ;;
  	esac
fi

case $# in
	0) echo "ERROR: No arguments supplied. Please specify {NAME} and {NAMESPACE}."
		  echo "$usage"
		  exit 1
	;;
	1) echo "ERROR: Please specify {NAMESPACE}."
		  echo "$usage"
		  exit 1
	;;
	*)
esac

NAME=$1
NAMESPACE=$2

if [ -z `kubectl get application | grep "${NAME}-app"` ];then
  echo "ERROR: Application '${NAME}-app' does not exist. Check specified {NAME}: '${NAME}' or check if GitHub Connector has been correctly deployed."
  exit 1
fi

cat <<EOF_ | kubectl create -f -
apiVersion: applicationconnector.kyma-project.io/v1alpha1
kind: ApplicationMapping
metadata:
  name: ${NAME}-app
  namespace: ${NAMESPACE}
  labels:
    app: ${NAME}-app
    chart: github-connector
    release: ${NAME}
EOF_

sleep 1

EXTERNALNAME=`kubectl get serviceclasses -n ${NAMESPACE} -o jsonpath="{.items[0].spec.externalName}"`

cat <<EOF_ | kubectl create -f -
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  name: ${NAME}
  namespace:  ${NAMESPACE}
spec:
  serviceClassExternalName: ${EXTERNALNAME}
EOF_

cat <<EOF | kubectl apply -f -
apiVersion: kubeless.io/v1beta1
kind: Function
metadata:
  name: ${NAME}-lambda
  namespace: ${NAMESPACE}
  labels:
    app: ${NAME}
spec:
  deployment:
    spec:
      template:
        spec:
          containers:
          - name: ""
            resources: {}
  deps: |-
    {
        "name": "example-1",
        "version": "0.0.1",
        "dependencies": {
          "request": "^2.85.0"
        }
    }
  function: |-
    module.exports = { main: function (event, context) {
        console.log("Issue opened")
    } }
  function-content-type: text
  handler: handler.main
  horizontalPodAutoscaler:
    spec:
      maxReplicas: 0
  runtime: nodejs8
  service:
    ports:
    - name: http-function-port
      port: 8080
      protocol: TCP
      targetPort: 8080
    selector:
      created-by: kubeless
      function: ${NAME}-lambda
  timeout: ""
  topic: issuesevent.opened
EOF

cat <<EOF | kubectl apply -f -
apiVersion: eventing.kyma-project.io/v1alpha1
kind: Subscription
metadata:
  labels:
    Function: ${NAME}-lambda
  name: ${NAME}-lambda-issuesevent-opened-v1sub
  namespace: ${NAMESPACE}
spec:
  endpoint: http://${NAME}-lambda.${NAMESPACE}:8080/
  event_type: issuesevent.opened
  event_type_version: v1
  include_subscription_name_header: true
  source_id: ${NAME}-app
EOF

echo "Happy GitHub Connecting!"
exit 0
