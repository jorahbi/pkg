package query

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func GenQuery(dsn, path, table, pname string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "",
		FieldWithIndexTag: true,
		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: true,
	})
	db, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(db)
	tableList, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}
	tables := strings.Split(table, ",")
	tabMap := lo.Associate(tables, func(f string) (string, struct{}) {
		return f, struct{}{}
	})
	for _, tableName := range tableList {
		if _, isok := tabMap[tableName]; !isok {
			continue
		}
		mate := g.GenerateModel(tableName)
		data := Query{StructName: mate.ModelStructName, TableName: mate.TableName, PName: pname}
		for _, value := range mate.Fields {
			if value.Type == "*time.Time" || value.Type == "time.Time" {
				data.WithTime = true
				if value.Name == "CreateTime" {
					data.IsCreateTime = true
				} else if value.Name == "UpdateTime" {
					data.IsUpdateTime = true
				}
			}
			// "int" "int8" "int16" "int32" "int64" "float32" "float64"
			field := Field{Name: value.Name, Type: value.Type, ColumnName: value.ColumnName,
				ColumnComment: value.ColumnComment, MultilineComment: value.MultilineComment,
				Tag: value.Tag, GORMTag: value.GORMTag, CustomGenType: value.CustomGenType, Relation: value.Relation}
			data.Fields = append(data.Fields, field)
		}
		tpl := template.Must(template.New("query").Funcs(template.FuncMap{
			"tolow": func(s string) string {
				return strings.ToLower(s[:1]) + s[1:]
			},
			"trim": func(s, cutset string) string {
				return strings.Trim(s, cutset)
			},
		}).Parse(tpl()))
		if err != nil {
			panic(err)
		}
		filePath := fmt.Sprintf("%s/query.%s.gen.go", path, mate.TableName)
		file, err := os.OpenFile(filePath, os.O_EXCL|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件已存在", err)
			continue
		}
		//及时关闭file句柄
		defer file.Close()
		err = tpl.Execute(file, data)
		if err != nil {
			panic(err)
		}
	}
}

func tpl() string {
	return `package query

import (
	{{ if .WithTime}}
	"time"
	{{ end }}
	"{{.PName}}/dao/model"

	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

{{ range $index, $value := .Fields -}}
{{ if eq $value.Type "*time.Time" "time.Time" }}
func (m *{{tolow $.StructName}}) With{{$value.Name}}(times ...{{trim $value.Type "*"}}) []gen.Condition {
	var conds []gen.Condition 
	if len(times) > 0 && !times[0].Local().IsZero(){
		conds = append(conds, m.{{$value.Name}}.Gte(times[0]))
	}
	if len(times) > 1  && !times[1].Local().IsZero(){
		conds = append(conds, m.{{$value.Name}}.Lte(times[1]))
	}
	return conds
} 
{{ end }}
{{- end -}}

func (m *{{tolow $.StructName}}) WithFilter(m{{.StructName}} model.{{.StructName}}) []gen.Condition {
	var conds []gen.Condition 
	{{ range $index, $value := .Fields -}}
	{{ if eq $value.Type "*string" "string" }}
	if m{{$.StructName}}.{{$value.Name}} != "" {
		conds = append(conds, m.{{$value.Name}}.Eq(m{{$.StructName}}.{{$value.Name}}))
	}
	{{ else if eq $value.Type "int" "int8" "int16" "int32" "int64" "float32" "float64" "*int" "*int8" "*int16" "*int32" "*int64" "*float32" "*float64" }}
	if m{{$.StructName}}.{{$value.Name}} != 0 {
		conds = append(conds, m.{{$value.Name}}.Eq(m{{$.StructName}}.{{$value.Name}}))
	}
	{{ end }}
	{{- end -}}

	return conds
}

func (t *{{tolow $.StructName}}Do) CreateOrUpdate(m {{tolow $.StructName}}, cols []string, m{{.StructName}} ...*model.{{.StructName}}) error {
	if len(m{{.StructName}}) == 0 {
		return nil
	}
	if len(cols) == 0 {
		cols = []string{
			{{- range $index, $value := .Fields -}}
			m.{{$value.Name}}.ColumnName().String(),
		{{ end }}}
	}
	
	pk := clause.Column{Name: m.ID.ColumnName().String()}
	return t.Clauses(clause.OnConflict{
		Columns:   []clause.Column{pk},
		DoUpdates: clause.AssignmentColumns(cols), // 更新哪些字段
	}).Create(m{{.StructName}}...)
}`
}

type Query struct {
	StructName   string
	Fields       []Field
	TableName    string
	IsCreateTime bool
	IsUpdateTime bool
	WithTime     bool
	PName        string
}

type Field struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	Tag              field.Tag
	GORMTag          field.GormTag
	CustomGenType    string
	Relation         *field.Relation
}
