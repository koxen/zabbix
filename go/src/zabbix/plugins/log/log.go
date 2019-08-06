/*
** Zabbix
** Copyright (C) 2001-2019 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

package log

import (
	"runtime"
	"time"
	"unsafe"
	"zabbix/internal/agent"
	"zabbix/internal/plugin"
	"zabbix/pkg/itemutil"
	"zabbix/pkg/zbxlib"
)

// Plugin -
type Plugin struct {
	plugin.Base
}

func (p *Plugin) Configure(options map[string]string) {
	zbxlib.SetMaxLinesPerSecond(agent.Options.MaxLinesPerSecond)
}

type metadata struct {
	key       string
	params    []string
	blob      unsafe.Pointer
	lastcheck time.Time
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	meta := ctx.Meta()

	var data *metadata
	if meta.Data == nil {
		data = &metadata{key: key, params: params}
		meta.Data = data
		runtime.SetFinalizer(data, func(d *metadata) { zbxlib.FreeActiveMetric(d.blob) })
	} else {
		data = meta.Data.(*metadata)
		if !itemutil.CompareKeysParams(key, params, data.key, data.params) {
			p.Debugf("item %d key has been changed, resetting log metadata", ctx.ItemID())
			zbxlib.FreeActiveMetric(data.blob)
			data.blob = nil
			data.key = key
			data.params = params
		}
	}

	if data.blob == nil {
		var err error
		if data.blob, err = zbxlib.NewActiveMetric(key, params, meta.LastLogsize(), meta.Mtime()); err != nil {
			return nil, err
		}
	}

	// with flexible checks there are no guaranteed refresh time,
	// so using number of seconds elapsed since last check
	now := time.Now()
	var refresh int
	if data.lastcheck.IsZero() {
		refresh = 1
	} else {
		refresh = int((now.Sub(data.lastcheck) + time.Second/2) / time.Second)
	}

	logitem := zbxlib.LogItem{Itemid: ctx.ItemID(), Results: make([]*plugin.Result, 0)}
	zbxlib.ProcessLogCheck(data.blob, &logitem, refresh)
	data.lastcheck = now

	if len(logitem.Results) != 0 {
		return logitem.Results, nil
	}
	return nil, nil
}

var impl Plugin

func init() {
	plugin.RegisterMetric(&impl, "log", "log", "Log file monitoring.")
	plugin.RegisterMetric(&impl, "log", "logrt", "Log file monitoring with log rotation support.")
	plugin.RegisterMetric(&impl, "log", "log.count", "Count of matched lines in log file monitoring.")
	plugin.RegisterMetric(&impl, "log", "logrt.count",
		"Count of matched lines in log file monitoring with log rotation support.")
}
