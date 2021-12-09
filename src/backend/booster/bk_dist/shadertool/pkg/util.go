/*
 * Copyright (c) 2021 THL A29 Limited, a Tencent company. All rights reserved
 *
 * This source code file is licensed under the MIT License, you may obtain a copy of the License at http://opensource.org/licenses/MIT
 *
 */

package pkg

import (
	"io/ioutil"
	"runtime"
	"strings"
	"sync/atomic"

	dcSDK "build-booster/bk_dist/common/sdk"
	"build-booster/bk_dist/shadertool/common"
	"build-booster/common/blog"
	"build-booster/common/codec"
)

func defaultCPULimit(custom int) int {
	if custom > 0 {
		return custom
	}
	return runtime.NumCPU() - 2
}

func resolveActionJSON(filename string) (*common.UE4Action, error) {
	blog.Infof("resolve action json file %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		blog.Errorf("failed to read action json file %s with error %v", filename, err)
		return nil, err
	}

	var t common.UE4Action
	if err = codec.DecJSON(data, &t); err != nil {
		blog.Errorf("failed to decode json content[%s] failed: %v", string(data), err)
		return nil, err
	}

	return &t, nil
}

func resolveApplyJSON(filename string) (*common.ApplyParameters, error) {
	blog.Infof("resolve apply json file %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		blog.Errorf("failed to read apply json file %s with error %v", filename, err)
		return nil, err
	}

	var t common.ApplyParameters
	if err = codec.DecJSON(data, &t); err != nil {
		blog.Errorf("failed to decode json content[%s] failed: %v", string(data), err)
		return nil, err
	}

	return &t, nil
}

func resolveOutputJSON(filename string) (*map[string]string, error) {
	blog.Infof("resolve output env json file %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		blog.Errorf("failed to read json file %s with error %v", filename, err)
		return nil, err
	}

	var m map[string]string
	if err = codec.DecJSON(data, &m); err != nil {
		blog.Errorf("failed to decode json content[%s] failed: %v", string(data), err)
		return nil, err
	}

	return &m, nil
}

func resolveToolChainJSON(filename string) (*dcSDK.Toolchain, error) {
	blog.Debugf("resolve tool chain json file %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		blog.Errorf("failed to read tool chain json file %s with error %v", filename, err)
		return nil, err
	}

	var t dcSDK.Toolchain
	if err = codec.DecJSON(data, &t); err != nil {
		blog.Errorf("failed to decode json content[%s] failed: %v", string(data), err)
		return nil, err
	}

	return &t, nil
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func removeaction(s []common.Action, i int) []common.Action {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

type count32 int32

func (c *count32) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

// replace which next is not in nextExcludes
func replaceWithNextExclude(s string, old byte, new string, nextExcludes []byte) string {
	if s == "" {
		return ""
	}

	if len(nextExcludes) == 0 {
		return strings.Replace(s, string(old), new, -1)
	}

	targetslice := make([]byte, 0, 0)
	nextexclude := false
	totallen := len(s)
	for i := 0; i < totallen; i++ {
		c := s[i]
		if c == old {
			nextexclude = false
			if i < totallen-1 {
				next := s[i+1]
				for _, e := range nextExcludes {
					if next == e {
						nextexclude = true
						break
					}
				}
			}
			if nextexclude {
				targetslice = append(targetslice, c)
				targetslice = append(targetslice, s[i+1])
				i++
			} else {
				targetslice = append(targetslice, []byte(new)...)
			}
		} else {
			targetslice = append(targetslice, c)
		}
	}

	return string(targetslice)
}
