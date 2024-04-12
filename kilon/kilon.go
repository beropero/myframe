package kilon

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		origin      *Origin
	}

	Origin struct {
		*RouterGroup
		router        *router
		groups        []*RouterGroup
		htmlTemplates *template.Template
		funcMap       template.FuncMap
	}
)

func New() *Origin {
	origin := &Origin{router: newRouter()}
	origin.RouterGroup = &RouterGroup{origin: origin}
	origin.groups = []*RouterGroup{origin.RouterGroup}
	return origin
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.origin.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, hander HandlerFunc) {
	group.addRoute("GET", pattern, hander)
}

func (group *RouterGroup) POST(pattern string, hander HandlerFunc) {
	group.addRoute("POST", pattern, hander)
}

func (origin *Origin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range origin.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	ctx := newContext(w, req)
	ctx.handlers = middlewares
	ctx.origin = origin
	origin.router.handle(ctx)
}

func (origin *Origin) Run(addr string) (err error) {
	return http.ListenAndServe(addr, origin)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	origin := group.origin
	newGroup := &RouterGroup{
		parent: group,
		prefix: group.prefix + prefix,
		origin: origin,
	}
	origin.groups = append(origin.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

// 创建一个处理文件请求的handler，用于路由注册
func (group *RouterGroup) creatStaticHandler(comp string, fs http.FileSystem) HandlerFunc {
	// 将分组路由前缀与用户定义的路由合并成路由地址
	pattern := path.Join(group.prefix, comp)
	// 创建文件服务器
	fileServer := http.StripPrefix(pattern, http.FileServer(fs))
	// 将处理文件请求的handler返回
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return // 如果找不到文件返回404状态码
		}
		// 调用文件服务器的ServeHTTP方法，里面已经实现将文件返回到逻辑
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (group *RouterGroup) Static(comp string, root string) {
	handler := group.creatStaticHandler(comp, http.Dir(root))
	urlPattern := path.Join(comp, "/*filepath")
	group.GET(urlPattern, handler)
}


func (origin *Origin) SetFuncMap(funcMap template.FuncMap){
	origin.funcMap = funcMap
}

func (origin *Origin) LoadHTMLGlob(pattern string){
	origin.htmlTemplates =  template.Must(template.New("").Funcs(origin.funcMap).ParseGlob(pattern))
}

func Default() *Origin {
	origin := New()
	origin.Use(Logger(), Recovery()) // 注册中间件
	return origin
}