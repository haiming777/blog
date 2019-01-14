Blog

api:

```json
interface:[/user/create]
requestMethod : post
param:name,password

interface:[/user/update]
reqeustMehtod:PUT
param:name,password

interface:[/signin]
requestMethod:post
param:name,password

interface:[/category/create]
requestMethod:post
param:name,parent_id

interface:[/categories]
requestMethod:get
param:{empty}

interface:[/sub-categories]
requestMethod:get
param:{id}//id：parent_id

interface:[/posts/create]
requestMethod:post
param:{
	"summary":"summary",
	"content":"content",
	"author":"author",
	"category_id":2,
	"status":"push"
}

interface:[/posts]
requestMethod:get
param:{1:{empty},2:{author},3:{category_id}}//1.查询全部，不传参数，2.根据作者查询，3.根据子分类id查询

interface:[/posts-detail]
requestMethod:get
param:{id}//帖子id
```