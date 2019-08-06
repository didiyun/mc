/*
 * Minio Client (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"
	"strings"

	"github.com/minio/minio/pkg/quick"
	"github.com/didiyun/mc/pkg/probe"
	"github.com/didiyun/mc/pkg/console"
)

func fixConfig() {
	// Fix config V3
	fixConfigV3()
	// Fix config V6
	fixConfigV6()
	// Fix config V6 for hosts
	fixConfigV6ForHosts()

	/* No more fixing job. Here after we bump the version for changes always.
	 */
}

/////////////////// Broken Config V3 ///////////////////
type brokenHostConfigV3 struct {
	AccessKeyID     string
	SecretAccessKey string
}

type brokenConfigV3 struct {
	Version string
	ACL     string
	Access  string
	Aliases map[string]string
	Hosts   map[string]brokenHostConfigV3
}

// newConfigV3 - get new config broken version 3.
func newBrokenConfigV3() *brokenConfigV3 {
	conf := new(brokenConfigV3)
	conf.Version = "3"
	conf.Aliases = make(map[string]string)
	conf.Hosts = make(map[string]brokenHostConfigV3)
	return conf
}

// Fix config version `3`. Some v3 config files are written without
// proper hostConfig JSON tags. They may also contain unused ACL and
// Access fields. Rewrite the hostConfig with proper fields using JSON
// tags and drop the unused (ACL, Access) fields.
func fixConfigV3() {
	if !isMcConfigExists() {
		return
	}
	brokenCfgV3 := newBrokenConfigV3()
	brokenMcCfgV3, e := quick.LoadConfig(mustGetMcConfigPath(), nil, brokenCfgV3)
	fatalIf(probe.NewError(e), "Unable to load config.")

	if brokenMcCfgV3.Version() != "3" {
		return
	}

	cfgV3 := newConfigV3()
	isMutated := false
	for k, v := range brokenMcCfgV3.Data().(*brokenConfigV3).Aliases {
		cfgV3.Aliases[k] = v
	}

	for host, brokenHostCfgV3 := range brokenMcCfgV3.Data().(*brokenConfigV3).Hosts {

		// If any of these fields contains any real value anytime,
		// it means we have already fixed the broken configuration.
		// We don't have to regenerate again.
		if brokenHostCfgV3.AccessKeyID != "" && brokenHostCfgV3.SecretAccessKey != "" {
			isMutated = true
		}

		// Use the correct hostConfig with JSON tags in it.
		cfgV3.Hosts[host] = hostConfigV3{
			AccessKeyID:     brokenHostCfgV3.AccessKeyID,
			SecretAccessKey: brokenHostCfgV3.SecretAccessKey,
		}
	}

	// We blindly drop ACL and Access fields from the broken config v3.

	if isMutated {
		mcNewConfigV3, e := quick.NewConfig(cfgV3, nil)
		fatalIf(probe.NewError(e), "Unable to initialize quick config for config version `3`.")

		e = mcNewConfigV3.Save(mustGetMcConfigPath())
		fatalIf(probe.NewError(e), "Unable to save config version `3`.")

		console.Infof("Successfully fixed %s broken config for version `3`.\n", mustGetMcConfigPath())
	}
}

// If the host key does not have http(s), fix it.
func fixConfigV6ForHosts() {
	if !isMcConfigExists() {
		return
	}

	brokenMcCfgV6, e := quick.LoadConfig(mustGetMcConfigPath(), nil, newConfigV6())
	fatalIf(probe.NewError(e), "Unable to load config.")

	if brokenMcCfgV6.Version() != "6" {
		return
	}

	newCfgV6 := newConfigV6()
	isMutated := false

	// Copy aliases.
	for k, v := range brokenMcCfgV6.Data().(*configV6).Aliases {
		newCfgV6.Aliases[k] = v
	}

	url := &clientURL{}
	// Copy hosts.
	for host, hostCfgV6 := range brokenMcCfgV6.Data().(*configV6).Hosts {
		// Already fixed - Copy and move on.
		if strings.HasPrefix(host, "https") || strings.HasPrefix(host, "http") {
			newCfgV6.Hosts[host] = hostCfgV6
			continue
		}

		// If host entry does not contain "http(s)", introduce a new entry and delete the old one.
		if host == "s3.amazonaws.com" || host == "storage.googleapis.com" ||
			host == "localhost:9000" || host == "127.0.0.1:9000" ||
			host == "play.minio.io:9000" || host == "dl.minio.io:9000" {
			console.Infoln("Found broken host entries, replacing " + host + " with https://" + host + ".")
			url.Host = host
			url.Scheme = "https"
			url.SchemeSeparator = "://"
			newCfgV6.Hosts[url.String()] = hostCfgV6
			isMutated = true
			continue
		}
	}

	if isMutated {
		// Save the new config back to the disk.
		mcCfgV6, e := quick.NewConfig(newCfgV6, nil)
		fatalIf(probe.NewError(e), "Unable to initialize quick config for config version `v6`.")

		e = mcCfgV6.Save(mustGetMcConfigPath())
		fatalIf(probe.NewError(e), "Unable to save config version `v6`.")
	}
}

// fixConfigV6 - fix all the unnecessary glob URLs present in existing config version 6.
func fixConfigV6() {
	if !isMcConfigExists() {
		return
	}
	config, e := quick.NewConfig(newConfigV6(), nil)
	fatalIf(probe.NewError(e), "Unable to initialize config.")

	e = config.Load(mustGetMcConfigPath())
	fatalIf(probe.NewError(e).Trace(mustGetMcConfigPath()), "Unable to load config.")

	if config.Data().(*configV6).Version != "6" {
		return
	}

	newConfig := new(configV6)
	isMutated := false
	newConfig.Aliases = make(map[string]string)
	newConfig.Hosts = make(map[string]hostConfigV6)
	newConfig.Version = "6"
	newConfig.Aliases = config.Data().(*configV6).Aliases
	for host, hostCfg := range config.Data().(*configV6).Hosts {
		if strings.Contains(host, "*") {
			fatalIf(errInvalidArgument(),
				fmt.Sprintf("Glob style `*` pattern matching is no longer supported. Please fix `%s` entry manually.", host))
		}
		if strings.Contains(host, "*s3*") || strings.Contains(host, "*.s3*") {
			console.Infoln("Found glob url, replacing " + host + " with s3.amazonaws.com")
			newConfig.Hosts["s3.amazonaws.com"] = hostCfg
			isMutated = true
			continue
		}
		if strings.Contains(host, "s3*") {
			console.Infoln("Found glob url, replacing " + host + " with s3.amazonaws.com")
			newConfig.Hosts["s3.amazonaws.com"] = hostCfg
			isMutated = true
			continue
		}
		if strings.Contains(host, "*amazonaws.com") || strings.Contains(host, "*.amazonaws.com") {
			console.Infoln("Found glob url, replacing " + host + " with s3.amazonaws.com")
			newConfig.Hosts["s3.amazonaws.com"] = hostCfg
			isMutated = true
			continue
		}
		if strings.Contains(host, "*storage.googleapis.com") {
			console.Infoln("Found glob url, replacing " + host + " with storage.googleapis.com")
			newConfig.Hosts["storage.googleapis.com"] = hostCfg
			isMutated = true
			continue
		}
		if strings.Contains(host, "localhost:*") {
			console.Infoln("Found glob url, replacing " + host + " with localhost:9000")
			newConfig.Hosts["localhost:9000"] = hostCfg
			isMutated = true
			continue
		}
		if strings.Contains(host, "127.0.0.1:*") {
			console.Infoln("Found glob url, replacing " + host + " with 127.0.0.1:9000")
			newConfig.Hosts["127.0.0.1:9000"] = hostCfg
			isMutated = true
			continue
		}
		// Other entries are hopefully OK. Copy them blindly.
		newConfig.Hosts[host] = hostCfg
	}

	if isMutated {
		newConf, e := quick.NewConfig(newConfig, nil)
		fatalIf(probe.NewError(e), "Unable to initialize newly fixed config.")

		e = newConf.Save(mustGetMcConfigPath())
		fatalIf(probe.NewError(e).Trace(mustGetMcConfigPath()), "Unable to save newly fixed config path.")
		console.Infof("Successfully fixed %s broken config for version `6`.\n", mustGetMcConfigPath())
	}
}
