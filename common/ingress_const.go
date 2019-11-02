//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package common

//IngressAnnotationClass ingress
const IngressAnnotationClass = "kubernetes.io/ingress.class"

//IngressAnnotationClassNginx nginx
const IngressAnnotationClassNginx = "nginx"

//IngressAnnotationRewriteTarget nginx
const IngressAnnotationRewriteTarget = "nginx.ingress.kubernetes.io/rewrite-target"

//IngressAnnotationRewriteTargetVal nginx
const IngressAnnotationRewriteTargetVal = "/"

//IngressAnnotationStaticIP static ip
const IngressAnnotationStaticIP = "kubernetes.io/ingress.global-static-ip-name"

//IngressRoutePrefixAPIGateway ingress apigatewayroute prefix
const IngressRoutePrefixAPIGateway = "/api"

//IngressRoutePrefixManagerAPI ingress manager api prefix
const IngressRoutePrefixManagerAPI = "/manager/api"

const ingressHostRoot = "quebic.io"

//IngressHostAPIGateway ingress api-gateway host
const IngressHostAPIGateway = "api" + "." + ingressHostRoot

//IngressHostManager ingress manager host
const IngressHostManager = "api" + "." + "mgr" + "." + ingressHostRoot
