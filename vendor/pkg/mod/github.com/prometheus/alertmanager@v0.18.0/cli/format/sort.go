// Copyright 2018 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package format

import (
	"time"

	"github.com/prometheus/alertmanager/api/v2/models"
)

type ByEndAt []models.GettableSilence

func (s ByEndAt) Len() int      { return len(s) }
func (s ByEndAt) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByEndAt) Less(i, j int) bool {
	return time.Time(*s[i].Silence.EndsAt).Before(time.Time(*s[j].Silence.EndsAt))
}

type ByStartsAt []*models.GettableAlert

func (s ByStartsAt) Len() int      { return len(s) }
func (s ByStartsAt) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByStartsAt) Less(i, j int) bool {
	return time.Time(*s[i].StartsAt).Before(time.Time(*s[j].StartsAt))
}
