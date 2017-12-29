// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"net/http"
	httputil "github.com/goodrain/rainbond/pkg/util/http"
	"github.com/go-chi/chi"
	"github.com/Sirupsen/logrus"

	"github.com/goodrain/rainbond/pkg/node/api/model"
	"strings"

)

var prometheusAPI  = model.PrometheusAPI{}
var toReplace ="#to_replace#"
//GetNode 获取一个节点详情
func GetMem(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	replaceTo:="instance='"+nodeID+"'"

	basic:="100 - ((node_memory_MemFree{job='rbd_node',#to_replace#} %2B node_memory_Cached{job='rbd_node',#to_replace#} %2B node_memory_Buffers{job='rbd_node',#to_replace#})/node_memory_MemTotal) * 100"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.Query(expr)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}
func replaceSelector(s,new string) string {
	return strings.Replace(s,toReplace,new,-1)
}
func GetCpu(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	replaceTo:="instance='"+nodeID+"'"
	basic:="100 - (avg by (instance) (irate(node_cpu{job='rbd_node', mode='idle',#to_replace#}[5m])) * 100)"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.Query(expr)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetDisk(w http.ResponseWriter, r *http.Request) {
	//basic:="100 - node_filesystem_free{job='rbd_node',#to_replace#,fstype!~'rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'} / node_filesystem_size{job='rbd_node',#to_replace#,fstype!~'rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'} * 100"
	basic:="100 - (avg by (instance) (node_filesystem_free{job='rbd_node',#to_replace#,fstype!~'rootfs|nsfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'} / node_filesystem_size{job='rbd_node',#to_replace#,fstype!~'rootfs|nsfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'}) * 100)"
	nodeID := chi.URLParam(r, "node_id")
	replaceTo:="instance='"+nodeID+"'"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.Query(expr)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetCpuRange(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	start := chi.URLParam(r, "start")
	end := chi.URLParam(r, "end")
	step := chi.URLParam(r, "step")
	replaceTo:="instance='"+nodeID+"'"
	basic:="100 - (avg by (instance) (irate(node_cpu{job='rbd_node', mode='idle',#to_replace#}[5m])) * 100)"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.QueryRange(expr,start,end,step)
	if err != nil {
		err.Handle(r, w)
		return
	}
	for _,v:=range resp.Data.Result{
		v.Value=v.Values
		v.Values=nil
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetMemRange(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	start := chi.URLParam(r, "start")
	end := chi.URLParam(r, "end")
	step := chi.URLParam(r, "step")
	replaceTo:="instance='"+nodeID+"'"

	basic:="100 - ((node_memory_MemFree{job='rbd_node',#to_replace#} %2B node_memory_Cached{job='rbd_node',#to_replace#} %2B node_memory_Buffers{job='rbd_node',#to_replace#})/node_memory_MemTotal) * 100"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.QueryRange(expr,start,end,step)
	if err != nil {
		err.Handle(r, w)
		return
	}
	for _,v:=range resp.Data.Result{
		v.Value=v.Values
		v.Values=nil
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetDiskRange(w http.ResponseWriter, r *http.Request) {
	basic:="100 - (avg by (instance) (node_filesystem_free{job='rbd_node',#to_replace#,fstype!~'rootfs|selinuxfs|nsfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'} / node_filesystem_size{job='rbd_node',#to_replace#,fstype!~'rootfs|selinuxfs|nsfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*'}) * 100)"
	nodeID := chi.URLParam(r, "node_id")
	start := chi.URLParam(r, "start")
	end := chi.URLParam(r, "end")
	step := chi.URLParam(r, "step")
	replaceTo:="instance='"+nodeID+"'"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.QueryRange(expr,start,end,step)
	if err != nil {
		err.Handle(r, w)
		return
	}
	for _,v:=range resp.Data.Result{
		v.Value=v.Values
		v.Values=nil
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetLoad1Range(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	start := chi.URLParam(r, "start")
	end := chi.URLParam(r, "end")
	step := chi.URLParam(r, "step")
	replaceTo:="instance='"+nodeID+"'"
	basic:="node_load1{job='rbd_node',#to_replace#}"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.QueryRange(expr,start,end,step)
	if err != nil {
		err.Handle(r, w)
		return
	}
	for _,v:=range resp.Data.Result{
		v.Value=v.Values
		v.Values=nil
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetExpr(w http.ResponseWriter, r *http.Request) {
	var expr model.Expr
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &expr.Body, nil); !ok {
		return
	}
	logrus.Infof(expr.Body.Expr)

	resp, err := prometheusAPI.Query(expr.Body.Expr)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}
func GetLoad1(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	replaceTo:="instance='"+nodeID+"'"
	basic:="node_load1{job='rbd_node',#to_replace#}"
	expr:=replaceSelector(basic,replaceTo)
	resp, err := prometheusAPI.Query(expr)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}



