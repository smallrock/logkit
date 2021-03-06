package mutate

import (
	"errors"

	"github.com/qiniu/logkit/transforms"
	. "github.com/qiniu/logkit/utils/models"
)

var (
	_ transforms.StatsTransformer = &PandoraKeyConvert{}
	_ transforms.Transformer      = &PandoraKeyConvert{}
)

type PandoraKeyConvert struct {
	stats StatsInfo
}

func (g *PandoraKeyConvert) RawTransform(datas []string) ([]string, error) {
	return datas, errors.New("pandora_key_convert transformer not support rawTransform")
}

func (g *PandoraKeyConvert) Transform(datas []Data) ([]Data, error) {
	for i, v := range datas {
		datas[i] = deepConvertKey(v)
	}

	g.stats, _ = transforms.SetStatsInfo(nil, g.stats, 0, int64(len(datas)), g.Type())
	return datas, nil
}

func deepConvertKey(data map[string]interface{}) map[string]interface{} {
	newData := make(map[string]interface{})
	for k, v := range data {
		nk := PandoraKey(k)
		if nv, ok := v.(map[string]interface{}); ok {
			v = deepConvertKey(nv)
		}
		newData[nk] = v
	}
	return newData
}

func (g *PandoraKeyConvert) Description() string {
	//return "pandora_key_convert can convert data key name to valid pandora key"
	return "将数据中的key名称中不合Pandora字段名规则的字符转为下划线, 如 a.b/c 改为 a_b_c"
}

func (g *PandoraKeyConvert) Type() string {
	return "pandora_key_convert"
}

func (g *PandoraKeyConvert) SampleConfig() string {
	return `{
		"type":"pandora_key_convert"
	}`
}

func (g *PandoraKeyConvert) ConfigOptions() []Option {
	return []Option{}
}

func (g *PandoraKeyConvert) Stage() string {
	return transforms.StageAfterParser
}

func (g *PandoraKeyConvert) Stats() StatsInfo {
	return g.stats
}

func (g *PandoraKeyConvert) SetStats(err string) StatsInfo {
	g.stats.LastError = err
	return g.stats
}

func init() {
	transforms.Add("pandora_key_convert", func() transforms.Transformer {
		return &PandoraKeyConvert{}
	})
}
