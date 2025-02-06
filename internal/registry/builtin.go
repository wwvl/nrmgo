package registry

// builtinRegistries 定义内置的 registry 列表
var builtinRegistries = map[string]*Info{
	"npm": {
		Name:        "npm",
		URL:         "https://registry.npmjs.org/",
		Home:        "https://www.npmjs.org",
		Description: "npm official",
	},
	"yarn": {
		Name:        "yarn",
		URL:         "https://registry.yarnpkg.com/",
		Home:        "https://yarnpkg.com",
		Description: "yarn official",
	},
	"taobao": {
		Name:        "taobao",
		URL:         "https://registry.npmmirror.com/",
		Home:        "https://npmmirror.com",
		Description: "Taobao npm mirror",
	},
	"tencent": {
		Name:        "tencent",
		URL:         "https://mirrors.tencent.com/npm/",
		Home:        "https://mirrors.tencent.com/help/npm.html",
		Description: "Tencent npm mirror",
	},
	"npmMirror": {
		Name:        "npmMirror",
		URL:         "https://skimdb.npmjs.com/registry/",
		Home:        "https://skimdb.npmjs.com/",
		Description: "npm mirror",
	},
	"huawei": {
		Name:        "huawei",
		URL:         "https://repo.huaweicloud.com/repository/npm/",
		Home:        "https://www.huaweicloud.com/special/npm-jingxiang.html",
		Description: "Huawei npm mirror",
	},
	"ustc": {
		Name:        "ustc",
		URL:         "https://npmreg.proxy.ustclug.org/",
		Home:        "https://mirrors.ustc.edu.cn/help/npm.html",
		Description: "USTC npm mirror",
	},
	"nju": {
		Name:        "nju",
		URL:         "https://repo.nju.edu.cn/repository/npm/",
		Home:        "https://doc.nju.edu.cn/books/35f4a/page/npm",
		Description: "NJU npm mirror",
	},
}

// GetBuiltinRegistry 获取指定名称的内置 registry
func GetBuiltinRegistry(name string) (*Info, bool) {
	reg, ok := builtinRegistries[name]
	return reg, ok
}

// ListBuiltinRegistries 列出所有内置的 registry
func ListBuiltinRegistries() []*Info {
	result := make([]*Info, 0, len(builtinRegistries))
	for _, reg := range builtinRegistries {
		result = append(result, reg)
	}
	return result
}
