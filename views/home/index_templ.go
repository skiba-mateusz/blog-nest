// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package home

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/skiba-mateusz/blog-nest/views/layouts"
import "github.com/skiba-mateusz/blog-nest/types"
import "github.com/skiba-mateusz/blog-nest/views/blogs"
import "fmt"

type IndexData struct {
	Title       string
	Categories  []types.Category
	Blogs       []types.Blog
	LatestBlogs []types.Blog
	TotalBlogs  int
	Page        int
	TotalPages  int
	User        *types.User
}

func Index(data IndexData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			templ_7745c5c3_Err = blogs.Latest(data.LatestBlogs).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"container container--large grid grid-2-1\"><section class=\"py-32\"><h2 class=\"heading\">More Blogs</h2>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = blogs.List(blogs.ListData{
				Blogs:      data.Blogs,
				TotalBlogs: data.TotalBlogs,
				Page:       data.Page,
				TotalPages: data.TotalPages,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</section><div class=\"flow py-32\"><form hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/blog/page/%d?search_query=%s", data.Page, ""))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/home/index.templ`, Line: 33, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#blogs\" hx-swap=\"outerHTML\" hx-trigger=\"change\"><fieldset class=\"radio-group\"><legend class=\"text-semi-bold\">Category</legend> <input id=\"Movies\" class=\"sr-only\" type=\"radio\" name=\"category\" value=\"Movies\"> <label for=\"Movies\" class=\"radio-group__btn\">Movies</label> <input id=\"Travelling\" class=\"sr-only\" type=\"radio\" name=\"category\" value=\"Travelling\"> <label for=\"Travelling\" class=\"radio-group__btn\">Travelling</label> <input id=\"Studying\" class=\"sr-only\" type=\"radio\" name=\"category\" value=\"Studying\"> <label for=\"Studying\" class=\"radio-group__btn\">Studying</label> <input id=\"Books\" class=\"sr-only\" type=\"radio\" name=\"category\" value=\"Books\"> <label for=\"Books\" class=\"radio-group__btn\">Books</label> <input id=\"All\" class=\"sr-only\" type=\"radio\" name=\"category\" value=\"\" checked> <label for=\"All\" class=\"radio-group__btn\">All</label></fieldset></form></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Base(data.Title, data.User).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
