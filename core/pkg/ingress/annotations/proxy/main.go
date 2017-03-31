/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"k8s.io/kubernetes/pkg/apis/extensions"

	"k8s.io/ingress/core/pkg/ingress/annotations/parser"
	"k8s.io/ingress/core/pkg/ingress/resolver"
)

const (
	bodySize     = "ingress.kubernetes.io/proxy-body-size"
	connect      = "ingress.kubernetes.io/proxy-connect-timeout"
	send         = "ingress.kubernetes.io/proxy-send-timeout"
	read         = "ingress.kubernetes.io/proxy-read-timeout"
	bufferSize   = "ingress.kubernetes.io/proxy-buffer-size"
	cookiePath   = "ingress.kubernetes.io/proxy-cookie-path"
	cookieDomain = "ingress.kubernetes.io/proxy-cookie-domain"
)

// Configuration returns the proxy timeout to use in the upstream server/s
type Configuration struct {
	BodySize       string `json:"bodySize"`
	ConnectTimeout int    `json:"conectTimeout"`
	SendTimeout    int    `json:"sendTimeout"`
	ReadTimeout    int    `json:"readTimeout"`
	BufferSize     string `json:"bufferSize"`
	CookieDomain   string `json:"cookieDomain"`
	CookiePath     string `json:"cookiePath"`
}

type proxy struct {
	backendResolver resolver.DefaultBackend
}

// NewParser creates a new reverse proxy configuration annotation parser
func NewParser(br resolver.DefaultBackend) parser.IngressAnnotation {
	return proxy{br}
}

// ParseAnnotations parses the annotations contained in the ingress
// rule used to configure upstream check parameters
func (a proxy) Parse(ing *extensions.Ingress) (interface{}, error) {
	defBackend := a.backendResolver.GetDefaultBackend()
	ct, err := parser.GetIntAnnotation(connect, ing)
	if err != nil {
		ct = defBackend.ProxyConnectTimeout
	}

	st, err := parser.GetIntAnnotation(send, ing)
	if err != nil {
		st = defBackend.ProxySendTimeout
	}

	rt, err := parser.GetIntAnnotation(read, ing)
	if err != nil {
		rt = defBackend.ProxyReadTimeout
	}

	bufs, err := parser.GetStringAnnotation(bufferSize, ing)
	if err != nil || bufs == "" {
		bufs = defBackend.ProxyBufferSize
	}

	cp, err := parser.GetStringAnnotation(cookiePath, ing)
	if err != nil || cp == "" {
		cp = defBackend.ProxyCookiePath
	}

	cd, err := parser.GetStringAnnotation(cookieDomain, ing)
	if err != nil || cd == "" {
		cd = defBackend.ProxyCookieDomain
	}

	bs, err := parser.GetStringAnnotation(bodySize, ing)
	if err != nil || bs == "" {
		bs = defBackend.ProxyBodySize
	}

	return &Configuration{bs, ct, st, rt, bufs, cd, cp}, nil
}
