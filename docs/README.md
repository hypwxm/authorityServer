swagger https://github.com/go-swagger/go-swagger/releases

swagger generate spec -o ./docs/swagger.json // 根据 swagger 规范 创建 swagger.json 规范文档
swagger serve -F=swagger ./docs/swagger.json // 启动一个 http 服务同时将 json 文档放入http://petstore.swagger.io 执行

解释
// swagger:operaion [请求方式(可以是 GET\PUT\DELETE\POST\PATCH)] [url:请求地址] [tag] [operation id] （同一标签的属于同一类，）
// --- 这个部分下面是 YAML 格式的 swagger 规范.确保您的缩进是一致的和正确的
// summary: 标题
// description: 描述
// parametres: 下面是参数了
// - name: 参数名
in: [header|body|query|path] 参数的位置 header 和 body 之 http 的 header 和 body 位置。 query 是 http://url?query path 就是 url 里面的替换信息
description: 描述
type: 类型
required: 是否必须
// responses: 响应
// 200：
// 404：
