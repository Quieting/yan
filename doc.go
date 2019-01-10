// Package yan 是一个 web 框架, 基于 net/http 包
// 处理 http 请求参数, 并逐步演化成 web 服务器
// 首先实现 mvc 结构中的控制层功能, 实现效果：RESTful 风格路由, 非侵入式框架, 达到的效果, 除了导入框架包外不会在业务逻辑中调用框架函数, 路由采用树形结构查找, 通过约定实现
// 接着实现 mvc 结构中的视图层
// mvc 的 model 层将不在此项目中实现
package yan
