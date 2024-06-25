// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package blogs

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/skiba-mateusz/blog-nest/views/layouts"
import "github.com/skiba-mateusz/blog-nest/types"
import "github.com/skiba-mateusz/blog-nest/forms"
import "github.com/skiba-mateusz/blog-nest/views/components"
import "strconv"

type CreateData struct {
	Title      string
	BlogForm   *forms.Form
	User       *types.User
	Categories []types.Category
}

func categoriesOptions(categories []types.Category) []components.OptionProps {
	var c []components.OptionProps
	for _, category := range categories {
		strID := strconv.Itoa(category.ID)
		c = append(c, components.OptionProps{Value: strID, Label: category.Name})
	}
	return c
}

func CreateBlogForm(categories []types.Category, form *forms.Form) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form id=\"create-blog-form\" class=\"grid\" enctype=\"multipart/form-data\" hx-post=\"/blog/create\" hx-trigger=\"create-blog\" hx-swap=\"outerHTML\"><input id=\"create-blog-content\" type=\"hidden\" name=\"content\"> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(form.Errors.Get("content")) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p class=\"alert\" role=\"alert\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(form.Errors.Get("content"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/blogs/create.templ`, Line: 29, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = components.Input(components.InputProps{
			Typ:         "file",
			Name:        "thumbnail",
			Label:       "Thumbnail",
			Placeholder: "",
			Value:       "",
			Error:       form.Errors.Get("thumbnail"),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Input(components.InputProps{
			Typ:         "text",
			Name:        "title",
			Label:       "Title",
			Placeholder: "Title",
			Value:       form.Values.Get("title"),
			Error:       form.Errors.Get("title"),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Select(components.SelectProps{
			Label:         "Category",
			Name:          "category",
			Error:         form.Errors.Get("category"),
			SelectedValue: form.Values.Get("category"),
			Placeholder: components.OptionProps{
				Label: "Select category",
			},
			Options: categoriesOptions(categories),
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func Create(data CreateData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var4 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"my-24\"><div class=\"container container--small\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.GoBack().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"panel flow\"><h1 class=\"heading\">Create blog</h1><div class=\"grid\" style=\"--grid-spacer: var(--size-8)\"><span>Blog content</span><div><div id=\"editor\"></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = CreateBlogForm(data.Categories, data.BlogForm).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button id=\"create-blog-btn\" class=\"btn btn--full btn--secondary\">Create Blog</button></div></div></section><script>\n            const toolbarOptions = [\n                ['bold', 'italic', 'underline'],\n                [{ 'list': 'ordered'}, { 'list': 'bullet' }],                     \n                [{ 'header': [2, 3, 4, 5, 6, false] }],\n                [{ 'color': [] }],          \n                [{ 'align': [] }],    \n                ['code-block'],                            \n            ];\n            const quill = new Quill('#editor', {\n                modules: {\n                    toolbar: toolbarOptions\n                },\n                theme: 'snow'\n            });\n            document.querySelector(\"#create-blog-btn\").addEventListener(\"click\", () => {\n                content = quill.getSemanticHTML()\n                document.querySelector(\"#create-blog-content\").value = content\n                htmx.trigger(document.querySelector(\"#create-blog-form\"), \"create-blog\")\n            })\n        </script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Base(data.Title, data.User).Render(templ.WithChildren(ctx, templ_7745c5c3_Var4), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
